package tui

import (
	"fmt"
	. "playmusic/helpers"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) playerBarView() string {
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

		status := ">"
		if m.paused {
			status = "||"
		}
		if m.isRandom {
			status += " rnd"
		}

		nowPlaying = currentStyle.Render(fmt.Sprintf("%s %s", status, track.Trackname))
		elapsed = dimmedStyle.Render(fmt.Sprintf("%s / %s", FormattedDuration(m.elapsed), track.FormatDuration()))
	}

	return barStyle.Width(m.width - 2).Render(
		fmt.Sprintf("%s\n%s\n%s\n%s", nowPlaying, elapsed, m.progress.ViewAs(percent), m.helpView()))
}

func (m Model) searchBarView() string {
	title := titleStyle.Padding(0, 2).Render(TITLE)
	var hint string
	var query string

	switch m.focus {
	case focusList:
		hint = lipgloss.NewStyle().Padding(0, 2).Render(dimmedStyle.Render(SEARCHBAR_LIST_PLACEHOLDER))
		if m.searchQuery != "" {
			query = dimmedStyle.Render("> " + m.searchQuery)
		}
	case focusSearch:
		hint = lipgloss.NewStyle().Padding(0, 2).Render(dimmedStyle.Render(SEARCHBAR_SEARCH_HINT))
		base := dimmedStyle.Render("> " + m.searchQuery)
		if m.searching {
			query = m.spinner.View() + " " + base
		} else {
			query = base
		}
	default:
		query = dimmedStyle.Render(SEARCHBAR_LIST_PLACEHOLDER)
	}

	if hint != "" {
		if query == "" {
			return fmt.Sprintf("%s\n%s\n", title, hint)
		}
		input := lipgloss.NewStyle().Padding(0, 2).Render(query)
		return fmt.Sprintf("%s\n%s\n%s\n", title, hint, input)
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
	return dimmedStyle.Render(GLOBAL_HELP_TEXT)
}
