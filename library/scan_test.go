package library

import (
	"os"
	"path/filepath"
	"testing"
)

func collectScannedTracks(t *testing.T, dirs []string) []Track {
	t.Helper()

	ch := make(chan Track)
	go ScanForMedia(dirs, ch)

	var tracks []Track
	for track := range ch {
		tracks = append(tracks, track)
	}

	return tracks
}

func TestScanForMediaStreamsSupportedFiles(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "song1.mp3"), []byte("one"), 0644); err != nil {
		t.Fatalf("failed to write song1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "song2.wav"), []byte("two"), 0644); err != nil {
		t.Fatalf("failed to write song2: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("skip"), 0644); err != nil {
		t.Fatalf("failed to write notes: %v", err)
	}

	tracks := collectScannedTracks(t, []string{dir})

	if len(tracks) != 2 {
		t.Fatalf("expected 2 supported tracks, got %d", len(tracks))
	}
}

func TestScanForMediaScansRecursively(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "nested", "deeper")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("failed to create nested directories: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nested, "song.mp3"), []byte("recursive"), 0644); err != nil {
		t.Fatalf("failed to write recursive track: %v", err)
	}

	tracks := collectScannedTracks(t, []string{root})

	if len(tracks) != 1 {
		t.Fatalf("expected 1 recursively discovered track, got %d", len(tracks))
	}
	if tracks[0].Filename != "song.mp3" {
		t.Fatalf("expected discovered file song.mp3, got %q", tracks[0].Filename)
	}
}

func TestScanForMediaSkipsMissingDirs(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "song.mp3"), []byte("ok"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	tracks := collectScannedTracks(t, []string{
		filepath.Join(dir, "missing"),
		dir,
	})

	if len(tracks) != 1 {
		t.Fatalf("expected scanner to skip missing dir and return 1 track, got %d", len(tracks))
	}
}

func TestScanForMediaPopulatesBasicTrackFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	tracks := collectScannedTracks(t, []string{dir})

	if len(tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks))
	}

	track := tracks[0]
	if track.Path != path {
		t.Fatalf("expected path %q, got %q", path, track.Path)
	}
	if track.Filename != "song.mp3" {
		t.Fatalf("expected filename song.mp3, got %q", track.Filename)
	}
	if track.Trackname == "" {
		t.Fatal("expected Trackname to be populated")
	}
}
