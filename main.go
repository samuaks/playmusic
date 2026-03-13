package main

import (
	"context"
	"fmt"
	"log"

	"os"
	lib "playmusic/library"
	"playmusic/search"
	"playmusic/tui"
	"playmusic/ytapi"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	err := godotenv.Load() //loading .env for global variables
	if err != nil {
		log.Fatal("Can't load secrets from .env")
	}

	tracks, err := lib.LoadDefaultLibrary()
	if err != nil {
		fmt.Printf("Error loading library: %v\n", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No tracks found in Media/ or default Music folder")
		return
	}

	searcher := search.New(search.MockSource{})

	ui := tea.NewProgram(
		tui.NewModel(tracks, searcher), tea.WithAltScreen())

	if _, err := ui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	ytapi.InitiateYTClient()

	//Checking for the yt-dlp binary and installing it if it's not present.
	//This ensures that the client is ready to use when the app starts (used in main()).
	ytdlp.MustInstall(context.TODO(), nil)
}

// fmt.Printf("Loaded %d tracks:\n", len(tracks))
// for i, track := range tracks {
// 	index := fmt.Sprintf("%s", "["+c.Colorize(fmt.Sprintf("%d", i+1), c.ColorBold+c.ColorCyan)+"]")
// 	title := c.Colorize(track.Title, c.ColorWhite)
// 	duration := fmt.Sprintf("%s", "("+c.Colorize(fmt.Sprintf("%s", track.FormatDuration()), c.ColorBold+c.ColorCyan)+")")
// 	fmt.Printf(" %s %s %s\n", index, title, duration)
// }

// p := &play.Player{}
// // spawn goroutine to handle user input while songs are playing
// // because i think this is the only reasonable way to achieve this without blocking main thread.
// go handleInput(p)

// for _, track := range tracks {
// 	fmt.Printf("\nPlaying: %s\n", track.Title)
// 	if err := p.Play(track.Path); err != nil {
// 		fmt.Printf("Error playing track: %v\n", err)
// 	}
// 	p.Wait()
// }
// fmt.Println("\nPlayList finished.")
