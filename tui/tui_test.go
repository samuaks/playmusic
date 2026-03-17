package tui

import (
	"errors"
	"testing"

	"playmusic/library"
	"playmusic/search"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdateAddsLibraryTrackFoundMsg(t *testing.T) {
	initial := []library.Track{
		{
			Trackname: "Existing",
			Path:      "/music/existing.mp3",
			Filename:  "existing.mp3",
		},
	}

	scanCh := make(chan library.ScanEvent)
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
		t.Fatal("expected Update to return a command that waits for the next library event")
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

	scanCh := make(chan library.ScanEvent)
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
		t.Fatal("expected Update to keep waiting for the next library event")
	}
}

func TestModelUpdateKeepsWaitingAfterLibraryScanError(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))

	updatedModel, cmd := model.Update(libraryScanErrorMsg{
		err: errors.New("scan failed"),
	})

	_ = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected Update to keep waiting after a scan error")
	}
}

func TestModelUpdateStopsWaitingWhenLibraryScanDone(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))

	updatedModel, cmd := model.Update(libraryScanDoneMsg{})

	_ = updatedModel.(Model)

	if cmd != nil {
		t.Fatal("expected no follow-up command after scan completion")
	}
}

func TestWaitForLibraryEventReturnsNilForNilChannel(t *testing.T) {
	cmd := waitForLibraryEvent(nil)

	if cmd != nil {
		t.Fatal("expected nil command for nil scan channel")
	}
}

func TestWaitForLibraryEventReturnsTrackFoundMsg(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	track := library.Track{
		Trackname: "Found",
		Path:      "/music/found.mp3",
		Filename:  "found.mp3",
	}
	ch <- library.ScanEvent{Track: &track}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	found, ok := msg.(libraryTrackFoundMsg)
	if !ok {
		t.Fatalf("expected libraryTrackFoundMsg, got %T", msg)
	}
	if found.track.Path != track.Path {
		t.Fatalf("expected track path %q, got %q", track.Path, found.track.Path)
	}
}

func TestWaitForLibraryEventReturnsScanErrorMsg(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	wantErr := errors.New("scan failed")
	ch <- library.ScanEvent{Err: wantErr}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	scanErr, ok := msg.(libraryScanErrorMsg)
	if !ok {
		t.Fatalf("expected libraryScanErrorMsg, got %T", msg)
	}
	if !errors.Is(scanErr.err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, scanErr.err)
	}
}

func TestWaitForLibraryEventReturnsScanDoneMsgWhenChannelClosed(t *testing.T) {
	ch := make(chan library.ScanEvent)
	close(ch)

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	if _, ok := msg.(libraryScanDoneMsg); !ok {
		t.Fatalf("expected libraryScanDoneMsg, got %T", msg)
	}
}

func TestModelInitReturnsStartupCmd(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))

	cmd := model.Init()

	if cmd == nil {
		t.Fatal("expected Init to return a startup command batch")
	}
	if _, ok := any(cmd).(tea.Cmd); !ok {
		t.Fatal("expected Init to return a tea.Cmd")
	}
}
