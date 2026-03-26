package tui

import (
	"playmusic/library"
	"time"
)

type tickMsg time.Time
type trackDoneMsg struct{}
type trackNextMsg struct{}
type trackPrevMsg struct{}

type searchDebounceMsg struct {
	query string
}

// searchTrackFoundMsg carries a track returned by the search subsystem.
type searchTrackFoundMsg struct {
	track library.Track
}

// libraryTrackFoundMsg delivers a track discovered by the background
// library scan into the Bubble Tea update loop.
type libraryTrackFoundMsg struct {
	track library.Track
}

// libraryTrackUpdatedMsg carries enrichment updates for an already
// discovered library track.
type libraryTrackUpdatedMsg struct {
	track library.Track
}

// libraryScanErrorMsg reports a non-fatal background scan error.
type libraryScanErrorMsg struct {
	err error
}

// libraryScanDoneMsg signals that the background library scan has finished.
type libraryScanDoneMsg struct{}

type searchDoneMsg struct{}

type searchResultsMsg struct {
	tracks []library.Track
}
