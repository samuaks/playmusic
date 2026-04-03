package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
// Sidenote: this function should probably be in Library module since it is only used there?
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

// duration formatting
func StringToDuration(duration string) (time.Duration, error) {
	//duration resp format: 4:16
	parts := strings.Split(duration, ":")

	if len(parts) < 1 || len(parts) > 2 {
		return 0, fmt.Errorf("wrong duration format: %s", duration)
	}

	var seconds int
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0, fmt.Errorf("invalid duration part: %w", err)
		}
		seconds = seconds*60 + n
	}

	return time.Duration(seconds) * time.Second, nil
}

func OpenURL(url string) error {
	var cmd string
	var args []string
	switch os := runtime.GOOS; os {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()

}
