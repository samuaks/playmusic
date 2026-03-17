package library

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ScanEvent struct {
	Type  ScanEventType
	Track *Track
	Err   error
}

type ScanEventType int

const (
	ScanEventDiscovered ScanEventType = iota
	ScanEventEnriched
)

// ScanForMedia scans the provided directories in the background and emits
// discovered and enriched tracks one by one into out. The channel is closed
// when scanning finishes. Missing directories are skipped.
func ScanForMedia(ctx context.Context, dirs []string, out chan<- ScanEvent) {
	defer close(out)
	state := newScanState(make(map[string]struct{}), make(map[string]struct{}))

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

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			candidate := processCandidate(state, path, d)
			if !candidate.include {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- ScanEvent{
				Type:  ScanEventDiscovered,
				Track: &candidate.discovered,
			}:
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- ScanEvent{
				Type:  ScanEventEnriched,
				Track: &candidate.enriched,
				Err:   candidate.enrichErr,
			}:
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
