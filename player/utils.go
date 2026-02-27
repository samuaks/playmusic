package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func getSourceSampleRate(path string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return 44100, err
	}

	var result struct {
		Streams []struct {
			SampleRate string `json:"sample_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &result); err != nil || len(result.Streams) == 0 {
		return 44100, fmt.Errorf("could not parse ffprobe output: %w", err)
	}
	sampleRate, err := strconv.Atoi(result.Streams[0].SampleRate)
	if err != nil {
		return 44100, fmt.Errorf("invalid sample rate in ffprobe output: %w", err)
	}
	return sampleRate, nil
}

func decodeWithFFmpeg(path string) (beep.StreamSeekCloser, beep.Format, error) {
	sampleRate, _ := getSourceSampleRate(path)

	tmp, err := os.CreateTemp("", "musicplayer-*.wav")
	if err != nil {
		return nil, beep.Format{}, err
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("ffmpeg", "-i", path, "-f", "wav", "-ar", strconv.Itoa(sampleRate), "-ac", "2", "-y", tmpPath)
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
