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
	for i, file := range files {
		os.WriteFile(filepath.Join(dir, file), []byte("content for file "+string(rune(i+97))), 0644)
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

func TestLoadLibraryDeduplication(t *testing.T) {
	dir := t.TempDir()

	content := []byte("testing content")
	os.WriteFile(filepath.Join(dir, "song1.mp3"), content, 0644)
	os.WriteFile(filepath.Join(dir, "song2.mp3"), content, 0644)
	os.WriteFile(filepath.Join(dir, "song3.mp3"), []byte("different content"), 0644)

	tracks, err := LoadLibrary(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("Expected 2 unique tracks, got %d", len(tracks))
	}
}

func TestLoadLibraryRecursive(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "nested", "deeper")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("failed to create nested dirs: %v", err)
	}

	if err := os.WriteFile(filepath.Join(nested, "song.mp3"), []byte("recursive"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("skip"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	tracks, err := LoadLibrary(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("expected 1 recursive track, got %d", len(tracks))
	}
	if tracks[0].Filename != "song.mp3" {
		t.Fatalf("expected song.mp3, got %s", tracks[0].Filename)
	}
}

func TestLoadLibrariesSkipsMissingAndDeduplicatesAcrossDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	same := []byte("same-content")

	if err := os.WriteFile(filepath.Join(dir1, "same.mp3"), same, 0644); err != nil {
		t.Fatalf("failed to write dir1 track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "same.mp3"), same, 0644); err != nil {
		t.Fatalf("failed to write dir2 track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "unique.mp3"), []byte("unique"), 0644); err != nil {
		t.Fatalf("failed to write unique track: %v", err)
	}

	tracks, err := LoadLibraries(filepath.Join(dir1, "missing"), dir1, dir2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 2 {
		t.Fatalf("expected 2 unique tracks, got %d", len(tracks))
	}
}

func TestBackgroundLibraryDirsExcludesMedia(t *testing.T) {
	dirs := BackgroundLibraryDirs()

	for _, dir := range dirs {
		if filepath.Clean(dir) == filepath.Clean("Media") {
			t.Fatalf("Media directory must not be included in background scan dirs")
		}
	}
}

func TestBackgroundLibraryDirsAreSubsetOfDefaultDirs(t *testing.T) {
	defaultDirs := DefaultLibraryDirs()
	backgroundDirs := BackgroundLibraryDirs()

	defaultSet := make(map[string]struct{}, len(defaultDirs))
	for _, dir := range defaultDirs {
		defaultSet[filepath.Clean(dir)] = struct{}{}
	}

	for _, dir := range backgroundDirs {
		cleaned := filepath.Clean(dir)
		if _, ok := defaultSet[cleaned]; !ok {
			t.Fatalf("background dir %q must be part of DefaultLibraryDirs()", dir)
		}
	}
}
