package library

import (
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	. "playmusic/helpers"
	"strings"
	"sync"
	"time"
)

type Track struct {
	Trackname string
	Artist    string
	Title     string
	Path      string
	Filename  string
	Duration  time.Duration
	Album     string
	Year      int
	Genre     string
}

func (t Track) FormatDuration() string {
	return FormattedDuration(t.Duration)
}

func LoadLibrary(dir string) ([]Track, error) {
	return loadFromDir(dir)
}

func LoadLibraries(dirs ...string) ([]Track, error) {
	var tracks []Track
	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		scanned, err := loadFromDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		tracks = append(tracks, scanned...)
	}
	if len(tracks) == 0 {
		return nil, nil
	}
	return sortingOfTracks(deduplicateTracks(enrichTracks(tracks))), nil
}

func DefaultLibraryDirs() []string {
	dirs := []string{
		filepath.Clean("Media"),
	}

	home, err := os.UserHomeDir()
	if err == nil && strings.TrimSpace(home) != "" {
		dirs = append(dirs, filepath.Join(home, "Music"))
	}

	return uniqueDirs(dirs)
}

/*  BackgroundLibraryDirs returns the directories that should be scanned
after the TUI has already started. The local Media directory is excluded
so startup stays fast and we do not scan the same source twice.*/

func BackgroundLibraryDirs() []string {
	var dirs []string
	mediaDir := filepath.Clean("Media")

	for _, dir := range DefaultLibraryDirs() {
		if filepath.Clean(dir) == mediaDir {
			continue
		}
		dirs = append(dirs, dir)
	}

	return dirs
}

func LoadDefaultLibrary() ([]Track, error) {
	return LoadLibraries(DefaultLibraryDirs()...)
}

func loadFromDir(dir string) ([]Track, error) {
	var tracks []Track

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !IsSupported(d.Name()) {
			return nil
		}

		tracks = append(tracks, BuildDiscoveredTrack(path))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deduplicateTracks(enrichTracks(tracks)), nil
}

func uniqueDirs(dirs []string) []string {
	seen := make(map[string]bool)
	uniq := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		cleaned := filepath.Clean(dir)
		if cleaned == "." || seen[cleaned] {
			continue
		}
		seen[cleaned] = true
		uniq = append(uniq, cleaned)
	}
	return uniq
}

func enrichTracks(tracks []Track) []Track {
	var wg sync.WaitGroup

	for i := range tracks {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			enriched, _ := EnrichTrack(tracks[idx])
			tracks[idx] = enriched
		}(i)
	}

	wg.Wait()
	return tracks
}

func deduplicateTracks(tracks []Track) []Track {
	var wg sync.WaitGroup
	hashes := make([]string, len(tracks))

	for i := range tracks {
		wg.Add(1)
		go func(idx int, t Track) {
			defer wg.Done()
			hashes[idx], _ = hashFile(t.Path)
		}(i, tracks[i])
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
	return uniqueTracks
}
