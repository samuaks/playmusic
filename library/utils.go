package library

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var supported = map[string]bool{
	".mp3":  true,
	".flac": true,
	".wav":  true,
	".m4a":  true,
	".ogg":  true,
	".aac":  true,
	".opus": true,
}

func isSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return supported[ext]
}

type ffprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func probeDuration(path string) (time.Duration, error) {
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
