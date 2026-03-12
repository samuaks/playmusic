package ytapi

import (
	"context"
	"fmt"
	"time"
)

func getVideoURLFromYt(query string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //setting timeout for the request
	defer cancel()

	apicall := ytClient.Search.List([]string{"id", "snippet"}).
		Q(query).
		Type("video").
		MaxResults(1).
		Context(ctx)

	response, err := apicall.Do()
	if err != nil {
		return "", "", err
	}

	if len(response.Items) == 0 {
		return "", "", fmt.Errorf("No video found for query: %s", query)
	}

	item := response.Items[0] // extracting video ID and title from the response
	videoId := item.Id.VideoId
	title := item.Snippet.Title

	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId) //constructing video URL
	return videoURL, title, nil
}
