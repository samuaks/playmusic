package library

import (
	"path/filepath"
	"strings"
)

func formatTrackName(artist, title, filename string) string {
	if artist == "" || title == "" {
		return strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	if strings.Contains(strings.ToLower(title), strings.ToLower(artist)) {
		return title
	}
	return artist + " - " + title
}
