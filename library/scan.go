package library

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	"runtime"
	"strconv"
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
	ScanEventDone
)

func normalizedPathKey(path string) string {
	cleaned := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		return strings.ToLower(cleaned)
	}
	return cleaned
}

func fileSignatureKey(d fs.DirEntry) (string, error) {
	info, err := d.Info()
	if err != nil {
		return "", err
	}
	return strings.ToLower(d.Name()) + ":" + strconv.FormatInt(info.Size(), 10), nil
}

/*
ScanForMedia scans the provided directories in the background and emits

	discovered tracks one by one into out. The channel is closed when scanning
	finishes. Missing directories are skipped.
*/
func ScanForMedia(ctx context.Context, dirs []string, out chan<- ScanEvent) {
	defer close(out)
	seenPaths := make(map[string]struct{})
	seenSignatures := make(map[string]struct{})

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

			pathKey := normalizedPathKey(path)
			if _, exists := seenPaths[pathKey]; exists {
				return nil
			}
			seenPaths[pathKey] = struct{}{}

			sigKey, sigErr := fileSignatureKey(d)
			if sigErr == nil && sigKey != "" {
				if _, exists := seenSignatures[sigKey]; exists {
					return nil
				}
				seenSignatures[sigKey] = struct{}{}
			}

			discovered := BuildDiscoveredTrack(path)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- ScanEvent{
				Type:  ScanEventDiscovered,
				Track: &discovered,
			}:
			}

			enriched, enrichErr := EnrichTrack(discovered)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- ScanEvent{
				Type:  ScanEventEnriched,
				Track: &enriched,
				Err:   enrichErr,
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
