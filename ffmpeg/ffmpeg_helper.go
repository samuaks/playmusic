package ffmpeg

import (
	"os/exec"
)

var ffmpegAvailable = false

func InitFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	if err == nil {
		ffmpegAvailable = true
	}
	return ffmpegAvailable
}

func IsFFmpegAvailable() bool {
	return ffmpegAvailable
}
