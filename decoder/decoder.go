package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"playmusic/ffmpeg"
	"playmusic/yt_dlp"
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

var beepFormats = map[string]bool{
	".mp3":  true,
	".wav":  true,
	".flac": true,
	".ogg":  true,
}

var ffmpegFormats = map[string]bool{
	".aac":  true,
	".m4a":  true,
	".opus": true,
}

func IsSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if beepFormats[ext] {
		return true
	}
	return ffmpegAvailable && ffmpegFormats[ext]
}

type decoderFunc func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var beepDecoders = map[string]decoderFunc{
	".mp3":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return mp3.Decode(r) },
	".flac": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return flac.Decode(r) },
	".wav":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return wav.Decode(r) },
	".ogg":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return vorbis.Decode(r) },
}

func Decode(path string) (beep.StreamSeekCloser, beep.Format, error) {
	ext := strings.ToLower(filepath.Ext(path))

	if decoder, ok := beepDecoders[ext]; ok {
		file, err := os.Open(path)
		if err != nil {
			return nil, beep.Format{}, err
		}
		streamer, format, err := decoder(file)
		if err != nil {
			file.Close()
			return nil, beep.Format{}, err
		}
		return streamer, format, nil
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

func decodeWithFFmpeg(path string) (beep.StreamSeekCloser, beep.Format, error) {
	sampleRate, _ := ffmpeg.GetSourceSampleRateFFmpeg(path)

	tmp, err := os.CreateTemp("", "musicplayer-*.wav")
	if err != nil {
		return nil, beep.Format{}, err
	}
	tmpPath := tmp.Name()
	tmp.Close()

	if err := ffmpeg.ConvertToWav(path, tmpPath, sampleRate); err != nil {
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

// it is actual streaming: yt-dlp returns stream with bytes and ffmpeg decodes it on the fly,
// so there is no need to wait until the whole file is downloaded
func DecodeStreamUrl(url string) (beep.Streamer, beep.Format, func(), error) {
	if !IsFFmpegAvailable() {
		err := fmt.Errorf("Can't docode stream without ffmpeg.")
		return nil, beep.Format{}, nil, err
	}

	ytdlpOut, ytdlpCmd, err := yt_dlp.GetAudioStreamPipe(url)
	if err != nil {
		return nil, beep.Format{}, nil, err
	}

	ffmpegOut, ffmpegCmd, err := ffmpeg.StreamFromPipe(ytdlpOut)
	if err != nil {
		return nil, beep.Format{}, nil, err
	}

	closeFn := func() {
		//killing downstrean first and then upstream is the right way
		closePipeKillProcess(ffmpegOut, ffmpegCmd)
		closePipeKillProcess(ytdlpOut, ytdlpCmd)
	}

	// PCM format: 44100 Hz, 2 channels, signed 16-bit
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	//making a streamer that reads from ffmpeg's output and converts it to beep's format
	streamer := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		buf := make([]byte, len(samples)*4)
		read, err := ffmpegOut.Read(buf)
		if err != nil {
			return 0, false
		}
		n = read / 4

		for i := 0; i < n; i++ {
			left := int16(binary.LittleEndian.Uint16(buf[i*4 : i*4+2]))
			right := int16(binary.LittleEndian.Uint16(buf[i*4+2 : i*4+4]))
			samples[i][0] = float64(left) / (1 << 15)
			samples[i][1] = float64(right) / (1 << 15)
		}
		return n, true
	})

	return streamer, format, closeFn, nil
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

func probeDurationWithBeep(path string) (time.Duration, error) {
	streamer, format, err := Decode(path)
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

func ProbeDuration(path string) (time.Duration, error) {
	duration, err := probeDurationWithBeep(path)
	if err == nil && duration > 0 {
		return duration, nil
	}
	if !IsFFmpegAvailable() {
		return 0, fmt.Errorf("could not determine duration and ffmpeg is not available")
	}
	return ffmpeg.ProbeDurationFFmpeg(path)
}

func closePipeKillProcess(pipe io.ReadCloser, cmd *exec.Cmd) {
	// closing pipe
	if pipe != nil {
		err := pipe.Close()
		if err != nil {
			fmt.Println("Error closing pipe:", err)
		}
	}

	// killing process
	if cmd == nil || cmd.Process == nil {
		return
	}

	done := make(chan struct{}) //channel to signal when Wait() is done

	go func() {
		_ = cmd.Wait() // blocks the thread so performing in separate goroutine
		close(done)
	}()

	// waiting for process to exit gracefully
	select {
	case <-done: //if done -> going out
		return
	case <-time.After(500 * time.Millisecond): //if timeout -> killing the process(500ms is the tick of the UI with less unstable)
		//fmt.Println("Forcing the process to end (Killing)")
	}

	err := cmd.Process.Kill()
	if err != nil && !strings.Contains(err.Error(), "Access is denied") {
		fmt.Println("Error killing the command:", err)
	}

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Wait() timeout, process might be stuck")
	}
}
