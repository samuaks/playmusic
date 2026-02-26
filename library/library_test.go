package library

import (
	"os"
	"path/filepath"
	"testing"
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
