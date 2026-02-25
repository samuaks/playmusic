package library

import (
	"os"
	"path/filepath"
	"strings"
)

type Track struct {
	Title    string
	Path     string
	Filename string
}

func LoadLibrary(dir string) ([]Track, error) {
	var tracks []Track

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !isSupported(entry.Name()) {
			continue
		}
		name := entry.Name()

		tracks = append(tracks, Track{
			Title:    strings.TrimSuffix(name, filepath.Ext(name)),
			Filename: name,
			Path:     filepath.Join(dir, name),
		})

	}
	return tracks, nil
}
