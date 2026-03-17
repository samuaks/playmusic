package library

import "io/fs"

type scanCandidate struct {
	discovered Track
	enriched   Track
	enrichErr  error
	include    bool
}

func processCandidate(state *scanState, path string, d fs.DirEntry) scanCandidate {
	if !state.shouldInclude(path, d) {
		return scanCandidate{}
	}

	discovered := BuildDiscoveredTrack(path)
	enriched, enrichErr := EnrichTrack(discovered)

	return scanCandidate{
		discovered: discovered,
		enriched:   enriched,
		enrichErr:  enrichErr,
		include:    true,
	}
}
