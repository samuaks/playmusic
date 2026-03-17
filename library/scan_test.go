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

func trackEventsByType(events []ScanEvent, typ ScanEventType) []ScanEvent {
	var filtered []ScanEvent
	for _, evt := range events {
		if evt.Type == typ && evt.Track != nil {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}

func firstTrackEventByType(events []ScanEvent, typ ScanEventType) (ScanEvent, bool) {
	for _, evt := range events {
		if evt.Type == typ && evt.Track != nil {
			return evt, true
		}
	}
	return ScanEvent{}, false
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

	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 2 {
		t.Fatalf("expected 2 discovered events, got %d", len(discovered))
	}
	if len(enriched) != 2 {
		t.Fatalf("expected 2 enriched events, got %d", len(enriched))
	}
}

func TestScanForMediaEmitsDiscoveredBeforeEnriched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("one"), 0644); err != nil {
		t.Fatalf("failed to write song: %v", err)
	}

	events := collectScanEvents(t, []string{dir})

	if len(events) < 2 {
		t.Fatalf("expected at least 2 events, got %d", len(events))
	}

	if events[0].Type != ScanEventDiscovered {
		t.Fatalf("expected first event to be discovered, got %v", events[0].Type)
	}
	if events[1].Type != ScanEventEnriched {
		t.Fatalf("expected second event to be enriched, got %v", events[1].Type)
	}
	if events[0].Track == nil || events[1].Track == nil {
		t.Fatal("expected both events to carry track payloads")
	}
	if events[0].Track.Path != events[1].Track.Path {
		t.Fatalf(
			"expected same path in discovered and enriched events, got %q and %q",
			events[0].Track.Path,
			events[1].Track.Path,
		)
	}
}

func TestScanForMediaDiscoveredEventContainsBaseTrackOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("one"), 0644); err != nil {
		t.Fatalf("failed to write song: %v", err)
	}

	events := collectScanEvents(t, []string{dir})

	evt, ok := firstTrackEventByType(events, ScanEventDiscovered)
	if !ok {
		t.Fatal("expected discovered event")
	}

	if evt.Track.Path != path {
		t.Fatalf("expected path %q, got %q", path, evt.Track.Path)
	}
	if evt.Track.Filename != "song.mp3" {
		t.Fatalf("expected filename song.mp3, got %q", evt.Track.Filename)
	}
	if evt.Track.Duration != 0 {
		t.Fatalf("expected zero duration before enrichment, got %v", evt.Track.Duration)
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
	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 1 {
		t.Fatalf("expected 1 recursively discovered event, got %d", len(discovered))
	}
	if len(enriched) != 1 {
		t.Fatalf("expected 1 recursively enriched event, got %d", len(enriched))
	}
	if discovered[0].Track.Filename != "song.mp3" {
		t.Fatalf("expected discovered file song.mp3, got %q", discovered[0].Track.Filename)
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

	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 1 {
		t.Fatalf("expected scanner to skip missing dir and emit 1 discovered event, got %d", len(discovered))
	}
	if len(enriched) != 1 {
		t.Fatalf("expected scanner to skip missing dir and emit 1 enriched event, got %d", len(enriched))
	}
}

func TestScanForMediaPopulatesBasicTrackFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	events := collectScanEvents(t, []string{dir})

	evt, ok := firstTrackEventByType(events, ScanEventEnriched)
	if !ok {
		t.Fatal("expected enriched event")
	}

	track := *evt.Track
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
