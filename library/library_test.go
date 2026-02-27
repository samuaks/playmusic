package library

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTracks(t *testing.T) {
	dir := t.TempDir()
	files := []string{"song1.mp3", "song2.wav", "document.txt"}
	for _, file := range files {
		os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	}

	tracks, err := LoadLibrary(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("expected 2 tracks, got %d", len(tracks))
	}
}

func TestLoadTracksEmptyDir(t *testing.T) {
	dir := t.TempDir()
	tracks, err := LoadLibrary(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 0 {
		t.Fatalf("expected 0 tracks, got %d", len(tracks))
	}
}

func TestLoadTracksInvalidDirectory(t *testing.T) {
	_, err := LoadLibrary("/non/existent/directory")
	if err == nil {
		t.Fatal("expected an error for non-existent directory, got nil")
	}
}

func TestIsSupported(t *testing.T) {
	cases := []struct {
		filename    string
		ffmpegAvail bool
		want        bool
	}{
		{"song.mp3", false, true},
		{"track.wav", false, true},
		{"audio.flac", false, true},
		{"song.MP3", false, true},
		{"document.txt", false, false},
		{"image.jpg", false, false},
		{"song.m4a", false, false},
		{"song.m4a", true, true},
		{"song.aac", true, true},
		{"song.aac", false, false},
	}
	for _, c := range cases {
		got := isSupported(c.filename, c.ffmpegAvail)
		if got != c.want {
			t.Errorf("isSupported(%q, ffmpeg=%v) = %v, want %v", c.filename, c.ffmpegAvail, got, c.want)
		}
	}
}

func TestFormattedDuration(t *testing.T) {
	cases := []struct {
		duration time.Duration
		want     string
	}{
		{3*time.Minute + 45*time.Second, "3:45"},
		{10*time.Minute + 5*time.Second, "10:05"},
		{0, "0:00"},
	}
	for _, c := range cases {
		track := Track{Duration: c.duration}
		got := track.FormatDuration()
		if got != c.want {
			t.Errorf("FormatDuration() = %q, want %q", got, c.want)
		}
	}
}
