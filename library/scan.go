package library

import (
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	. "playmusic/helpers"
	"strings"
)

/* ScanForMedia scans the provided directories in the background and emits
   discovered tracks one by one into out. The channel is closed when scanning
   finishes. Missing directories are skipped.*/
func ScanForMedia(dirs []string, out chan<- Track) {
	defer close(out)

	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				// Skip unreadable entries so one bad subtree does not abort the whole scan.
				return nil
			}
			if d.IsDir() || !IsSupported(d.Name()) {
				return nil
			}

			name := d.Name()
			metadata, _ := GetMetadata(path)
			duration, _ := ProbeDuration(path)

			out <- Track{
				Trackname: formatTrackName(metadata.Artist, metadata.Title, name),
				Title:     metadata.Title,
				Artist:    metadata.Artist,
				Filename:  name,
				Path:      path,
				Duration:  duration,
				Year:      metadata.Year,
				Genre:     metadata.Genre,
			}

			return nil
		})

		if err != nil && !os.IsNotExist(err) {
			// Ignore directory-level failures for now
			continue
		}
	}
}
