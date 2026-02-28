package library

import (
	"fmt"
	"os"
	"path/filepath"
	. "playmusic/decoder"
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

	for _, entry := range entries {
		if entry.IsDir() || !IsSupported(entry.Name()) {
			continue
		}
		name := entry.Name()

		tracks = append(tracks, Track{
			Title:    strings.TrimSuffix(name, filepath.Ext(name)),
			Filename: name,
			Path:     filepath.Join(dir, name),
		})
	}

	var wg sync.WaitGroup
	hashes := make([]string, len(tracks))
	for i := range tracks {
		wg.Add(1)
		go func(idx int, t *Track) {
			defer wg.Done()
			t.Duration, _ = ProbeDuration(t.Path)
			hashes[idx], _ = hashFile(t.Path)
		}(i, &tracks[i])
	}
	wg.Wait()

	seen := make(map[string]bool)
	uniqueTracks := tracks[:0]
	for i, track := range tracks {
		if hashes[i] == "" || seen[hashes[i]] {
			continue
		}
		seen[hashes[i]] = true
		uniqueTracks = append(uniqueTracks, track)
	}
	return uniqueTracks, nil
}
