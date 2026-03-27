package ffmpeg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"time"

	ff "github.com/u2takey/ffmpeg-go"
)

type ffprobeDuration struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

type ffprobeSampleRateStruct struct {
	Streams []struct {
		SampleRate string `json:"sample_rate"`
	} `json:"streams"`
}

func ProbeDurationFFmpeg(ffInt FFmpegInterface, path string) (time.Duration, error) {
	output, err := ffInt.Probe(path)
	if err != nil {
		return 0, err
	}

	var result ffprobeDuration
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return 0, err
	}
	seconds, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

func GetSourceSampleRateFFmpeg(ffInt FFmpegInterface, path string) (int, error) {
	out, err := ffInt.Probe(path)
	if err != nil {
		return 44100, err
	}

	var result ffprobeSampleRateStruct
	if err := json.Unmarshal([]byte(out), &result); err != nil || len(result.Streams) == 0 {
		return 44100, fmt.Errorf("could not parse ffprobe output: %w", err)
	}
	sampleRate, err := strconv.Atoi(result.Streams[0].SampleRate)
	if err != nil {
		return 44100, fmt.Errorf("invalid sample rate in ffprobe output: %w", err)
	}
	return sampleRate, nil
}

func ConvertToWav(inputPath, outputPath string, sampleRate int) error {
	//turning off logs, as spams in console when converting
	prev := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(prev)

	return ff.Input(inputPath).
		Output(outputPath, ff.KwArgs{
			"f":  "wav",
			"ar": strconv.Itoa(sampleRate),
			"ac": "2",
		}).
		OverWriteOutput().
		Run()
}

func StreamFromPipe(input io.Reader) (io.ReadCloser, *exec.Cmd, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0",
		"-f", "s16le",
		"-ac", "2",
		"-ar", "44100",
		"pipe:1",
	)

	cmd.Stderr = io.Discard
	cmd.Stdin = input

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return out, cmd, nil
}
