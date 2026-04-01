package tui

import (
	"fmt"
	"playmusic/library"
	"playmusic/player"
	"playmusic/search"
	"strings"
	"time"
	"unicode"

	. "playmusic/helpers"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OnlineModel struct {
	tracks      []library.Track
	current     int
	elapsed     time.Duration
	paused      bool
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
	return tick()
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
		case "ctrl+p":
			if m.paused {
				m.player.Resume()
				m.paused = false
			} else {
				m.player.Pause()
				m.paused = true
			}
			return m, nil
		case "ctrl+n":
			m.player.Next()
			return m, nil
		case "ctrl+b":
			m.player.Prev()
			return m, nil

		default:
			if len(msg.String()) > 0 {
				runes := []rune(msg.String())

				if len(runes) == 1 {
					r := runes[0]

					if unicode.IsGraphic(r) {
						m.searchQuery += msg.String()
						m.updateListItems()
						return m, debounceSearch(m.searchQuery)
					}
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

	case searchResultsMsg:
		m.searching = false
		m.tracks = msg.tracks
		m.current = 0

		if len(m.tracks) > 0 {
			m.result = &m.tracks[0]
			return m, m.playCurrent()
		}

		return m, nil

	case trackDoneMsg:
		m.current++
		m.elapsed = 0

		if m.current < len(m.tracks) {
			m.result = &m.tracks[m.current]
			return m, m.playCurrent()
		}
		return m, nil

	case trackNextMsg:
		m.current++
		m.elapsed = 0

		if m.current > len(m.tracks) {
			m.current = len(m.tracks) - 1
		}
		m.result = &m.tracks[m.current]
		return m, m.playCurrent()

	case trackPrevMsg:
		m.elapsed = 0

		if m.elapsed > 3*time.Second {
			return m, m.playCurrent()
		}

		m.current--
		if m.current < 0 {
			m.current = 0
		}
		m.result = &m.tracks[m.current]
		return m, m.playCurrent()

	case tickMsg:
		if !m.paused && m.player.IsPlaying() {
			m.elapsed += 500 * time.Millisecond
		}
		return m, tick()

	}

	return m, nil
}

func (m OnlineModel) runOnlineSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.searcher.Search(query)
		if err != nil || len(tracks) == 0 {
			return searchDoneMsg{}
		}
		return searchResultsMsg{tracks: tracks}
	}
}

func (m OnlineModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Padding(0, 2).Render("Music Player — Online") + "\n")

	query := "> " + m.searchQuery
	if m.searchQuery == "" {
		query = dimmedStyle.Render("> type artist or genre query to play music...")
	}
	if m.searching {
		query = m.spinner.View() + " " + dimmedStyle.Render(m.searchQuery)
	}
	sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(query) + "\n\n")

	_, elapsedStr := m.radioTrackProgress()

	if m.result != nil {
		symb := "▶ "

		if m.paused {
			symb = "▮▮ "
		}
		sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(
			currentStyle.Render(symb+m.result.Trackname) + "\n" +
				dimmedStyle.Render(m.result.YTVideoURL) + "\n" +
				dimmedStyle.Render(elapsedStr),
		))
	}

	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(
		dimmedStyle.Render("• enter: search • esc: clear query \n" +
			"• ctrl+p: pause/play • ctrl+n: next track • ctrl+b: previous track • ctrl+d: download the track • ctrl+q: quit"),
	))

	return sb.String()
}

// not getting why progress bar is not updating when
func (m OnlineModel) radioTrackProgress() (percent float64, elapsed string) {
	if m.current == -1 || len(m.tracks) == 0 {
		return 0, "0:00 / 0:00"
	}

	track := m.tracks[m.current]

	if track.Duration > 0 {
		percent = float64(m.elapsed) / float64(track.Duration)

		if percent > 1 {
			percent = 1
		}
	}

	elapsed = fmt.Sprintf("%s / %s", FormattedDuration(m.elapsed), track.FormatDuration())
	return
}

func (m OnlineModel) playCurrent() tea.Cmd {
	if len(m.tracks) == 0 || m.current >= len(m.tracks) {
		return nil
	}

	track := m.tracks[m.current]
	pl := m.player

	return func() tea.Msg {
		err := pl.PlayFromSearch(track.YTVideoURL)
		if err != nil {
			return searchDoneMsg{}
		}

		switch pl.Wait() {
		case player.TrackFinished:
			return trackDoneMsg{}
		case player.TrackNext:
			return trackNextMsg{}
		case player.TrackPrev:
			return trackPrevMsg{}
		}

		return nil
	}
}

func waitTrackDone(p *player.Player) tea.Cmd {
	return func() tea.Msg {
		if p.Done() == nil {
			return nil
		}
		<-p.Done()
		return trackDoneMsg{}
	}
}
