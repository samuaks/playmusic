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
