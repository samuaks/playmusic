package library

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func collectScanEvents(t *testing.T, dirs []string) []ScanEvent {
	t.Helper()

	ch := make(chan ScanEvent)
	go ScanForMedia(context.Background(), dirs, ch)

	var events []ScanEvent
	for evt := range ch {
		events = append(events, evt)
	}

	return events
}

func onlyTracks(events []ScanEvent) []Track {
	var tracks []Track
	for _, evt := range events {
		if evt.Track != nil {
			tracks = append(tracks, *evt.Track)
		}
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

	events := collectScanEvents(t, []string{dir})
	tracks := onlyTracks(events)

	if len(tracks) != 2 {
		t.Fatalf("expected 2 supported tracks, got %d", len(tracks))
	}
}

func TestScanForMediaEmitsTypedTrackEvents(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "song.mp3"), []byte("one"), 0644); err != nil {
		t.Fatalf("failed to write song: %v", err)
	}

	events := collectScanEvents(t, []string{dir})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt := events[0]
	if evt.Type != ScanEventEnriched {
		t.Fatalf("expected event type %v, got %v", ScanEventEnriched, evt.Type)
	}
	if evt.Track == nil {
		t.Fatal("expected track payload")
	}
	if evt.Err != nil {
		t.Fatalf("expected nil error, got %v", evt.Err)
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

	events := collectScanEvents(t, []string{root})
	tracks := onlyTracks(events)

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

	events := collectScanEvents(t, []string{
		filepath.Join(dir, "missing"),
		dir,
	})
	tracks := onlyTracks(events)

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

	events := collectScanEvents(t, []string{dir})
	tracks := onlyTracks(events)

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

func TestScanForMediaStopsWhenContextCanceled(t *testing.T) {
	dir := t.TempDir()
	ch := make(chan ScanEvent)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	go ScanForMedia(ctx, []string{dir}, ch)

	for range ch {
		t.Fatal("expected no events after cancellation")
	}
}
