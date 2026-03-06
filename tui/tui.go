package tui

import (
	"fmt"
	. "playmusic/helpers"
	. "playmusic/library"
	. "playmusic/player"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time
type trackDoneMsg struct{}
type trackItem struct {
	track Track
}

func (t trackItem) Title() string       { return t.track.Title }
func (t trackItem) Description() string { return t.track.FormatDuration() }
func (t trackItem) FilterValue() string { return t.track.Title }

type Model struct {
	tracks   []Track
	current  int
	elapsed  time.Duration
	paused   bool
	player   *Player
	list     list.Model
	progress progress.Model
	width    int
	height   int
}

func NewModel(tracks []Track) Model {
	items := make([]list.Item, len(tracks))

	for i, t := range tracks {
		items[i] = trackItem{t}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedTitleStyle
	delegate.Styles.SelectedDesc = selectedDescStyle

	songs := list.New(items, delegate, 0, 0)
	songs.Title = "PlayMusic"
	songs.SetShowStatusBar(false)
	songs.SetShowHelp(false)
	songs.SetFilteringEnabled(false)
	songs.Styles.Title = titleStyle

	return Model{
		tracks:   tracks,
		player:   &Player{},
		list:     songs,
		progress: progress.New(progress.WithDefaultGradient()),
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

func (m Model) move(direction int) Model {
	if len(m.tracks) == 0 {
		return m
	}
	m.player.Stop()
	m.elapsed = 0
	m.paused = false
	m.current = (m.current + direction + len(m.tracks)) % len(m.tracks)
	m.list.Select(m.current)
	return m
}

func (m Model) Init() tea.Cmd {
	setTerminalTitle("Playing Music 🎶")
	return tea.Batch(m.playCurrent(), tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		playerBarHeight := 7
		m.list.SetSize(msg.Width, msg.Height-playerBarHeight)
		m.progress.Width = msg.Width - 6

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
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
		case "p", "left":
			m = m.move(-1)
			m.player.Next()
			return m, m.playCurrent()
		case "n", "right":
			m = m.move(1)
			m.player.Next()
			return m, m.playCurrent()

		case "enter":
			selected := m.list.Index()
			if selected != m.current {
				m.player.Stop()
				m.elapsed = 0
				m.paused = false
				m.current = selected
				m.list.Select(m.current)
				return m, m.playCurrent()
			}
			return m, nil

		}
	case tickMsg:
		if !m.paused {
			m.elapsed += 500 * time.Millisecond
		}
		return m, tick()
	case trackDoneMsg:
		m.elapsed = 0
		m.paused = false
		m.current = (m.current + 1) % len(m.tracks)
		m.list.Select(m.current)
		return m, m.playCurrent()

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	if len(m.tracks) == 0 {
		return "No tracks\n"
	}

	track := m.tracks[m.current]

	var percent float64
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

	playing := currentStyle.Render(fmt.Sprintf("%s %s", status, track.Title))
	elapsed := dimmedStyle.Render(fmt.Sprintf("%s / %s", FormattedDuration(m.elapsed), track.FormatDuration()))
	help := dimmedStyle.Render("space pause/resume • ←/→ prev/next • enter play • q quit")
	progressBar := barStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", playing, elapsed, m.progress.ViewAs(percent), help))

	return fmt.Sprintf("%s\n%s", m.list.View(), progressBar)

}
