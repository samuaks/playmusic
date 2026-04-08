package tui

import (
	"fmt"
	"math/rand"
	"path/filepath"
	. "playmusic/library"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func setTerminalTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func isExternalPlaybackTrack(track Track) bool {
	return strings.EqualFold(filepath.Ext(track.Path), ".mp4")
}
func mediaTypeLabel(track Track) string {
	if isExternalPlaybackTrack(track) {
		return "[VIDEO]"
	}
	return "[AUDIO]"
}

func (m *Model) syncPlaybackModeForCurrentTrack() {
	if m.current < 0 || m.current >= len(m.tracks) {
		m.externalPlayback = false
		return
	}
	m.externalPlayback = isExternalPlaybackTrack(m.tracks[m.current])
}

func (m Model) filteredTracks() []Track {
	if m.searchQuery == "" {
		return m.tracks
	}
	query := strings.ToLower(m.searchQuery)

	var results []Track

	for _, t := range m.tracks {
		if strings.Contains(strings.ToLower(t.Trackname), query) {
			results = append(results, t)
		}
	}
	return results
}

func (m *Model) updateListItems() {
	tracks := m.filteredTracks()
	items := make([]list.Item, len(tracks))

	for i, t := range tracks {
		items[i] = trackItem{t}
	}
	m.list.SetItems(items)
	if m.current != -1 {
		m.list.SetDelegate(newDelegate(m.tracks[m.current].Path, m.searchQuery))
	} else {
		m.list.SetDelegate(newDelegate("", m.searchQuery))
	}
}

func (m *Model) sortTracks() {
	currentID := ""
	if m.current >= 0 && m.current < len(m.tracks) {
		currentID = m.tracks[m.current].Identifier()
	}

	sort.SliceStable(m.tracks, func(i, j int) bool {
		return m.tracks[i].Trackname < m.tracks[j].Trackname
	})

	if currentID == "" {
		return
	}

	for i, t := range m.tracks {
		if t.Identifier() == currentID {
			m.current = i
			return
		}
	}
	m.current = -1
}

func highlightMatch(trackName, query string) string {
	if query == "" {
		return trackName
	}

	lower := strings.ToLower(trackName)
	index := strings.Index(lower, strings.ToLower(query))
	if index == -1 {
		return trackName
	}

	before := trackName[:index]
	match := trackName[index : index+len(query)]
	after := trackName[index+len(query):]

	return before + "\033[4m" + match + "\033[24m" + after

}

func (m Model) findNext() int {
	filtered := m.filteredTracks()

	if len(filtered) == 0 {
		return m.current
	}

	if m.isRandom {
		randomTrack := filtered[rand.Intn(len(filtered))]

		for i, t := range m.tracks {
			if t.Path == randomTrack.Path {
				return i
			}
		}
		return m.current
	}

	currentIndex := -1

	for i, t := range filtered {
		if t.Path == m.tracks[m.current].Path {
			currentIndex = i
			break
		}
	}

	nextFiltered := filtered[(currentIndex+1)%len(filtered)]

	for i, t := range m.tracks {
		if t.Path == nextFiltered.Path {
			return i
		}
	}
	return m.current

}
