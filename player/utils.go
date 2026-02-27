package player

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

var ffmpegAvailable bool

func init() {
	_, err := exec.LookPath("ffmpeg")
	ffmpegAvailable = err == nil
}

func IsFFmpegAvailable() bool {
	return ffmpegAvailable
}

type decoderFunc func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var beepDecoders = map[string]decoderFunc{
	".mp3":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return mp3.Decode(r) },
	".flac": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return flac.Decode(r) },
	".wav":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return wav.Decode(r) },
	".ogg":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return vorbis.Decode(r) },
}

func decode(path string) (beep.StreamSeekCloser, beep.Format, error) {
	ext := strings.ToLower(filepath.Ext(path))

	if decoder, ok := beepDecoders[ext]; ok {
		file, err := os.Open(path)
		if err != nil {
			return nil, beep.Format{}, err
		}
		return decoder(file)
	}

	if !IsFFmpegAvailable() {
		return nil, beep.Format{}, fmt.Errorf("unsupported file format: %s and ffmpeg is not available", ext)
	}
	return decodeWithFFmpeg(path)

}

type readSeekCloser struct {
	*bytes.Reader
}

type bufferedStreamer struct {
	beep.StreamSeeker
}

func (b bufferedStreamer) Close() error {
	return nil
}

func (r readSeekCloser) Close() error { return nil }

func decodeWithFFmpeg2(path string) (beep.StreamSeekCloser, beep.Format, error) {
	fmt.Printf("debug: decodeWithFFmpeg called for %s\n", path)
	cmd := exec.Command("ffmpeg",
		"-i", path,
		"-f", "wav",
		"-ar", "44100",
		"-ac", "2",
		"pipe:1",
	)
	cmd.Stderr = io.Discard

	out, err := cmd.Output()
	if err != nil {
		return nil, beep.Format{}, err
	}

	streamer, format, err := wav.Decode(readSeekCloser{bytes.NewReader(out)})

	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("wav.Decode failed: %w", err)
	}
	defer streamer.Close()

	buf := beep.NewBuffer(format)
	buf.Append(streamer)

	fmt.Printf("debug: buffer length: %d\n", buf.Len())

	fmt.Printf("debug: format sample rate: %v, channels: %v\n", format.SampleRate, format.NumChannels)
	fmt.Printf("debug: streamer length: %d samples\n", streamer.Len())

	return bufferedStreamer{buf.Streamer(0, buf.Len())}, format, nil

}

func decodeWithFFmpeg(path string) (beep.StreamSeekCloser, beep.Format, error) {
	tmp, err := os.CreateTemp("", "musicplayer-*.wav")
	if err != nil {
		return nil, beep.Format{}, err
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("ffmpeg", "-i", path, "-f", "wav", "-ar", "44100", "-ac", "2", "-y", tmpPath)
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		os.Remove(tmpPath)
		return nil, beep.Format{}, err
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return nil, beep.Format{}, err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		os.Remove(tmpPath)
		return nil, beep.Format{}, err
	}

	return &tempFileStreamer{StreamSeekCloser: streamer, path: tmpPath}, format, nil
}

type tempFileStreamer struct {
	beep.StreamSeekCloser
	path string
}

func (t *tempFileStreamer) Close() error {
	err := t.StreamSeekCloser.Close()
	os.Remove(t.path)
	return err
}

func ProbeDuration(path string) (time.Duration, error) {
	streamer, format, err := decode(path)
	if err != nil {
		return 0, fmt.Errorf("Could not determine duration: %w", err)
	}
	defer streamer.Close()

	samples := streamer.Len()
	if samples <= 0 {
		return 0, fmt.Errorf("Could not determine duration: invalid sample length")
	}
	return format.SampleRate.D(samples), nil
}
