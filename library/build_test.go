package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDiscoveredTrackPopulatesBaseFields(t *testing.T) {
	path := filepath.Join("music", "song.mp3")

	track := BuildDiscoveredTrack(path)

	if track.Path != path {
		t.Fatalf("expected path %q, got %q", path, track.Path)
	}
	if track.Filename != "song.mp3" {
		t.Fatalf("expected filename %q, got %q", "song.mp3", track.Filename)
	}
	if track.Trackname == "" {
		t.Fatal("expected Trackname to be populated")
	}
	if track.Duration != 0 {
		t.Fatalf("expected zero duration for discovered track, got %v", track.Duration)
	}
	if track.Title != "" {
		t.Fatalf("expected empty title for discovered track, got %q", track.Title)
	}
	if track.Artist != "" {
		t.Fatalf("expected empty artist for discovered track, got %q", track.Artist)
	}
	if track.Album != "" {
		t.Fatalf("expected empty album for discovered track, got %q", track.Album)
	}
	if track.Genre != "" {
		t.Fatalf("expected empty genre for discovered track, got %q", track.Genre)
	}
	if track.Year != 0 {
		t.Fatalf("expected zero year for discovered track, got %d", track.Year)
	}
}

func TestEnrichTrackPreservesIdentityFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	discovered := BuildDiscoveredTrack(path)
	enriched, err := EnrichTrack(discovered)
	if err != nil {
		t.Logf("enrichment returned non-fatal error: %v", err)
	}

	if enriched.Path != discovered.Path {
		t.Fatalf("expected path %q, got %q", discovered.Path, enriched.Path)
	}
	if enriched.Filename != discovered.Filename {
		t.Fatalf("expected filename %q, got %q", discovered.Filename, enriched.Filename)
	}
	if enriched.Trackname == "" {
		t.Fatal("expected Trackname to remain populated after enrichment")
	}
}

func TestEnrichTrackKeepsTrackUsableWhenMetadataIsUnavailable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	discovered := BuildDiscoveredTrack(path)
	enriched, err := EnrichTrack(discovered)
	if err != nil {
		t.Logf("enrichment returned non-fatal error: %v", err)
	}

	if enriched.Path != path {
		t.Fatalf("expected path %q, got %q", path, enriched.Path)
	}
	if enriched.Filename != "song.mp3" {
		t.Fatalf("expected filename song.mp3, got %q", enriched.Filename)
	}
	if enriched.Trackname == "" {
		t.Fatal("expected Trackname to be populated")
	}
}

func TestEnrichTrackKeepsFallbackTracknameWhenMetadataMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "song.mp3")

	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write track: %v", err)
	}

	discovered := BuildDiscoveredTrack(path)
	enriched, err := EnrichTrack(discovered)
	if err != nil {
		t.Logf("enrichment returned non-fatal error: %v", err)
	}

	if discovered.Trackname == "" {
		t.Fatal("expected discovered track to have fallback name")
	}
	if enriched.Trackname == "" {
		t.Fatal("expected enriched track to retain a usable name")
	}
}
