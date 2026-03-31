package yt_dlp

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"playmusic/helpers"
	"playmusic/paths"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func init() {
	jsCheck()
}

var ytDlpBinaryPath string

func SetBinaryPath(path string) {
	ytDlpBinaryPath = path
}

var jsAvailable bool

func jsCheck() {
	_, err := exec.LookPath("node")
	jsAvailable = err == nil
}

// yt-dlp wants a JavaScrypt runtime to execute some of its features.
// If Node.js is in the PATH, we can use it as the runtime for yt-dlp.
func IsJSAvailable() bool {
	return jsAvailable
}

type TrackInfo struct {
	Trackname  string
	YTVideoURL string
	Duration   time.Duration
}

// not used in the current impl
func GetStreamURL(ytVideoURL string) (string, error) {
	ytdlpCommand := ytdlp.New()

	if jsAvailable {
		ytdlpCommand = ytdlpCommand.JsRuntimes("node")
	}

	ytdlpCommand = ytdlpCommand.
		Format("bestaudio").
		Print("urls")

	output, err := ytdlpCommand.Run(context.TODO(), ytVideoURL)
	if err != nil {
		return "", fmt.Errorf("Failed to extract audio URL: %w", err)
	}

	audioURL := strings.TrimSpace(output.Stdout)

	return audioURL, nil
}

// configuration of the pipe with cmd and not package is used as
// this way connecting "producer" and "receiver" is made easier to config the "receiver"
func GetAudioStreamPipe(ytVideoURL string) (io.ReadCloser, *exec.Cmd, error) {
	cmd := exec.Command(
		ytDlpBinaryPath,
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

// not used in the current impl
func GetYTVideoInfo(query string) (string, string, time.Duration, error) {
	ytdlpCommand := ytdlp.New()

	if jsAvailable {
		ytdlpCommand = ytdlpCommand.JsRuntimes("node")
	}

	ytdlpCommand = ytdlpCommand.Print("%(webpage_url)s<<>>%(uploader)s<<>>%(title)s<<>>%(duration)s")

	out, err := ytdlpCommand.Run(context.TODO(), "ytsearch1:"+query)
	if err != nil {
		return "", "", 0, err
	}

	return extractInfoFromResult(*out)
}

func GetMusicJamPlaylistWithQueryJson(query string) ([]TrackInfo, error) {
	var playlist []TrackInfo
	ytdlpCommand := ytdlp.New()

	if jsAvailable {
		ytdlpCommand = ytdlpCommand.JsRuntimes("node")
	}

	ytdlpCommand = ytdlpCommand.
		NoWarnings().
		Quiet().
		FlatPlaylist().
		MatchFilters("duration > 120 & duration < 540").
		DumpJSON()

	out, err := ytdlpCommand.Run(context.TODO(), "ytsearch10:"+query+" topic ") //10 results
	if err != nil {
		return nil, err
	}

	entries, err := YTVideoInfoParser(*out)
	if err != nil {
		return nil, err
	}

	for _, data := range entries {
		playlist = append(playlist, TrackInfo{
			Trackname:  data.Author + " - " + data.Name,
			YTVideoURL: data.URL,
			Duration:   time.Duration(data.Duration * float64(time.Second)),
		})
	}

	return playlist, nil
}

// playlist with kinda similar tracks
func GetRecomendedWithYTVideoURL(ytVideoURL string) ([]TrackInfo, error) {
	var tracklist []TrackInfo
	ytdlpCommand := ytdlp.New()

	if jsAvailable {
		ytdlpCommand = ytdlpCommand.JsRuntimes("node")
	}

	ytdlpCommand = ytdlpCommand.
		NoWarnings().
		Quiet().
		FlatPlaylist().
		MatchFilters("duration > 120 & duration < 540").
		PlaylistItems("1-20").
		DumpJSON()

	splitURL := strings.Split(ytVideoURL, "?v=")
	ytVideoID := splitURL[1]
	mixQueryURL := "https://www.youtube.com/watch?v=" + ytVideoID + "&list=RD" + ytVideoID

	out, err := ytdlpCommand.Run(context.TODO(), mixQueryURL)
	if err != nil {
		return nil, err
	}

	entries, err := YTVideoInfoParser(*out)
	if err != nil {
		return nil, err
	}

	for _, data := range entries {
		tracklist = append(tracklist, TrackInfo{
			Trackname:  data.Author + " - " + data.Name,
			YTVideoURL: data.URL,
			Duration:   time.Duration(data.Duration * float64(time.Second)),
		})
	}

	return tracklist, nil
}

func DownloadAudio(ytVideoURL string) error {
	userMediaDir, err := paths.UserMediaDir()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(userMediaDir, "%(title)s.%(ext)s")

	ytdlpCmd := ytdlp.New()

	if jsAvailable {
		ytdlpCmd = ytdlpCmd.JsRuntimes("node")
	}

	ytdlpCmd = ytdlpCmd.
		Format("bestaudio").
		ExtractAudio().
		AudioFormat("mp3").
		Output(outputPath)

	_, err = ytdlpCmd.Run(context.TODO(), ytVideoURL)
	return err
}

func extractInfoFromResult(queryRes ytdlp.Result) (string, string, time.Duration, error) {
	parts := strings.SplitN(strings.TrimSpace(queryRes.Stdout), "<<>>", 4)
	if len(parts) < 4 {
		return "", "", 0, fmt.Errorf("invalid output")
	}

	url := parts[0]
	author := parts[1]
	name := parts[2]
	duration := strings.TrimSpace(parts[3])
	formatted, err := helpers.StringToDuration(duration)
	if err != nil {
		return "", "", 0, err
	}

	return url, author + " - " + name, formatted, nil
}
