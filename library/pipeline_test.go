package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessCandidateSeparatesDiscoveryFromEnrichment(t *testing.T) {
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

	candidate := processCandidate(
		newScanState(make(map[string]struct{}), make(map[string]struct{})),
		path,
		entry[0],
	)

	if !candidate.include {
		t.Fatal("expected candidate to be included")
	}
	if candidate.discovered.Path != path {
		t.Fatalf("expected discovered path %q, got %q", path, candidate.discovered.Path)
	}
	if candidate.discovered.Filename != "song.mp3" {
		t.Fatalf("expected discovered filename song.mp3, got %q", candidate.discovered.Filename)
	}
	if candidate.discovered.Duration != 0 {
		t.Fatalf("expected discovered track to have zero duration, got %v", candidate.discovered.Duration)
	}
	if candidate.enriched.Path != candidate.discovered.Path {
		t.Fatalf("expected enriched path %q, got %q", candidate.discovered.Path, candidate.enriched.Path)
	}
	if candidate.enriched.Filename != candidate.discovered.Filename {
		t.Fatalf(
			"expected enriched filename %q, got %q",
			candidate.discovered.Filename,
			candidate.enriched.Filename,
		)
	}
	if candidate.enriched.Trackname == "" {
		t.Fatal("expected enriched trackname to stay populated")
	}
}
