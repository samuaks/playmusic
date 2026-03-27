package search

import (
	"playmusic/library"
	"playmusic/yt_dlp"
)

type YTRadioSource struct{}

func (y YTRadioSource) Name() string { return "YouTubeRadio" }

func (y YTRadioSource) Search(query string) ([]library.Track, error) {
	var tracks []library.Track

	trackInfo, err := yt_dlp.GetMusicJamPlaylistWithQueryJson(query)
	if err != nil {
		return nil, err
	}

	for _, info := range trackInfo {
		track := library.Track{
			Trackname:  info.Trackname,
			YTVideoURL: info.YTVideoURL,
			Duration:   info.Duration,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
