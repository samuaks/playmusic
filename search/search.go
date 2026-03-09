package search

import (
	"playmusic/library"
	"sync"
)

type Source interface {
	Search(query string) ([]library.Track, error)
	Name() string
}

type Searcher struct {
	sources []Source
}

func New(sources ...Source) *Searcher {
	return &Searcher{sources: sources}
}

func (s *Searcher) Search(query string) ([]library.Track, error) {
	results := make(chan []library.Track, len(s.sources))

	var wg sync.WaitGroup
	for _, source := range s.sources {
		wg.Add(1)
		go func(src Source) {
			defer wg.Done()
			tracks, err := src.Search(query)
			if err != nil {
				return
			}
			results <- tracks
		}(source)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	var all []library.Track
	for tracks := range results {
		all = append(all, tracks...)
	}
	return all, nil
}
