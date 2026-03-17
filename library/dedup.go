package library

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
