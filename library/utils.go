package library

import (
	"path/filepath"
	"strings"
)

var supported = map[string]bool{
	".mp3":  true,
	".flac": true,
	".wav":  true,
	".m4a":  true,
	".ogg":  true,
	".aac":  true,
	".opus": true,
}

func isSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return supported[ext]
}
