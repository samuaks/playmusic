package library

import (
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/helpers"
	"strings"
	"time"
)

const mediaLibraryDir = "Media"

type Track struct {
	Trackname  string
	Artist     string
	Title      string
	Path       string
	Filename   string
	Duration   time.Duration
	Album      string
	Year       int
	Genre      string
	YTVideoURl string
}

func (t Track) FormatDuration() string {
	return FormattedDuration(t.Duration)
}

func (t Track) Identifier() string {
	if t.Path != "" {
		return t.Path
	}
	return t.YTVideoURl
}

func LoadLibrary(dir string) ([]Track, error) {
	return loadFromDir(dir, make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{}))
}

func LoadLibraries(dirs ...string) ([]Track, error) {
	var tracks []Track
	seenPaths := make(map[string]struct{})
	seenSignatures := make(map[string]struct{})
	seenContents := make(map[string]struct{})

	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		scanned, err := loadFromDir(dir, seenPaths, seenSignatures, seenContents)
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
	return sortingOfTracks(tracks), nil
}

func DefaultLibraryDirs() []string {
	dirs := []string{
		filepath.Clean(mediaLibraryDir),
	}

	home, err := os.UserHomeDir()
	if err == nil && strings.TrimSpace(home) != "" {
		dirs = append(dirs, filepath.Join(home, "Music"))
	}

	return uniqueDirs(dirs)
}

// BackgroundLibraryDirs returns the directories that should be scanned after
// the TUI has already started. The local Media directory is excluded so
// startup stays fast and we do not scan the same source twice.
func BackgroundLibraryDirs() []string {
	var dirs []string
	mediaDir := filepath.Clean(mediaLibraryDir)

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

func loadFromDir(dir string, seenPaths, seenSignatures, seenContents map[string]struct{}) ([]Track, error) {
	var tracks []Track
	state := newScanState(seenPaths, seenSignatures, seenContents)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if !state.shouldInclude(path, d) {
			return nil
		}

		discovered := BuildDiscoveredTrack(path)
		enriched, _ := EnrichTrack(discovered)
		tracks = append(tracks, enriched)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tracks, nil
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
