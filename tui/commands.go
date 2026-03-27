package tui

import (
	"fmt"
	"playmusic/library"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func debounceSearch(query string) tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return searchDebounceMsg{query}
	})
}

func clearNotificationAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearNotificationMsg{}
	})
}

// waitForLibraryEvent blocks until the next background-scanned track arrives.
// if the scan channel is closed, it emits a completion message instead
func waitForLibraryEvent(ch <-chan library.ScanEvent) tea.Cmd {
	if ch == nil {
		return nil
	}

	return func() tea.Msg {
		evt, ok := <-ch
		if !ok {
			return libraryScanDoneMsg{}
		}
		if evt.Track == nil {
			if evt.Err != nil {
				return libraryScanErrorMsg{err: evt.Err}
			}
			return libraryScanErrorMsg{err: fmt.Errorf("library scan event missing track and error")}
		}

		switch evt.Type {
		case library.ScanEventDiscovered:
			return libraryTrackFoundMsg{track: *evt.Track}
		case library.ScanEventEnriched:
			return libraryTrackUpdatedMsg{track: *evt.Track}
		default:
			return libraryScanErrorMsg{err: fmt.Errorf("unknown library scan event type: %d", evt.Type)}
		}
	}
}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) playCurrent() tea.Cmd {
	track := m.tracks[m.current]
	player := m.player

	return func() tea.Msg {
		err := player.Play(track.Path)

		if err != nil {
			fmt.Println("Error playing track:", err)
			return trackDoneMsg{}
		}
		if player.SimpleWait() {
			return trackDoneMsg{}
		}
		return nil
	}
}

/* func (m Model) runSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.searcher.Search(query)
		if err != nil || len(tracks) == 0 {
			return searchDoneMsg{}
		}
		for _, t := range tracks {
			return searchTrackFoundMsg{track: t}
		}
		return searchDoneMsg{}
	}
}
*/
