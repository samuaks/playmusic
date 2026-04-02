package yt_dlp

import (
	"encoding/json"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

type YTResponse struct {
	Entries []YTEntry `json:"entries"`
}

type YTEntry struct {
	URL      string  `json:"url"`
	Author   string  `json:"uploader"`
	Name     string  `json:"title"`
	Duration float64 `json:"duration"`
}

func YTVideoInfoParser(queryRes ytdlp.Result) ([]YTEntry, error) {
	lines := strings.Split(strings.TrimSpace(queryRes.Stdout), "\n")
	entries := make([]YTEntry, 0, len(lines))

	for _, line := range lines {
		var e YTEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

var banned = []string{
	"react",
	"reaction",
	"first time hearing",
	"first listen",
	"analysis",
	"REACT",
	"REACTION",
}

func ClearUnrelatedToMusicGarbage(entries []YTEntry) []YTEntry {
	var filtered = make([]YTEntry, 0, len(entries))

	for _, entry := range entries {
		title := strings.ToLower(entry.Name)

		skip := false
		for _, key := range banned {
			if strings.Contains(title, key) {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		filtered = append(filtered, entry)
	}

	return filtered
}
