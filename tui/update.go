package tui

import (
	"fmt"
	"playmusic/library"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	setTerminalTitle(TITLE + " ??")

	return tea.Batch(
		tick(),
		waitForLibraryEvent(m.scanCh),
	)
}

func (m Model) selectedTrack() (library.Track, int, bool) {
	item, ok := m.list.SelectedItem().(trackItem)
	if !ok {
		return library.Track{}, 0, false
	}
	for i, t := range m.tracks {
		if t.Identifier() == item.track.Identifier() {
			return t, i, true
		}
	}
	return library.Track{}, 0, false
}

func (m Model) trackIndexByPath(path string) int {
	for i, t := range m.tracks {
		if t.Path == path {
			return i
		}
	}
	return -1
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 6
		m.list.SetSize(msg.Width, msg.Height-playerBarHeight-searchBarHeight-scanBarHeight)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		case "q", "?":
			if m.focus == focusList {
				m.focus = focusSearch
				return m, nil
			}
			if m.focus == focusSearch {
				m.searchQuery += msg.String()
				m.updateListItems()
				return m, nil
			}
		case "esc":
			if m.focus == focusSearch {
				m.focus = focusList
				m.searchQuery = ""
				m.updateListItems()
				return m, nil
			}
			return m, nil
		case " ":
			if m.paused {
				m.player.Resume()
				m.paused = false
			} else {
				m.player.Pause()
				m.paused = true
			}
			return m, nil
		case "ctrl+r":
			m.isRandom = !m.isRandom
			return m, nil
		case "enter":
			if m.focus == focusSearch {
				query := strings.TrimSpace(m.searchQuery)
				m.focus = focusList
				m.searchQuery = ""
				m.updateListItems()

				if query == "" {
					return m, nil
				}

				m.searchRequestID++
				reqID := m.searchRequestID
				m.searching = true
				m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight-scanBarHeight)
				return m, tea.Batch(m.runSearch(query, reqID), m.spinner.Tick)
			}
			if _, idx, ok := m.selectedTrack(); ok && idx != m.current {
				m.elapsed = 0
				m.paused = false
				m.current = idx
				m.list.SetDelegate(newDelegate(m.tracks[m.current].Identifier(), m.searchQuery))
				m.player.Next()
				return m, m.playCurrent()
			}
		case "backspace":
			if m.focus == focusSearch && len(m.searchQuery) > 0 {
				queryRunes := []rune(m.searchQuery)
				if len(queryRunes) > 0 {
					queryRunes = queryRunes[:len(queryRunes)-1]
					m.searchQuery = string(queryRunes)
					m.updateListItems()
					return m, nil
				}
			}
			return m, nil
		default:
			if m.focus == focusSearch && len(msg.String()) > 0 {
				runes := []rune(msg.String())
				if len(runes) == 1 {
					r := runes[0]
					if unicode.IsGraphic(r) && msg.String() != " " {
						m.searchQuery += msg.String()
						m.updateListItems()
						return m, nil
					}
				}
			}
		}
	case tickMsg:
		if !m.paused && m.player.IsPlaying() {
			m.elapsed += 500 * time.Millisecond
		}
		return m, tick()
	case trackDoneMsg:
		m.elapsed = 0
		m.paused = false
		m.current = m.findNext()
		m.list.SetDelegate(newDelegate(m.tracks[m.current].Identifier(), m.searchQuery))
		m.list.Select(m.current)
		return m, m.playCurrent()

	case spinner.TickMsg:
		if m.searching {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case searchTrackFoundMsg:
		if msg.reqID != m.searchRequestID {
			return m, nil
		}
		m.searching = false
		m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight-scanBarHeight)
		for _, t := range m.tracks {
			if t.Identifier() == msg.track.Identifier() {
				return m, nil
			}
		}
		m.tracks = append(m.tracks, msg.track)
		cmd = m.list.InsertItem(len(m.tracks)-1, trackItem{msg.track})
		return m, cmd
	case libraryTrackFoundMsg:
		// Deduplicate by path because startup tracks and background scan
		// may overlap or the scanner may revisit the same location.
		if m.trackIndexByPath(msg.track.Path) == -1 {
			m.tracks = append(m.tracks, msg.track)
			m.scanAdded++
			m.updateListItems()
		}
		return m, waitForLibraryEvent(m.scanCh)
	case libraryTrackUpdatedMsg:
		idx := m.trackIndexByPath(msg.track.Path)
		if idx >= 0 {
			m.tracks[idx] = msg.track
		} else {
			m.tracks = append(m.tracks, msg.track)
			m.scanAdded++
		}
		m.updateListItems()
		return m, waitForLibraryEvent(m.scanCh)
	case libraryScanErrorMsg:
		m.scanError = msg.err
		return m, waitForLibraryEvent(m.scanCh)
	case libraryScanDoneMsg:
		m.scanning = false
		m.scanDone = true
		return m, nil
	case searchDoneMsg:
		if msg.reqID != m.searchRequestID {
			return m, nil
		}
		m.searching = false
		m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight-scanBarHeight)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		m.searchBarView(),
		m.libraryScanStatusView(),
		m.list.View(),
		m.playerBarView())
}
