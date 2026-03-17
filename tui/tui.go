package tui

import (
	"fmt"
	"playmusic/library"
	. "playmusic/player"
	"playmusic/search"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type trackDoneMsg struct{}
type trackItem struct {
	track library.Track
}

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

func debounceSearch(query string) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return searchDebounceMsg{query}
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
			return libraryScanDoneMsg{}
		}

		switch evt.Type {
		case library.ScanEventDiscovered:
			return libraryTrackFoundMsg{track: *evt.Track}
		case library.ScanEventEnriched:
			return libraryTrackUpdatedMsg{track: *evt.Track}
		default:
			if evt.Err != nil {
				return libraryScanErrorMsg{err: evt.Err}
			}
			return libraryScanDoneMsg{}
		}
	}
}

func (t trackItem) Title() string { return t.track.Trackname } // u can use some of the additional metadata for the UI

func (t trackItem) Description() string { return t.track.FormatDuration() }
func (t trackItem) FilterValue() string { return t.track.Trackname }

type Model struct {
	tracks      []library.Track
	current     int
	elapsed     time.Duration
	paused      bool
	player      *Player
	list        list.Model
	progress    progress.Model
	width       int
	height      int
	searcher    *search.Searcher
	spinner     spinner.Model
	searching   bool
	searchQuery string
	scanCh      <-chan library.ScanEvent
}

func NewModel(tracks []library.Track, searcher *search.Searcher, scanCh <-chan library.ScanEvent) Model {
	items := make([]list.Item, len(tracks))

	for i, t := range tracks {
		items[i] = trackItem{t}
	}

	newList := list.New(items, newDelegate("", ""), 0, 0)

	newList.SetShowStatusBar(false)
	newList.SetShowTitle(false)
	newList.SetShowHelp(false)

	newList.Styles.NoItems = emptyStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	return Model{
		current:  -1,
		tracks:   tracks,
		player:   &Player{},
		list:     newList,
		progress: progress.New(progress.WithDefaultGradient()),
		searcher: searcher,
		spinner:  s,
		scanCh:   scanCh,
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
		if err := player.Play(track.Path); err != nil {
			return trackDoneMsg{}
		}
		<-player.Done()
		return trackDoneMsg{}
	}
}

func (m Model) selectedTrack() (library.Track, int, bool) {
	item, ok := m.list.SelectedItem().(trackItem)
	if !ok {
		return library.Track{}, 0, false
	}
	for i, t := range m.tracks {
		if t.Path == item.track.Path {
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

func (m Model) runSearch(query string) tea.Cmd {
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

func (m Model) Init() tea.Cmd {
	setTerminalTitle(TITLE + " 🎶")

	return tea.Batch(
		tick(),
		waitForLibraryEvent(m.scanCh),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 6
		m.list.SetSize(msg.Width, msg.Height-playerBarHeight-searchBarHeight)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		case "esc":
			m.searchQuery = ""
			m.updateListItems()
			return m, debounceSearch("")
		case " ":
			if m.paused {
				m.player.Resume()
				m.paused = false
			} else {
				m.player.Pause()
				m.paused = true
			}
			return m, nil
		case "enter":
			if _, idx, ok := m.selectedTrack(); ok && idx != m.current {
				m.elapsed = 0
				m.paused = false
				m.current = idx
				m.list.SetDelegate(newDelegate(m.tracks[m.current].Path, m.searchQuery))
				m.player.Next()
				return m, m.playCurrent()
			}
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.updateListItems()
				return m, debounceSearch(m.searchQuery)
			}
			return m, nil
		default:
			if len(msg.String()) == 1 && msg.String() != " " {
				r := rune(msg.String()[0])
				if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) {
					m.searchQuery += msg.String()
					m.updateListItems()
					return m, debounceSearch(m.searchQuery)
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
		m.list.SetDelegate(newDelegate(m.tracks[m.current].Path, m.searchQuery))
		m.list.Select(m.current)
		return m, m.playCurrent()

	case spinner.TickMsg:
		if m.searching {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case searchDebounceMsg:
		if msg.query == m.searchQuery && msg.query != "" {
			m.searching = true
			m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight)

			return m, tea.Batch(m.runSearch(msg.query), m.spinner.Tick)
		}
	case searchTrackFoundMsg:
		m.searching = false
		m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight)
		for _, t := range m.tracks {
			if t.Path == msg.track.Path {
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
			m.updateListItems()
		}
		return m, waitForLibraryEvent(m.scanCh)
	case libraryTrackUpdatedMsg:
		idx := m.trackIndexByPath(msg.track.Path)
		if idx >= 0 {
			m.tracks[idx] = msg.track
		} else {
			m.tracks = append(m.tracks, msg.track)
		}
		m.updateListItems()
		return m, waitForLibraryEvent(m.scanCh)
	case libraryScanErrorMsg:
		return m, waitForLibraryEvent(m.scanCh)
	case libraryScanDoneMsg:
		return m, nil
	case searchDoneMsg:
		m.searching = false
		m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s",
		m.searchBarView(),
		m.list.View(),
		m.playerBarView())

}
