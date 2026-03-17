package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildAndEnrichTrackPreserveIdentity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	entry, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	if len(entry) != 1 {
		t.Fatalf("expected 1 dir entry, got %d", len(entry))
	}

	discovered := BuildDiscoveredTrack(path)
	enriched, _ := EnrichTrack(discovered)

	if discovered.Path != path {
		t.Fatalf("expected discovered path %q, got %q", path, discovered.Path)
	}
	if discovered.Filename != "song.mp3" {
		t.Fatalf("expected discovered filename song.mp3, got %q", discovered.Filename)
	}
	if discovered.Duration != 0 {
		t.Fatalf("expected discovered track to have zero duration, got %v", discovered.Duration)
	}
	if enriched.Path != discovered.Path {
		t.Fatalf("expected enriched path %q, got %q", discovered.Path, enriched.Path)
	}
	if enriched.Filename != discovered.Filename {
		t.Fatalf(
			"expected enriched filename %q, got %q",
			discovered.Filename,
			enriched.Filename,
		)
	}
	if enriched.Trackname == "" {
		t.Fatal("expected enriched trackname to stay populated")
	}
}
