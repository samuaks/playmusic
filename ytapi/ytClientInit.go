package ytapi

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var ytClient *youtube.Service //initializing client that will be used for all requests to youtube API

func initiateYTClient() {
	ctx := context.Background() //putting ytservice into long-living context of the app

	var err error
	ytClient, err = youtube.NewService(
		ctx,
		option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")),
	)

	if err != nil {
		log.Fatal(err)
	}
}
