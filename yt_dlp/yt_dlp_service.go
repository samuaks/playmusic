package yt_dlp

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"playmusic/helpers"
	"strings"
	"time"

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

// not used in the current impl
func GetStreamURL(ytVideoURL string) (string, error) {
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

func GetAudioStreamPipe(ytVideoURL string) (io.ReadCloser, *exec.Cmd, error) {
	cmd := exec.Command(
		"yt-dlp",
		"-f", "bestaudio",
		"-o", "-",
		ytVideoURL,
	)

	cmd.Stderr = io.Discard

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return stdout, cmd, nil
}

func GetYTVideoInfo(query string) (string, string, time.Duration, error) {
	var ytdlpCommand *ytdlp.Command

	if !jsAvalable {
		ytdlpCommand = ytdlp.New().
			Print("%(webpage_url)s|%(title)s|%(duration)s")
	} else {
		ytdlpCommand = ytdlp.New().
			JsRuntimes("node").
			Print("%(webpage_url)s|%(title)s|%(duration)s")
	}

	out, err := ytdlpCommand.Run(context.TODO(), "ytsearch1:"+query)
	if err != nil {
		return "", "", 0, err
	}

	parts := strings.SplitN(strings.TrimSpace(out.Stdout), "|", 3)
	if len(parts) < 3 {
		return "", "", 0, fmt.Errorf("invalid output")
	}

	url := parts[0]
	title := parts[1]
	duration := strings.TrimSpace(parts[2])
	formatted, err := helpers.StringToDuration(duration)

	return url, title, formatted, nil
}
