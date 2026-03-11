package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
)

type trackDelegate struct {
	list.DefaultDelegate
	currentPath string
}

func newDelegate(currentPath string) trackDelegate {
	d := trackDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
		currentPath:     currentPath,
	}
	d.Styles.SelectedTitle = selectedTitleStyle
	d.Styles.SelectedDesc = selectedDescStyle
	return d
}

func (d trackDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(trackItem)
	if !ok {
		return
	}

	title := t.Title()
	desc := t.Description()

	//isPlaying := index == d.current && d.current != -1
	isPlaying := t.track.Path == d.currentPath
	isSelected := index == m.Index()

	switch {
	case isPlaying && isSelected:
		title = "▶ " + selectedTitleStyle.Render(title)
		desc = "  " + selectedDescStyle.Render(desc)
	case isPlaying:
		title = "▶ " + currentStyle.Render(title)
		desc = "  " + selectedDescStyle.Render(desc)
	case isSelected:
		title = "  " + selectedTitleStyle.Render(title)
		desc = "  " + selectedDescStyle.Render(desc)
	default:
		title = "  " + dimmedStyle.Render(title)
		desc = "  " + dimmedStyle.Render(desc)
	}
	fmt.Fprintf(w, "%s\n%s", title, desc)
}
