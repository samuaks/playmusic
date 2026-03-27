package library

import (
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	"runtime"
	"strconv"
	"strings"
)

type scanState struct {
	seenPaths      map[string]struct{}
	seenSignatures map[string]struct{}
}

func newScanState(seenPaths, seenSignatures map[string]struct{}) *scanState {
	return &scanState{
		seenPaths:      seenPaths,
		seenSignatures: seenSignatures,
	}
}

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

func (s *scanState) shouldInclude(path string, d fs.DirEntry) bool {
	if d.IsDir() || !IsSupported(d.Name()) {
		return false
	}

	pathKey := normalizedPathKey(path)
	if _, exists := s.seenPaths[pathKey]; exists {
		return false
	}
	s.seenPaths[pathKey] = struct{}{}

	// Use a cheap filename+size heuristic before enrichment so obvious
	// duplicates are filtered without hashing whole files on the startup path.
	sigKey, sigErr := fileSignatureKey(d)
	if sigErr == nil && sigKey != "" {
		if _, exists := s.seenSignatures[sigKey]; exists {
			return false
		}
		s.seenSignatures[sigKey] = struct{}{}
	}

	return true
}

func (s *scanState) rememberTrack(track Track) {
	if track.Path == "" {
		return
	}

	pathKey := normalizedPathKey(track.Path)
	s.seenPaths[pathKey] = struct{}{}

	info, err := os.Stat(track.Path)
	if err != nil {
		return
	}

	name := track.Filename
	if strings.TrimSpace(name) == "" {
		name = filepath.Base(track.Path)
	}
	sigKey := strings.ToLower(name) + ":" + strconv.FormatInt(info.Size(), 10)
	s.seenSignatures[sigKey] = struct{}{}
}
