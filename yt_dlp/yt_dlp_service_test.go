package yt_dlp

import (
	"testing"
)

func TestReturnsURL(t *testing.T) {
	url, err := GetStreamURLFromYtDlp("https://www.youtube.com/watch?v=5EpyN_6dqyk")
	if err != nil {
		t.Fatal("error occurred while getting audio stream URL:", err)
	}

	if url == "" {
		t.Fatal("expected an audio stream URL, got an empty string")
	}

	t.Logf("%s", url)
}
