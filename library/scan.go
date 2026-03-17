package library

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	. "playmusic/helpers"
	"strings"
)

type ScanEvent struct {
	Track *Track
	Err   error
}

/*
ScanForMedia scans the provided directories in the background and emits

	discovered tracks one by one into out. The channel is closed when scanning
	finishes. Missing directories are skipped.
*/
func ScanForMedia(ctx context.Context, dirs []string, out chan<- ScanEvent) {
	defer close(out)

	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- ScanEvent{Err: walkErr}:
				}
				return nil
			}
			if d.IsDir() || !IsSupported(d.Name()) {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			name := d.Name()
			metadata, _ := GetMetadata(path)
			duration, _ := ProbeDuration(path)

			track := Track{
				Trackname: formatTrackName(metadata.Artist, metadata.Title, name),
				Title:     metadata.Title,
				Artist:    metadata.Artist,
				Filename:  name,
				Path:      path,
				Duration:  duration,
				Year:      metadata.Year,
				Genre:     metadata.Genre,
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- ScanEvent{Track: &track}:
				return nil
			}
		})

		if err != nil && !os.IsNotExist(err) && err != context.Canceled {
			select {
			case <-ctx.Done():
				return
			case out <- ScanEvent{Err: err}:
			}
		}
	}
}
