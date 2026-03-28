package paths

import (
	"os"
	"path/filepath"
)

const MediaLibraryDir = "Media"

func UserMediaDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, MediaLibraryDir), nil
}
