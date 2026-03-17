package library

import "io/fs"

type scanCandidate struct {
	discovered Track
	enriched   Track
	enrichErr  error
	include    bool
}

// processCandidate keeps discovery and enrichment as separate stages:
// first we cheaply decide whether the file should be considered at all,
// then we build a minimal track shape, and only after that do metadata and
// duration probing for the enriched result.
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
