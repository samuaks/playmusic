package tui

import (
	"testing"

	"playmusic/library"
	"playmusic/search"
)

func TestModelUpdateAddsLibraryTrackFoundMsg(t *testing.T) {
	initial := []library.Track{
		{
			Trackname: "Existing",
			Path:      "/music/existing.mp3",
			Filename:  "existing.mp3",
		},
	}

	scanCh := make(chan library.Track)
	model := NewModel(initial, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackFoundMsg{
		track: library.Track{
			Trackname: "New Track",
			Path:      "/music/new.mp3",
			Filename:  "new.mp3",
		},
	})

	got := updatedModel.(Model)

	if len(got.tracks) != 2 {
		t.Fatalf("expected 2 tracks after update, got %d", len(got.tracks))
	}
	if got.tracks[1].Path != "/music/new.mp3" {
		t.Fatalf("expected new track path to be appended, got %q", got.tracks[1].Path)
	}
	if cmd == nil {
		t.Fatal("expected Update to return a command that waits for the next scanned track")
	}
}

func TestModelUpdateSkipsDuplicateLibraryTrackByPath(t *testing.T) {
	initial := []library.Track{
		{
			Trackname: "Existing",
			Path:      "/music/existing.mp3",
			Filename:  "existing.mp3",
		},
	}

	scanCh := make(chan library.Track)
	model := NewModel(initial, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackFoundMsg{
		track: library.Track{
			Trackname: "Existing Duplicate",
			Path:      "/music/existing.mp3",
			Filename:  "existing.mp3",
		},
	})

	got := updatedModel.(Model)

	if len(got.tracks) != 1 {
		t.Fatalf("expected duplicate path to be ignored, got %d tracks", len(got.tracks))
	}
	if cmd == nil {
		t.Fatal("expected Update to keep waiting for the next scanned track")
	}
}
