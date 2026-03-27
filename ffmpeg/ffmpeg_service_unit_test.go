package ffmpeg

import (
	"fmt"
	"testing"
)

type MockFFmpeg struct {
	Output string
	Err    error
}

func (m MockFFmpeg) Probe(path string) (string, error) {
	return m.Output, m.Err
}

func TestProbeDurationSuccess(t *testing.T) {
	mock := MockFFmpeg{
		Output: `{
			"format": {
				"duration": "260.1"
			}
		}`,
		Err: nil,
	}

	duration, err := ProbeDurationFFmpeg(mock, "mockPath")
	if err != nil {
		t.Fatal("Should be nil, but it's not.", err)
	}

	t.Logf("%s", duration)
}

func TestProbeDurationErr(t *testing.T) {
	mock := MockFFmpeg{
		Output: "",
		Err:    fmt.Errorf("Ka-boom!"),
	}

	_, err := ProbeDurationFFmpeg(mock, "mockPath")
	if err == nil {
		t.Fatal("Should be nil, but it's not.", err)
	}
}
