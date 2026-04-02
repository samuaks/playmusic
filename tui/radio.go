package tui

import (
	"fmt"
	"playmusic/library"
	"playmusic/player"
	"playmusic/search"
	"playmusic/yt_dlp"
	"strings"
	"time"
	"unicode"

	. "playmusic/helpers"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusMode int

const (
	focusPlayer focusMode = iota
	focusSearch
)

type OnlineModel struct {
	tracks          []library.Track
	current         int
	elapsed         time.Duration
	paused          bool
	searcher        *search.Searcher
	player          *player.Player
	searchQuery     string
	searching       bool
	result          *library.Track
	spinner         spinner.Model
	width           int
	height          int
	trackDownloaded bool
	loading         bool
	focus           focusMode
}

func NewOnlineModel(tracks []library.Track, searcher *search.Searcher) OnlineModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	return OnlineModel{
		searcher: searcher,
		player:   &player.Player{},
		spinner:  s,
		focus:    focusSearch,
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
		case "q", "?":
			if m.focus == focusSearch {
				m.searchQuery += msg.String()
				return m, nil
			}
			if m.focus == focusPlayer {
				m.focus = focusSearch
				return m, nil
			}
		case "backspace":
			if m.focus == focusSearch && len(m.searchQuery) > 0 {
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				return m, nil
			}
		case "esc":
			if m.focus == focusPlayer {
				m.searchQuery = ""
				m.focus = focusSearch
				return m, nil
			}
			if m.focus == focusSearch {
				m.focus = focusPlayer
				return m, nil
			}
		case "enter":
			if m.focus == focusSearch && m.searchQuery != "" {
				if m.result != nil {
					m.player.Stop()
				}
				m.focus = focusPlayer
				return m, debounceSearch(m.searchQuery)
			}
		case "ctrl+p", " ":
			if m.focus == focusSearch {
				m.searchQuery += msg.String()
				return m, nil
			}
			if m.paused {
				m.player.Resume()
				m.paused = false
			} else {
				m.player.Pause()
				m.paused = true
			}
			return m, nil
		case "ctrl+n", "right":
			if m.focus == focusPlayer {
				m.player.Next()
				return m, nil
			}
		case "ctrl+b", "left":
			if m.focus == focusPlayer {
				m.player.Prev()
				return m, nil
			}
		case "ctrl+d":
			return m, m.downloadTrack()
		default:
			if m.focus == focusSearch && len(msg.String()) > 0 {
				if len(msg.String()) == 1 {
					r := rune(msg.String()[0])
					if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || unicode.IsPunct(r) {
						m.searchQuery += msg.String()
						return m, nil
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
		m.elapsed = 0

		if len(m.tracks) > 0 {
			m.result = &m.tracks[0]
			cmd := m.playCurrent()

			if len(m.tracks) < 10 {
				m.loading = true
				cmd = tea.Batch(cmd, m.addTracksToRadioPlaylist())
			}

			return m, cmd
		}

		return m, nil

	case trackDoneMsg:
		m.current++
		m.elapsed = 0

		if m.current < len(m.tracks) {
			m.result = &m.tracks[m.current]
			// return m, m.playCurrent()

			var cmd tea.Cmd = m.playCurrent()

			if m.current >= len(m.tracks)-3 && !m.loading {
				m.loading = true
				cmd = tea.Batch(cmd, m.addTracksToRadioPlaylist())
			}

			return m, cmd

		}
		return m, nil

	case trackNextMsg:
		m.current++
		m.elapsed = 0

		var cmd tea.Cmd = m.playCurrent()

		if m.current >= len(m.tracks) {
			m.current = len(m.tracks) - 1
		}

		if m.current >= len(m.tracks)-3 && !m.loading {
			m.loading = true
			cmd = tea.Batch(cmd, m.addTracksToRadioPlaylist())
		}

		m.result = &m.tracks[m.current]

		return m, cmd

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

	case trackDownloadedMsg:
		m.trackDownloaded = true
		return m, clearNotificationAfter(3 * time.Second)

	case clearNotificationMsg:
		m.trackDownloaded = false
		return m, nil

	case timeToAddTracksMsq:
		m.loading = false
		m.tracks = append(m.tracks, msg.tracks...)
	}

	return m, nil
}

func (m *OnlineModel) runOnlineSearch(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.searcher.Search(query)
		if err != nil || len(tracks) == 0 {
			return searchDoneMsg{}
		}
		return searchResultsMsg{tracks: tracks}
	}
}

func (m OnlineModel) downloadTrack() tea.Cmd {
	return func() tea.Msg {
		_ = yt_dlp.DownloadAudio(m.tracks[m.current].YTVideoURL)
		return trackDownloadedMsg{}
	}
}

func (m OnlineModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Padding(0, 2).Render("Music Player — Online Radio") + "\n")

	var query string
	switch m.focus {
	case focusSearch:
		query = currentStyle.Render("> type artist or genre query to play music...")
		if m.searchQuery != "" {
			query = currentStyle.Render("> " + m.searchQuery)
		}
	case focusPlayer:
		if m.searchQuery == "" {
			query = dimmedStyle.Render("> type artist or genre query to play music...")
		}
		query = dimmedStyle.Render("> " + m.searchQuery)
	}

	if m.searching {
		query = m.spinner.View() + dimmedStyle.Render(m.searchQuery)
	}

	sb.WriteString(lipgloss.NewStyle().Padding(0, 2).Render(query) + "\n\n")

	if m.trackDownloaded {
		sb.WriteString(dimmedStyle.Padding(0, 2).Render("Track downloaded") + "\n")
	}
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
		dimmedStyle.Render("• enter: search • esc: clear query • q / ?: to search\n" +
			"• space: pause/play • -->: next track • <--: previous track • ctrl+d: download the track • ctrl+q: quit"),
	))

	return sb.String()
}

// not getting why % progress bar is not showing
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

func (m *OnlineModel) addTracksToRadioPlaylist() tea.Cmd {
	return func() tea.Msg {
		recTracks, err := yt_dlp.GetRecomendedWithYTVideoURL(m.tracks[m.current].YTVideoURL)
		if err != nil || len(recTracks) == 0 {
			return nil //maybe should do smth in this case
		}

		dedupTracks := m.deduplicateFoundTracks(recTracks)

		var tracks []library.Track
		for _, info := range dedupTracks {
			tracks = append(tracks, library.Track{
				Trackname:  info.Trackname,
				YTVideoURL: info.YTVideoURL,
				Duration:   info.Duration,
			})
		}

		m.loading = false

		return timeToAddTracksMsq{tracks: tracks}
	}
}

func (m *OnlineModel) deduplicateFoundTracks(tracks []yt_dlp.TrackInfo) []yt_dlp.TrackInfo {
	seen := make(map[string]struct{})
	var filtered []yt_dlp.TrackInfo

	//filling with existing tracks
	//should be global?
	for _, track := range m.tracks {
		seen[track.YTVideoURL] = struct{}{}
	}

	for _, foundtrack := range tracks {
		//apparently map has value and bool field(key check) if value exists - returns true in bool
		if _, ok := seen[foundtrack.YTVideoURL]; !ok {
			seen[foundtrack.YTVideoURL] = struct{}{}
			filtered = append(filtered, foundtrack)
		}
	}

	return filtered
}
