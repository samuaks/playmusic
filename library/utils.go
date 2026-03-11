package library

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func formatTrackName(artist, title, filename string) string {
	if artist == "" || title == "" {
		return strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	if strings.Contains(strings.ToLower(title), strings.ToLower(artist)) {
		return title
	}
	return artist + " - " + title
}

// ascending alphabetical Trackname at the moment
func sortingOfTracks(tracks []Track) []Track {
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Trackname < tracks[j].Trackname
	})
	return tracks
}
