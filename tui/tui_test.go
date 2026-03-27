package tui

import (
	"errors"
	"strings"
	"testing"
	"time"

	"playmusic/library"
	"playmusic/search"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdateAddsDiscoveredLibraryTrack(t *testing.T) {
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

func TestModelUpdateSkipsDuplicateDiscoveredLibraryTrack(t *testing.T) {
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

func TestModelUpdateReplacesTrackOnLibraryTrackUpdatedMsg(t *testing.T) {
	initial := []library.Track{
		{
			Trackname: "song.mp3",
			Path:      "/music/song.mp3",
			Filename:  "song.mp3",
		},
	}

	scanCh := make(chan library.ScanEvent)
	model := NewModel(initial, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackUpdatedMsg{
		track: library.Track{
			Trackname: "Artist - Title",
			Path:      "/music/song.mp3",
			Filename:  "song.mp3",
			Artist:    "Artist",
			Title:     "Title",
			Duration:  3 * time.Minute,
		},
	})

	got := updatedModel.(Model)

	if len(got.tracks) != 1 {
		t.Fatalf("expected updated track to replace existing one, got %d tracks", len(got.tracks))
	}
	if got.tracks[0].Artist != "Artist" {
		t.Fatalf("expected artist to be updated, got %q", got.tracks[0].Artist)
	}
	if got.tracks[0].Duration != 3*time.Minute {
		t.Fatalf("expected duration to be updated, got %v", got.tracks[0].Duration)
	}
	if cmd == nil {
		t.Fatal("expected Update to keep waiting for the next library event")
	}
}

func TestModelUpdateAppendsUpdatedTrackWhenOriginalMissing(t *testing.T) {
	scanCh := make(chan library.ScanEvent)
	model := NewModel(nil, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackUpdatedMsg{
		track: library.Track{
			Trackname: "Artist - Title",
			Path:      "/music/song.mp3",
			Filename:  "song.mp3",
			Artist:    "Artist",
			Title:     "Title",
		},
	})

	got := updatedModel.(Model)

	if len(got.tracks) != 1 {
		t.Fatalf("expected missing updated track to be appended, got %d tracks", len(got.tracks))
	}
	if got.tracks[0].Path != "/music/song.mp3" {
		t.Fatalf("expected appended updated track path, got %q", got.tracks[0].Path)
	}
	if got.scanAdded != 1 {
		t.Fatalf("expected scanAdded to increment for appended update, got %d", got.scanAdded)
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

	got := updatedModel.(Model)

	if got.scanError == nil {
		t.Fatal("expected scan error to be stored in model")
	}

	if cmd == nil {
		t.Fatal("expected Update to keep waiting after a scan error")
	}
}

func TestModelUpdateStopsWaitingWhenLibraryScanDone(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))

	updatedModel, cmd := model.Update(libraryScanDoneMsg{})

	got := updatedModel.(Model)

	if got.scanning {
		t.Fatal("expected scanning to be disabled after completion")
	}
	if !got.scanDone {
		t.Fatal("expected scanDone to be set after completion")
	}

	if cmd != nil {
		t.Fatal("expected no follow-up command after scan completion")
	}
}

func TestModelUpdateTracksAddedFromBackgroundScan(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))

	updatedModel, _ := model.Update(libraryTrackFoundMsg{
		track: library.Track{
			Trackname: "New Track",
			Path:      "/music/new.mp3",
			Filename:  "new.mp3",
		},
	})

	got := updatedModel.(Model)
	if got.scanAdded != 1 {
		t.Fatalf("expected scanAdded to increment, got %d", got.scanAdded)
	}
}

func TestLibraryScanStatusViewShowsActiveState(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))
	model.scanning = true
	model.scanAdded = 2

	view := model.libraryScanStatusView()
	if !strings.Contains(view, "Scanning library") {
		t.Fatalf("expected active scan status, got %q", view)
	}
}

func TestLibraryScanStatusViewShowsDoneState(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), nil)
	model.scanDone = true
	model.scanAdded = 3

	view := model.libraryScanStatusView()
	if !strings.Contains(view, "Library scan complete") {
		t.Fatalf("expected done scan status, got %q", view)
	}
}

func TestLibraryScanStatusViewShowsWarningState(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))
	model.scanning = true
	model.scanError = errors.New("scan failed")

	view := model.libraryScanStatusView()
	if !strings.Contains(view, "warning") {
		t.Fatalf("expected warning scan status, got %q", view)
	}
}

func TestWaitForLibraryEventReturnsNilForNilChannel(t *testing.T) {
	cmd := waitForLibraryEvent(nil)

	if cmd != nil {
		t.Fatal("expected nil command for nil scan channel")
	}
}

func TestWaitForLibraryEventReturnsTrackFoundMsgForDiscoveredEvent(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	track := library.Track{
		Trackname: "Found",
		Path:      "/music/found.mp3",
		Filename:  "found.mp3",
	}
	ch <- library.ScanEvent{
		Type:  library.ScanEventDiscovered,
		Track: &track,
	}

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

func TestWaitForLibraryEventReturnsTrackUpdatedMsgForEnrichedEvent(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	track := library.Track{
		Trackname: "Artist - Found",
		Path:      "/music/found.mp3",
		Filename:  "found.mp3",
		Artist:    "Artist",
	}
	ch <- library.ScanEvent{
		Type:  library.ScanEventEnriched,
		Track: &track,
	}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	updated, ok := msg.(libraryTrackUpdatedMsg)
	if !ok {
		t.Fatalf("expected libraryTrackUpdatedMsg, got %T", msg)
	}
	if updated.track.Path != track.Path {
		t.Fatalf("expected track path %q, got %q", track.Path, updated.track.Path)
	}
	if updated.track.Artist != "Artist" {
		t.Fatalf("expected updated artist %q, got %q", "Artist", updated.track.Artist)
	}
}

func TestModelUpdateKeepsWaitingAfterLibraryTrackFoundMsg(t *testing.T) {
	scanCh := make(chan library.ScanEvent)
	model := NewModel(nil, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackFoundMsg{
		track: library.Track{
			Trackname: "Found",
			Path:      "/music/found.mp3",
			Filename:  "found.mp3",
		},
	})

	got := updatedModel.(Model)
	if len(got.tracks) != 1 {
		t.Fatalf("expected 1 track after discovered event, got %d", len(got.tracks))
	}
	if cmd == nil {
		t.Fatal("expected Update to keep waiting for the next background scan event")
	}
}

func TestModelUpdateKeepsWaitingAfterLibraryTrackUpdatedMsg(t *testing.T) {
	scanCh := make(chan library.ScanEvent)
	model := NewModel([]library.Track{{
		Trackname: "song.mp3",
		Path:      "/music/song.mp3",
		Filename:  "song.mp3",
	}}, search.New(search.MockSource{}), scanCh)

	updatedModel, cmd := model.Update(libraryTrackUpdatedMsg{
		track: library.Track{
			Trackname: "Artist - Title",
			Path:      "/music/song.mp3",
			Filename:  "song.mp3",
			Artist:    "Artist",
		},
	})

	got := updatedModel.(Model)
	if got.tracks[0].Artist != "Artist" {
		t.Fatalf("expected updated artist %q, got %q", "Artist", got.tracks[0].Artist)
	}
	if cmd == nil {
		t.Fatal("expected Update to keep waiting for the next background scan event")
	}
}

func TestWaitForLibraryEventPrefersTrackUpdateOverErrorWhenTrackPresent(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	track := library.Track{
		Trackname: "Artist - Found",
		Path:      "/music/found.mp3",
		Filename:  "found.mp3",
		Artist:    "Artist",
	}
	ch <- library.ScanEvent{
		Type:  library.ScanEventEnriched,
		Track: &track,
		Err:   errors.New("ffprobe failed"),
	}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	updated, ok := msg.(libraryTrackUpdatedMsg)
	if !ok {
		t.Fatalf("expected libraryTrackUpdatedMsg, got %T", msg)
	}
	if updated.track.Path != track.Path {
		t.Fatalf("expected track path %q, got %q", track.Path, updated.track.Path)
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

func TestWaitForLibraryEventReturnsErrorForMalformedEvent(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	ch <- library.ScanEvent{}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	scanErr, ok := msg.(libraryScanErrorMsg)
	if !ok {
		t.Fatalf("expected libraryScanErrorMsg, got %T", msg)
	}
	if scanErr.err == nil || !strings.Contains(scanErr.err.Error(), "missing track and error") {
		t.Fatalf("expected malformed event error, got %v", scanErr.err)
	}
}

func TestWaitForLibraryEventReturnsErrorForUnknownEventType(t *testing.T) {
	ch := make(chan library.ScanEvent, 1)
	track := library.Track{
		Trackname: "Found",
		Path:      "/music/found.mp3",
		Filename:  "found.mp3",
	}
	ch <- library.ScanEvent{
		Type:  library.ScanEventType(99),
		Track: &track,
	}

	cmd := waitForLibraryEvent(ch)
	if cmd == nil {
		t.Fatal("expected command for scan channel")
	}

	msg := cmd()
	scanErr, ok := msg.(libraryScanErrorMsg)
	if !ok {
		t.Fatalf("expected libraryScanErrorMsg, got %T", msg)
	}
	if scanErr.err == nil || !strings.Contains(scanErr.err.Error(), "unknown library scan event type") {
		t.Fatalf("expected unknown event type error, got %v", scanErr.err)
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
}

func TestNewModelSetsScanningStateFromChannel(t *testing.T) {
	modelWithScan := NewModel(nil, search.New(search.MockSource{}), make(chan library.ScanEvent))
	if !modelWithScan.scanning {
		t.Fatal("expected model to start in scanning state when scan channel is provided")
	}

	modelWithoutScan := NewModel(nil, search.New(search.MockSource{}), nil)
	if modelWithoutScan.scanning {
		t.Fatal("expected model to start with scanning disabled when no scan channel is provided")
	}
}

func TestNewModelStartsInListFocus(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), nil)

	if model.focus != focusList {
		t.Fatalf("expected default focus to be focusList, got %v", model.focus)
	}
}

func TestModelUpdateEntersSearchFocusFromListOnQOrQuestion(t *testing.T) {
	cases := []string{"q", "?"}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			model := NewModel(nil, search.New(search.MockSource{}), nil)
			model.focus = focusList

			updatedModel, cmd := model.Update(tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(key),
			})
			got := updatedModel.(Model)

			if got.focus != focusSearch {
				t.Fatalf("expected focusSearch after %q, got %v", key, got.focus)
			}
			if cmd != nil {
				t.Fatalf("expected nil cmd on focus switch, got %v", cmd)
			}
		})
	}
}

func TestModelUpdateLeavesSearchFocusOnEsc(t *testing.T) {
	model := NewModel(nil, search.New(search.MockSource{}), nil)
	model.focus = focusSearch
	model.searchQuery = "abc"

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got := updatedModel.(Model)

	if got.focus != focusList {
		t.Fatalf("expected focusList after esc, got %v", got.focus)
	}
	if got.searchQuery != "" {
		t.Fatalf("expected search query to be cleared on esc, got %q", got.searchQuery)
	}
	if cmd != nil {
		t.Fatalf("expected nil cmd on esc focus switch, got %v", cmd)
	}
}

func TestModelUpdateEnterInListFocusStartsPlayback(t *testing.T) {
	model := NewModel([]library.Track{
		{
			Trackname: "Local Track",
			Path:      "/music/local.mp3",
			Filename:  "local.mp3",
		},
	}, search.New(search.MockSource{}), nil)
	model.focus = focusList

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected playback cmd on enter in list focus")
	}
	if got.current != 0 {
		t.Fatalf("expected current index 0 after enter in list focus, got %d", got.current)
	}
	if got.searching {
		t.Fatal("did not expect searching state when pressing enter in list focus")
	}
}

func TestModelUpdateEnterInSearchFocusRunsSearch(t *testing.T) {
	model := NewModel([]library.Track{
		{
			Trackname: "Local Track",
			Path:      "/music/local.mp3",
			Filename:  "local.mp3",
		},
	}, search.New(search.MockSource{}), nil)
	model.focus = focusSearch
	model.searchQuery = "beatles"

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected search cmd on enter in search focus with non-empty query")
	}
	if !got.searching {
		t.Fatal("expected searching state to be enabled on enter in search focus")
	}
	if got.current != -1 {
		t.Fatalf("did not expect playback selection change in search focus, got current=%d", got.current)
	}
}

func TestModelUpdateEnterInSearchFocusWithEmptyQueryDoesNothing(t *testing.T) {
	model := NewModel([]library.Track{
		{
			Trackname: "Local Track",
			Path:      "/music/local.mp3",
			Filename:  "local.mp3",
		},
	}, search.New(search.MockSource{}), nil)
	model.focus = focusSearch
	model.searchQuery = ""

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updatedModel.(Model)

	if cmd != nil {
		t.Fatalf("expected nil cmd on enter in search focus with empty query, got %v", cmd)
	}
	if got.searching {
		t.Fatal("expected searching state to remain false with empty query")
	}
	if got.current != -1 {
		t.Fatalf("did not expect playback selection change, got current=%d", got.current)
	}
}
