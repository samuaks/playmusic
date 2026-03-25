package search

import (
	"playmusic/library"
	"playmusic/yt_dlp"
)

type YTSource struct{}

func (y YTSource) Name() string { return "YouTube" }

func (y YTSource) Search(query string) ([]library.Track, error) {
	videoUrl, title, duration, err := yt_dlp.GetYTVideoInfo(query)
	if err != nil {
		return nil, err
	}

	track := library.Track{
		Trackname:  "Streaming: " + title,
		YTVideoURL: videoUrl,
		Duration:   duration,
	}

	return []library.Track{track}, nil
}

func (y YTSource) MakePlaylistWithQuery(query string) ([]library.Track, error) {
	var tracks []library.Track

	trackInfo, err := yt_dlp.GetMusicJamPlaylistWithQuery(query)
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
