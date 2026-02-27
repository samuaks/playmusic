package library

import (
	"fmt"
	"os"
	"path/filepath"
	"playmusic/player"
	"strings"
	"sync"
	"time"
)

type Track struct {
	Title    string
	Path     string
	Filename string
	Duration time.Duration
}

func (t Track) FormatDuration() string {
	minutes := int(t.Duration.Minutes())
	seconds := int(t.Duration.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func LoadLibrary(dir string) ([]Track, error) {
	var tracks []Track

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ffmpegAvailable := player.IsFFmpegAvailable()

	for _, entry := range entries {
		if entry.IsDir() || !isSupported(entry.Name(), ffmpegAvailable) {
			continue
		}
		name := entry.Name()

		tracks = append(tracks, Track{
			Title:    strings.TrimSuffix(name, filepath.Ext(name)),
			Filename: name,
			Path:     filepath.Join(dir, name),
		})

		var wg sync.WaitGroup
		for i := range tracks {
			wg.Add(1)
			go func(t *Track) {
				defer wg.Done()
				t.Duration, _ = probeDuration(t.Path)
			}(&tracks[i])
		}
		wg.Wait()

	}
	return tracks, nil
}
