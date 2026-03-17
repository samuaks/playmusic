package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type multiValueFlag []string

func (m *multiValueFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiValueFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("directory value cannot be empty")
	}
	*m = append(*m, value)
	return nil
}

func resolveLibraryRoots(cliDirs []string, systemMusic bool, scanAll bool) ([]string, error) {
	roots := make([]string, 0, len(cliDirs)+4)
	if len(cliDirs) > 0 {
		roots = append(roots, cliDirs...)
	} else if env := strings.TrimSpace(os.Getenv("PLAYMUSIC_DIR")); env != "" {
		for _, dir := range strings.Split(env, string(os.PathListSeparator)) {
			if trimmed := strings.TrimSpace(dir); trimmed != "" {
				roots = append(roots, trimmed)
			}
		}
	} else {
		roots = append(roots, "Media")
	}

	if systemMusic {
		dir, err := getSystemMusicDir()
		if err != nil {
			return nil, err
		}
		roots = append(roots, dir)
	}

	if scanAll {
		roots = append(roots, getScanAllRoots()...)
	}

	unique := make([]string, 0, len(roots))
	seen := make(map[string]bool, len(roots))
	for _, root := range roots {
		clean := filepath.Clean(root)
		abs, err := filepath.Abs(clean)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path %q: %w", root, err)
		}
		key := strings.ToLower(abs)
		if runtime.GOOS != "windows" {
			key = abs
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, abs)
	}
	return unique, nil
}

func getSystemMusicDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user home directory: %w", err)
	}
	return filepath.Join(home, "Music"), nil
}

func getScanAllRoots() []string {
	if runtime.GOOS == "windows" {
		roots := make([]string, 0, 4)
		for letter := 'A'; letter <= 'Z'; letter++ {
			root := fmt.Sprintf("%c:\\", letter)
			if info, err := os.Stat(root); err == nil && info.IsDir() {
				roots = append(roots, root)
			}
		}
		return roots
	}
	return []string{string(os.PathSeparator)}
}
