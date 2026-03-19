package ffmpeg

import (
	"encoding/json"
	"fmt"
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

func ProbeDurationWithPackage(path string) (time.Duration, error) {
	output, err := ff.Probe(path)
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

func GetSourceSampleRatePackage(path string) (int, error) {
	out, err := ff.Probe(path)
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
