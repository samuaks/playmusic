package ytapi

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestLoadTracks(t *testing.T) {
	initiateYTClientForTest()

	ytVideoId, title, err := getVideoIdFromYt("sleep token chokehold")
	if err != nil {
		t.Fatal("error occurred while fetching video ID:", err)
	}

	if ytVideoId == "" {
		t.Fatal("expected a video ID, got an empty string")
	}
	if title == "" {
		t.Fatal("expected a title, got an empty string")
	}

	t.Logf("Video ID: %s, Title: %s", ytVideoId, title)
}

func initiateYTClientForTest() {
	//have to load .env manualy as main() func is not called in tests
	errInEnv := godotenv.Load("../.env")
	if errInEnv != nil {
		log.Fatal("Can't load secrets from .env")
	}

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
