package library

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	. "playmusic/decoder"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type scanState struct {
	seenPaths      map[string]struct{}
	seenSignatures map[string]struct{}
	seenContents   map[string]struct{}
}

const fingerprintChunkSize int64 = 64 * 1024

func newScanState(seenPaths, seenSignatures, seenContents map[string]struct{}) *scanState {
	if seenContents == nil {
		seenContents = make(map[string]struct{})
	}
	return &scanState{
		seenPaths:      seenPaths,
		seenSignatures: seenSignatures,
		seenContents:   seenContents,
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

func contentKey(path string, size int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()

	offsets := sampleOffsets(size)
	buf := make([]byte, fingerprintChunkSize)
	for _, offset := range offsets {
		if _, err := file.Seek(offset, io.SeekStart); err != nil {
			return "", err
		}

		remaining := size - offset
		if remaining <= 0 {
			continue
		}

		chunkLen := min(remaining, fingerprintChunkSize)

		n, readErr := io.ReadFull(file, buf[:chunkLen])
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			return "", readErr
		}
		if n > 0 {
			if _, err := hasher.Write(buf[:n]); err != nil {
				return "", err
			}
		}
	}

	return strconv.FormatInt(size, 10) + ":" + hex.EncodeToString(hasher.Sum(nil)), nil
}

func sampleOffsets(size int64) []int64 {
	if size <= 0 {
		return []int64{0}
	}

	offsets := []int64{0}

	if size > fingerprintChunkSize {
		middle := max(size/2-fingerprintChunkSize/2, 0)
		last := max(size-fingerprintChunkSize, 0)

		offsets = appendUniqueOffset(offsets, middle)
		offsets = appendUniqueOffset(offsets, last)
	}

	return offsets
}

func appendUniqueOffset(offsets []int64, candidate int64) []int64 {
	if slices.Contains(offsets, candidate) {
		return offsets
	}
	return append(offsets, candidate)
}

func (s *scanState) shouldInclude(path string, d fs.DirEntry) bool {
	if d.IsDir() || !IsSupported(d.Name()) {
		return false
	}

	pathKey := normalizedPathKey(path)
	if _, exists := s.seenPaths[pathKey]; exists {
		return false
	}

	// Use a cheap filename+size heuristic before enrichment so obvious
	// duplicates are filtered without hashing whole files on the startup path.
	sigKey, sigErr := fileSignatureKey(d)
	if sigErr == nil && sigKey != "" {
		if _, exists := s.seenSignatures[sigKey]; exists {
			return false
		}
	}

	info, infoErr := d.Info()
	if infoErr != nil {
		return false
	}
	if exactKey, err := contentKey(path, info.Size()); err == nil && exactKey != "" {
		if _, exists := s.seenContents[exactKey]; exists {
			return false
		}
		s.seenContents[exactKey] = struct{}{}
	}

	s.seenPaths[pathKey] = struct{}{}
	if sigErr == nil && sigKey != "" {
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

	if exactKey, err := contentKey(track.Path, info.Size()); err == nil && exactKey != "" {
		s.seenContents[exactKey] = struct{}{}
	}
}
