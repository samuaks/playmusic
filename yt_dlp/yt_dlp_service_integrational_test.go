package yt_dlp

import (
	"context"
	"testing"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func TestReturnsURL(t *testing.T) {
	ytdlp.MustInstall(context.TODO(), nil)

	url, err := GetStreamURL("https://www.youtube.com/watch?v=5EpyN_6dqyk")
	if err != nil {
		t.Fatal("error occurred while getting audio stream URL:", err)
	}

	if url == "" {
		t.Fatal("expected an audio stream URL, got an empty string")
	}

	t.Logf("%s", url)
}

func TestShouldGetYTVideoInfo(t *testing.T) {
	ytdlp.MustInstall(context.TODO(), nil)

	url, title, duration, err := GetYTVideoInfo("Doomed") //will actually for test
	if err != nil {
		t.Fatal("Can't get YT video info:", err)
	}

	if url == "" || title == "" || duration == 0 {
		t.Fatal("Incomplete info from the request.")
	}

	t.Logf("%s, %s, %s", url, title, duration)
}

// not mine test :C
func TestShouldGetThePipe(t *testing.T) {
	resInst, err := ytdlp.Install(context.TODO(), nil)
	if err != nil {
		t.Fatal("Can't install and use yt-dlp.")
	}

	SetBinaryPath(resInst.Executable)
	ytVideoURL := "https://www.youtube.com/watch?v=5EpyN_6dqyk"

	stdout, cmd, err := GetAudioStreamPipe(ytVideoURL)
	if err != nil {
		t.Fatalf("Failed to start stream pipe: %v", err)
	}

	// delaying killing/closing
	defer func() {
		_ = stdout.Close()
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// limiting pipe output
	buf := make([]byte, 8192)

	// making channel to signal when Wait() is done
	done := make(chan struct{})
	go func() {
		defer close(done)

		n, err := stdout.Read(buf)
		if err != nil {
			t.Errorf("Read error: %v", err)
			return
		}

		if n == 0 {
			t.Errorf("No data received from stream")
		}
	}()

	select {
	case <-done:
		// OK
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout: stream is stuck")
	}
}

func TestShouldGetListOfTracks(t *testing.T) {
	resInst, err := ytdlp.Install(context.TODO(), nil)
	if err != nil {
		t.Fatal("Can't install and use yt-dlp.")
	}

	SetBinaryPath(resInst.Executable)

	out, err := GetMusicJamPlaylistWithQuery("the weeknd")
	if err != nil {
		t.Fatalf("Failed to receive list of tracks: %v", err)
	}

	if out == nil {
		t.Fatal("No result from the request.")
	}

	if len(out) == 0 {
		t.Fatal("Not enough elements in the result.")
	}

	t.Logf("%s", out)
}
