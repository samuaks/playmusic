package search

import (
	"playmusic/library"
	"strings"
	"time"
)

type MockSource struct{}

func (m MockSource) Name() string { return "mock" }

func (m MockSource) Search(query string) ([]library.Track, error) {
	time.Sleep(1 * time.Second) // simulate API delay

	tracks, err := library.LoadLibrary("Mock")
	if err != nil {
		return nil, err
	}

	var results []library.Track
	for _, t := range tracks {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(query)) {
			results = append(results, t)
		}
	}
	return results, nil
}
