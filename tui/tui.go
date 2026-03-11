package tui

import (
	"fmt"
	. "playmusic/helpers"
	. "playmusic/library"
	. "playmusic/player"
	"playmusic/search"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type trackDoneMsg struct{}
type trackItem struct {
	track Track
}

type searchDebounceMsg struct {
	query string
}

func debounceSearch(query string) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return searchDebounceMsg{query}
	})
}

func (t trackItem) Title() string { return t.track.Trackname } // u can use some of the additional metadata for the UI

func (t trackItem) Description() string { return t.track.FormatDuration() }
func (t trackItem) FilterValue() string { return t.track.Title }

type Model struct {
	tracks    []Track
	current   int
	elapsed   time.Duration
	paused    bool
	player    *Player
	list      list.Model
	progress  progress.Model
	width     int
	height    int
	searcher  *search.Searcher
	spinner   spinner.Model
	searching bool
}

func NewModel(tracks []Track, searcher *search.Searcher) Model {
	items := make([]list.Item, len(tracks))

	for i, t := range tracks {
		items[i] = trackItem{t}
	}

	newList := list.New(items, newDelegate(""), 0, 0)

	newList.Title = TITLE
	newList.SetShowStatusBar(false)
	newList.SetShowHelp(true)
	newList.SetFilteringEnabled(true)
	newList.Styles.Title = titleStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	return Model{
		// doesnt pre-select any track
		current:  -1,
		tracks:   tracks,
		player:   &Player{},
		list:     newList,
		progress: progress.New(progress.WithDefaultGradient()),
		searcher: searcher,
		spinner:  s,
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

func (m Model) selectedTrack() (Track, int, bool) {
	item, ok := m.list.SelectedItem().(trackItem)
	if !ok {
		return Track{}, 0, false
	}
	for i, t := range m.tracks {
		if t.Path == item.track.Path {
			return t, i, true
		}
	}
	return Track{}, 0, false
}

func (m Model) filterQuery() string {
	return m.list.FilterValue()
}

type newTrackMsg struct {
	track Track
}

func (m Model) runSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.searcher.Search(query)
		if err != nil || len(tracks) == 0 {
			return nil
		}
		for _, t := range tracks {
			return newTrackMsg{t}
		}
		return nil
	}
}

func (m Model) move(direction int) Model {
	if len(m.tracks) == 0 {
		return m
	}

	m.elapsed = 0
	m.paused = false
	if m.list.FilterState() == list.FilterApplied {
		if direction > 0 {
			m.list.CursorDown()
		} else {
			m.list.CursorUp()
		}
		if _, idx, ok := m.selectedTrack(); ok {
			m.current = idx
			m.list.SetDelegate(newDelegate(m.tracks[m.current].Path))

		}
	} else {
		m.current = (m.current + direction + len(m.tracks)) % len(m.tracks)
		m.list.SetDelegate(newDelegate(m.tracks[m.current].Path))
		m.list.Select(m.current)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	setTerminalTitle("Playing Music 🎶")
	//return tea.Batch(m.playCurrent(), tick())
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 6
		m.list.SetSize(msg.Width, msg.Height-playerBarHeight)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			m.list, cmd = m.list.Update(msg)
			return m, tea.Batch(cmd, debounceSearch(m.filterQuery()))
		}
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		case " ":
			if m.paused {
				m.player.Resume()
				m.paused = false
			} else {
				m.player.Pause()
				m.paused = true
			}
			return m, nil
		case "n", "right":
			m = m.move(1)
			m.player.Next()
			return m, m.playCurrent()

		case "p", "left":
			m = m.move(-1)
			m.player.Next()
			return m, m.playCurrent()

		case "enter":
			if _, idx, ok := m.selectedTrack(); ok && idx != m.current {
				m.elapsed = 0
				m.paused = false
				m.current = idx
				m.list.SetDelegate(newDelegate(m.tracks[m.current].Path))
				m.player.Next()
				return m, m.playCurrent()
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
		m.current = (m.current + 1) % len(m.tracks)
		m.list.SetDelegate(newDelegate(m.tracks[m.current].Path))
		m.list.Select(m.current)
		return m, m.playCurrent()

	case spinner.TickMsg:
		if m.searching {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case searchDebounceMsg:
		if msg.query == m.filterQuery() && msg.query != "" {
			m.searching = true
			m.list.SetSize(m.width, m.height-playerBarHeight-searchBarHeight)

			return m, tea.Batch(m.runSearch(msg.query), m.spinner.Tick)
		}
	case newTrackMsg:
		m.searching = false
		m.list.SetSize(m.width, m.height-playerBarHeight)
		for _, t := range m.tracks {
			if t.Path == msg.track.Path {
				return m, nil
			}
		}
		m.tracks = append(m.tracks, msg.track)
		cmd = m.list.InsertItem(len(m.tracks)-1, trackItem{msg.track})
		return m, cmd

	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	if len(m.tracks) == 0 {
		return "No tracks\n"
	}

	nowPlaying := ""
	elapsed := ""
	var percent float64

	if m.current == -1 {
		nowPlaying = dimmedStyle.Render("No track selected")
		elapsed = dimmedStyle.Render("0:00 / 0:00")
	} else {
		track := m.tracks[m.current]
		if track.Duration > 0 {
			percent = float64(m.elapsed) / float64(track.Duration)
			if percent > 1 {
				percent = 1
			}
		}

		status := "▶"
		if m.paused {
			status = "⏸"
		}

		nowPlaying = currentStyle.Render(fmt.Sprintf("%s %s", status, track.Title))
		elapsed = dimmedStyle.Render(fmt.Sprintf("%s / %s", FormattedDuration(m.elapsed), track.FormatDuration()))
	}
	searchBar := ""
	if m.searching {
		searchBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Padding(0, 2).
			Render(m.spinner.View() + " searching external sources...")
		searchBar += "\n"
	}
	help := dimmedStyle.Render("space pause/resume • enter play")
	progressBar := barStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", nowPlaying, elapsed, m.progress.ViewAs(percent), help))

	return m.list.View() + "\n" + searchBar + progressBar

}
