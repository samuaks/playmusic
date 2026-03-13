package search

import (
	"playmusic/library"
	"playmusic/yt_dlp"
	"playmusic/ytapi"
)

type YTSource struct{}

func (y YTSource) Name() string { return "YouTube" }

func (y YTSource) Search(query string) ([]library.Track, error) {
	videoUrl, title, err := ytapi.GetVideoURLFromYt(query)
	if err != nil {
		return nil, err
	}

	audioStreamURL, err := yt_dlp.GetStreamURLFromYtDlp(videoUrl)
	if err != nil {
		return nil, err
	}

	track := library.Track{
		Trackname:      title,
		AudioStreamURL: audioStreamURL,
	}

	return []library.Track{track}, nil
}
