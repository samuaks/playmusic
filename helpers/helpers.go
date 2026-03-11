package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/dhowden/tag"
)

func FormattedDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

type Metadata struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Genre  string
}

// GetMetadata is used in library.loadFromDir() func after we have checked that format is suitable for fetching the data
func GetMetadata(path string) (Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return Metadata{}, err
	}
	defer file.Close()

	data, err := tag.ReadFrom(file)
	if err != nil {
		return Metadata{}, err
	}

	metadata := Metadata{
		Title:  data.Title(),
		Artist: data.Artist(),
		Album:  data.Album(),
		Year:   data.Year(),
		Genre:  data.Genre(),
	}

	return metadata, nil
}
