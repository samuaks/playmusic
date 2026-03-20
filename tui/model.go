package tui

import (
	"playmusic/library"
	. "playmusic/player"
	"playmusic/search"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type trackItem struct {
	track library.Track
}

func (t trackItem) Title() string { return t.track.Trackname } // u can use some of the additional metadata for the UI

func (t trackItem) Description() string { return t.track.FormatDuration() }
func (t trackItem) FilterValue() string { return t.track.Trackname }

type focusMode int

const (
	focusList focusMode = iota
	focusSearch
)

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
	focus       focusMode
	searcher    *search.Searcher
	spinner     spinner.Model
	searching   bool
	searchQuery string
	scanCh      <-chan library.ScanEvent
	scanning    bool
	scanDone    bool
	scanError   error
	scanAdded   int
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
		focus:    focusList,
		searcher: searcher,
		spinner:  s,
		scanCh:   scanCh,
		scanning: scanCh != nil,
	}
}
