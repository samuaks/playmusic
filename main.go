package main

import (
	"fmt"
	"os"
	lib "playmusic/library"
	"playmusic/search"
	"playmusic/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load the local Media directory synchronously so the TUI can start fast
	// with an immediate playlist even on cold startup.
	tracks, err := lib.LoadLibrary("Media")
	if err != nil {
		fmt.Printf("Error loading library: %v\n", err)
		return
	}

	searcher := search.New(search.MockSource{})

	// Scan the rest of the library in the background and stream tracks into the TUI.
	scanCh := make(chan lib.Track)
	go lib.ScanForMedia(lib.BackgroundLibraryDirs(), scanCh)

	ui := tea.NewProgram(
		tui.NewModel(tracks, searcher, scanCh),
		tea.WithAltScreen(),
	)

	if _, err := ui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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
