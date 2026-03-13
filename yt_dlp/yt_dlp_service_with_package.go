package yt_dlp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

func init() {
	jsCheck()
}

var jsAvalable bool

func jsCheck() {
	_, err := exec.LookPath("node")
	jsAvalable = err == nil
}

// yt-dlp wants a JavaScrypt runtime to execute some of its features.
// If Node.js is in the PATH, we can use it as the runtime for yt-dlp.
func IsJSAvailable() bool {
	return jsAvalable
}

func GetStreamURLFromYtDlp(ytVideoURL string) (string, error) {
	var ytdlpCommand *ytdlp.Command

	if !jsAvalable {
		ytdlpCommand = ytdlp.New().
			Format("bestaudio").
			Print("urls")
	} else {
		ytdlpCommand = ytdlp.New().
			JsRuntimes("node").
			Format("bestaudio").
			Print("urls")
	}

	output, err := ytdlpCommand.Run(context.TODO(), ytVideoURL)
	if err != nil {
		return "", fmt.Errorf("Failed to extract audio URL: %w", err)
	}

	audioURL := strings.TrimSpace(output.Stdout)

	return audioURL, nil
}
