package library

import (
	"errors"
	"path/filepath"

	. "playmusic/decoder"
	. "playmusic/helpers"
)

// BuildDiscoveredTrack creates the minimal track shape that is safe to show
// before metadata and duration probing have completed.
func BuildDiscoveredTrack(path string) Track {
	filename := filepath.Base(path)

	return Track{
		Path:      path,
		Filename:  filename,
		Trackname: formatTrackName("", "", filename),
	}
}

// EnrichTrack fills the metadata-dependent fields for an already discovered
// track while preserving its identity fields.
func EnrichTrack(track Track) (Track, error) {
	metadata, metadataErr := GetMetadata(track.Path)
	duration, durationErr := ProbeDuration(track.Path)

	track.Title = metadata.Title
	track.Artist = metadata.Artist
	track.Album = metadata.Album
	track.Year = metadata.Year
	track.Genre = metadata.Genre
	track.Duration = duration
	track.Trackname = formatTrackName(metadata.Artist, metadata.Title, track.Filename)

	return track, errors.Join(metadataErr, durationErr)
}
