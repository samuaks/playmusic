package yt_dlp

import (
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	ytdlpCheck()
	jsCheck()
}

var ytdlpAvailable bool
var jsAvalable bool

func ytdlpCheck() {
	_, err := exec.LookPath("yt-dlp")
	ytdlpAvailable = err == nil
}

func IsYTdlpAvailable() bool {
	return ytdlpAvailable
}

func jsCheck() {
	_, err := exec.LookPath("node")
	jsAvalable = err == nil
}

func IsJSAvailable() bool {
	return jsAvalable
}

func GetAudioURLFromYtDlp(videoURL string) (string, error) {
	if !ytdlpAvailable {
		return "", fmt.Errorf("yt-dlp is not available in the system. Install yt-dlp to use this feature.")
	}

	var cmd *exec.Cmd

	if !jsAvalable {
		cmd = exec.Command("yt-dlp", "-f", "bestaudio", "-g", videoURL)
	} else {
		cmd = exec.Command("yt-dlp", "--js-runtimes", "node", "-f", "bestaudio", "-g", videoURL)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to extract audio URL: %w", err)
	}

	url := strings.TrimSpace(string(output))

	return url, nil
}
