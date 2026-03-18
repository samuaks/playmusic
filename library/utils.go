package library

import (
	"path/filepath"
	"sort"
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

// initial ascending alphabetical Trackname at the moment
func sortingOfTracks(tracks []Track) []Track {
	sort.SliceStable(tracks, func(i, j int) bool {
		return tracks[i].Trackname < tracks[j].Trackname
	})
	return tracks
}
