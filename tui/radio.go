package tui

import (
	"playmusic/library"
	"playmusic/player"
	"playmusic/search"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OnlineModel struct {
	searcher    *search.Searcher
	player      *player.Player
	searchQuery string
	searching   bool
	result      *library.Track
	spinner     spinner.Model
	width       int
	height      int
}

func NewOnlineModel(tracks []library.Track, searcher *search.Searcher) OnlineModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	return OnlineModel{
		searcher: searcher,
		player:   &player.Player{},
		spinner:  s,
	}
}

func (m OnlineModel) Init() tea.Cmd {
	return nil
}

func (m OnlineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
			return m, nil
		case "esc":
			m.searchQuery = ""
			return m, nil
		case "enter":
			if m.result != nil {
				m.player.Stop()
				// playback handled separately
			}
			return m, debounceSearch(m.searchQuery)
		default:
			if len(msg.String()) == 1 {
				r := rune(msg.String()[0])
				if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || unicode.IsPunct(r) {
					m.searchQuery += msg.String()
					return m, debounceSearch(m.searchQuery)
				}
			}
		}

	case searchDebounceMsg:
		if msg.query == m.searchQuery && msg.query != "" {
			m.searching = true
			return m, tea.Batch(m.runOnlineSearch(msg.query), m.spinner.Tick)
		}

	case searchTrackFoundMsg:
		m.searching = false
		m.result = &msg.track
		return m, nil

	case searchDoneMsg:
		m.searching = false
		return m, nil

	case spinner.TickMsg:
		if m.searching {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m OnlineModel) runOnlineSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.searcher.Search(query)
		if err != nil || len(tracks) == 0 {
			return searchDoneMsg{}
		}
		return searchTrackFoundMsg{tracks[0]}
	}
}
func (m OnlineModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Padding(0, 2).Render("Music Player — Online") + "\n")

	query := "> " + m.searchQuery
	if m.searchQuery == "" {
		query = dimmedStyle.Render("> type to search youtube...")
	}
	if m.searching {
		query = m.spinner.View() + " " + dimmedStyle.Render(m.searchQuery)
	}
	sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(query) + "\n\n")

	if m.result != nil {
		sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(
			currentStyle.Render("▶ "+m.result.Trackname) + "\n" +
				dimmedStyle.Render(m.result.YTVideoURL),
		))
	}

	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(
		dimmedStyle.Render("enter search • esc clear • ctrl+q quit"),
	))

	return sb.String()
}
