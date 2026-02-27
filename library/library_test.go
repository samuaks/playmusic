package library

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadLibrary(t *testing.T) {
	dir := t.TempDir()
	files := []string{"song1.mp3", "song2.wav", "document.txt"}
	for _, file := range files {
		os.WriteFile(filepath.Join(dir, file), []byte{}, 0644)
	}

	tracks, err := LoadLibrary(dir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("Expected 2 tracks, got %d", len(tracks))
	}
}

func TestLoadEmptyDir(t *testing.T) {
	dir := t.TempDir()
	tracks, err := LoadLibrary(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(tracks) != 0 {
		t.Fatalf("Expected 0 tracks, got %d", len(tracks))
	}
}

func TestLoadTracksInvalidDirectory(t *testing.T) {
	_, err := LoadLibrary("/non/existent/directory")
	if err == nil {
		t.Fatal("Expected an error for non-existent directory, got nil")
	}
}

func TestIsSupported(t *testing.T) {
	cases := []struct {
		filename string
		want     bool
	}{
		{"song.mp3", true},
		{"track.wav", true},
		{"document.txt", false},
		{"image.jpg", false},
		{"audio.flac", true},
		{"song.MP3", true},
	}
	for _, c := range cases {
		if got := isSupported(c.filename); got != c.want {
			t.Errorf("isSupported(%q) = %v, want %v", c.filename, got, c.want)
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
