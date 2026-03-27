package library

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	if evt.Track.Trackname != "song" {
		t.Fatalf("expected discovered trackname song, got %q", evt.Track.Trackname)
	}
	if evt.Track.Artist != "" {
		t.Fatalf("expected empty artist before enrichment, got %q", evt.Track.Artist)
	}
	if evt.Track.Title != "" {
		t.Fatalf("expected empty title before enrichment, got %q", evt.Track.Title)
	}
	if evt.Track.Album != "" {
		t.Fatalf("expected empty album before enrichment, got %q", evt.Track.Album)
	}
	if evt.Track.Genre != "" {
		t.Fatalf("expected empty genre before enrichment, got %q", evt.Track.Genre)
	}
	if evt.Track.Year != 0 {
		t.Fatalf("expected zero year before enrichment, got %d", evt.Track.Year)
	}
	if evt.Track.Duration != 0 {
		t.Fatalf("expected zero duration before enrichment, got %v", evt.Track.Duration)
	}
}

func TestScanForMediaEmitsDiscoveredBeforeEnrichmentCompletes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	originalEnrichTrack := enrichTrack
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	defer func() {
		enrichTrack = originalEnrichTrack
	}()

	enrichTrack = func(track Track) (Track, error) {
		started <- struct{}{}
		<-release
		return originalEnrichTrack(track)
	}

	ch := make(chan ScanEvent)
	go ScanForMedia(context.Background(), []string{dir}, ch)

	select {
	case evt, ok := <-ch:
		if !ok {
			t.Fatal("expected discovered event before channel close")
		}
		if evt.Type != ScanEventDiscovered {
			t.Fatalf("expected first event to be discovered, got %v", evt.Type)
		}
		if evt.Track == nil {
			t.Fatal("expected discovered event to include track payload")
		}
		if evt.Track.Path != path {
			t.Fatalf("expected discovered path %q, got %q", path, evt.Track.Path)
		}
		if evt.Track.Duration != 0 {
			t.Fatalf("expected zero duration before enrichment, got %v", evt.Track.Duration)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected discovered event before enrichment completed")
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("expected enrichment to start after discovered event")
	}

	select {
	case evt := <-ch:
		t.Fatalf("expected no enriched event before releasing enrichment, got %v", evt.Type)
	case <-time.After(100 * time.Millisecond):
	}

	close(release)

	select {
	case evt, ok := <-ch:
		if !ok {
			t.Fatal("expected enriched event after releasing enrichment")
		}
		if evt.Type != ScanEventEnriched {
			t.Fatalf("expected second event to be enriched, got %v", evt.Type)
		}
		if evt.Track == nil {
			t.Fatal("expected enriched event to include track payload")
		}
		if evt.Track.Path != path {
			t.Fatalf("expected enriched path %q, got %q", path, evt.Track.Path)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected enriched event after releasing enrichment")
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

func TestScanForMediaSignalsCompletionByClosingChannel(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "song.mp3"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	ch := make(chan ScanEvent)
	go ScanForMedia(context.Background(), []string{dir}, ch)

	count := 0
	for range ch {
		count++
	}

	if count != 2 {
		t.Fatalf("expected discovered and enriched events before close, got %d", count)
	}
}

func TestScanForMediaEmitsDiscoveredBeforeFailedEnrichment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	originalEnrichTrack := enrichTrack
	defer func() {
		enrichTrack = originalEnrichTrack
	}()

	wantErr := errors.New("enrichment failed")
	enrichTrack = func(track Track) (Track, error) {
		enriched := track
		enriched.Trackname = "Broken Song"
		return enriched, wantErr
	}

	ch := make(chan ScanEvent)
	go ScanForMedia(context.Background(), []string{dir}, ch)

	first, ok := <-ch
	if !ok {
		t.Fatal("expected discovered event")
	}
	if first.Type != ScanEventDiscovered {
		t.Fatalf("expected first event to be discovered, got %v", first.Type)
	}
	if first.Track == nil {
		t.Fatal("expected discovered event to include track payload")
	}
	if first.Track.Path != path {
		t.Fatalf("expected discovered path %q, got %q", path, first.Track.Path)
	}

	second, ok := <-ch
	if !ok {
		t.Fatal("expected enriched event")
	}
	if second.Type != ScanEventEnriched {
		t.Fatalf("expected second event to be enriched, got %v", second.Type)
	}
	if second.Track == nil {
		t.Fatal("expected enriched event to include track payload")
	}
	if second.Track.Path != path {
		t.Fatalf("expected enriched path %q, got %q", path, second.Track.Path)
	}
	if !errors.Is(second.Err, wantErr) {
		t.Fatalf("expected enrichment error %v, got %v", wantErr, second.Err)
	}
}

func TestScanForMediaSkipsDuplicateFilenameAndSizeBeforeEnrichment(t *testing.T) {
	root := t.TempDir()
	dirA := filepath.Join(root, "a")
	dirB := filepath.Join(root, "b")

	if err := os.MkdirAll(dirA, 0755); err != nil {
		t.Fatalf("failed to create dirA: %v", err)
	}
	if err := os.MkdirAll(dirB, 0755); err != nil {
		t.Fatalf("failed to create dirB: %v", err)
	}

	content := []byte("same-size")
	if err := os.WriteFile(filepath.Join(dirA, "song.mp3"), content, 0644); err != nil {
		t.Fatalf("failed to write dirA track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dirB, "song.mp3"), content, 0644); err != nil {
		t.Fatalf("failed to write dirB track: %v", err)
	}

	events := collectScanEvents(t, []string{dirA, dirB})
	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 1 {
		t.Fatalf("expected 1 discovered event after dedup, got %d", len(discovered))
	}
	if len(enriched) != 1 {
		t.Fatalf("expected 1 enriched event after dedup, got %d", len(enriched))
	}
}

func TestScanForMediaProcessesDistinctFilesWithSameName(t *testing.T) {
	root := t.TempDir()
	dirA := filepath.Join(root, "a")
	dirB := filepath.Join(root, "b")

	if err := os.MkdirAll(dirA, 0755); err != nil {
		t.Fatalf("failed to create dirA: %v", err)
	}
	if err := os.MkdirAll(dirB, 0755); err != nil {
		t.Fatalf("failed to create dirB: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dirA, "song.mp3"), []byte("short"), 0644); err != nil {
		t.Fatalf("failed to write dirA track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dirB, "song.mp3"), []byte("longer-content"), 0644); err != nil {
		t.Fatalf("failed to write dirB track: %v", err)
	}

	events := collectScanEvents(t, []string{dirA, dirB})
	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 2 {
		t.Fatalf("expected 2 discovered events for distinct files, got %d", len(discovered))
	}
	if len(enriched) != 2 {
		t.Fatalf("expected 2 enriched events for distinct files, got %d", len(enriched))
	}
}

func TestScanForMediaSkipsDuplicateContentWithDifferentNames(t *testing.T) {
	root := t.TempDir()
	dirA := filepath.Join(root, "a")
	dirB := filepath.Join(root, "b")

	if err := os.MkdirAll(dirA, 0755); err != nil {
		t.Fatalf("failed to create dirA: %v", err)
	}
	if err := os.MkdirAll(dirB, 0755); err != nil {
		t.Fatalf("failed to create dirB: %v", err)
	}

	content := []byte("same-audio-content")
	if err := os.WriteFile(filepath.Join(dirA, "song.mp3"), content, 0644); err != nil {
		t.Fatalf("failed to write dirA track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dirB, "song (1).mp3"), content, 0644); err != nil {
		t.Fatalf("failed to write dirB track: %v", err)
	}

	events := collectScanEvents(t, []string{dirA, dirB})
	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 1 {
		t.Fatalf("expected 1 discovered event after content dedup, got %d", len(discovered))
	}
	if len(enriched) != 1 {
		t.Fatalf("expected 1 enriched event after content dedup, got %d", len(enriched))
	}
}

func TestScanForMediaWithSeedSkipsDuplicateFromPreviouslyLoadedTrack(t *testing.T) {
	root := t.TempDir()
	mediaDir := filepath.Join(root, "Media")
	otherDir := filepath.Join(root, "Other")

	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("failed to create mediaDir: %v", err)
	}
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatalf("failed to create otherDir: %v", err)
	}

	content := []byte("same-content")
	seedPath := filepath.Join(mediaDir, "song.mp3")
	dupPath := filepath.Join(otherDir, "song.mp3")

	if err := os.WriteFile(seedPath, content, 0644); err != nil {
		t.Fatalf("failed to write seed track: %v", err)
	}
	if err := os.WriteFile(dupPath, content, 0644); err != nil {
		t.Fatalf("failed to write duplicate track: %v", err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "unique.mp3"), []byte("unique"), 0644); err != nil {
		t.Fatalf("failed to write unique track: %v", err)
	}

	ch := make(chan ScanEvent)
	go ScanForMediaWithSeed(context.Background(), []string{otherDir}, []Track{BuildDiscoveredTrack(seedPath)}, ch)

	var events []ScanEvent
	for evt := range ch {
		events = append(events, evt)
	}

	discovered := trackEventsByType(events, ScanEventDiscovered)
	enriched := trackEventsByType(events, ScanEventEnriched)

	if len(discovered) != 1 {
		t.Fatalf("expected only the unique discovered event after seeding, got %d", len(discovered))
	}
	if len(enriched) != 1 {
		t.Fatalf("expected only the unique enriched event after seeding, got %d", len(enriched))
	}
	if discovered[0].Track == nil || discovered[0].Track.Filename != "unique.mp3" {
		t.Fatalf("expected unique.mp3 to remain after seeded dedup, got %+v", discovered[0].Track)
	}
}
