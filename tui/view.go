package tui

import (
	"fmt"
	. "playmusic/helpers"
	"playmusic/library"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) playerBarView() string {
	nowPlaying := ""
	percent, elapsedStr := m.trackProgress()
	elapsed := dimmedStyle.Render(elapsedStr)

	if m.current == -1 {
		nowPlaying = dimmedStyle.Render("No track selected")
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
		if m.isRandom {
			status += " 🔀"
		}

		nowPlaying = currentStyle.Render(fmt.Sprintf("%s %s", status, track.Trackname))
	}

	return barStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", nowPlaying, elapsed, m.progress.ViewAs(percent), m.helpView()))
}

func GetProgressPercentage(elapsed time.Duration, current int, tracks []library.Track) float64 {
	if current == -1 || len(tracks) == 0 {
		return 0
	}

	track := tracks[current]
	if track.Duration > 0 {
		percent := float64(elapsed) / float64(track.Duration)
		if percent > 1 {
			return 1
		}
		return percent
	}
	return 0
}

func (m Model) trackProgress() (percent float64, elapsed string) {
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

func (m Model) searchBarView() string {
	title := titleStyle.Padding(0, 2).Render(TITLE)
	var query string

	switch {
	case m.searchQuery == "":
		query = dimmedStyle.Render("> " + SEARCHBAR_TEXT)
	default:
		query = dimmedStyle.Render("> ") + dimmedStyle.Render(m.searchQuery)
	}
	input := lipgloss.NewStyle().Padding(0, 2).Render(query)
	return fmt.Sprintf("%s\n%s\n", title, input)
}

func (m Model) libraryScanStatusView() string {
	switch {
	case m.scanning && m.scanError != nil:
		return scanWarningStyle.Padding(0, 2).Render(fmt.Sprintf("Library scan warning: %v", m.scanError))
	case m.scanning:
		if m.scanAdded > 0 {
			return scanStatusStyle.Padding(0, 2).Render(fmt.Sprintf("Scanning library... +%d tracks", m.scanAdded))
		}
		return scanStatusStyle.Padding(0, 2).Render("Scanning library...")
	case m.scanDone && m.scanError != nil:
		return scanWarningStyle.Padding(0, 2).Render(fmt.Sprintf("Library scan complete with warning: %v", m.scanError))
	case m.scanDone:
		return scanStatusStyle.Padding(0, 2).Render(fmt.Sprintf("Library scan complete. Added %d tracks.", m.scanAdded))
	default:
		return ""
	}
}

func (m Model) helpView() string {
	return dimmedStyle.Render(HELP_TEXT)
}
