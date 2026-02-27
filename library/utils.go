package library

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"playmusic/player"
	"strconv"
	"strings"
	"time"
)

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

func isSupported(filename string, ffmpegAvailable bool) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if beepFormats[ext] {
		return true
	}
	return ffmpegAvailable && ffmpegFormats[ext]
}

func probeDuration(path string) (time.Duration, error) {
	duration, err := player.ProbeDuration(path)
	if err == nil && duration > 0 {
		return duration, nil
	}
	if !player.IsFFmpegAvailable() {
		return 0, fmt.Errorf("could not determine duration and ffmpeg is not available")
	}
	return probeWithFFprobe(path)
}

type ffprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func probeWithFFprobe(path string) (time.Duration, error) {
	command := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", path)
	output, err := command.Output()
	if err != nil {
		return 0, err
	}

	var result ffprobeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, err
	}
	seconds, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds * float64(time.Second)), nil
}
