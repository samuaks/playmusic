package yt_dlp

import (
	"context"
	"testing"

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
