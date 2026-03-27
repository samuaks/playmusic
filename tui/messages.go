package tui

import (
	"playmusic/library"
	"time"
)

type tickMsg time.Time
type trackDoneMsg struct {
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
