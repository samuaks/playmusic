package main

import (
	"context"
	"fmt"
	"os"
	lib "playmusic/library"
	"playmusic/search"
	"playmusic/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const localMediaDir = "Media"

func main() {
	// Load the local Media directory synchronously so the TUI can start fast
	// with an immediate playlist even on cold startup.
	tracks, err := lib.LoadLibrary(localMediaDir)
	if err != nil {
		fmt.Printf("Error loading library: %v\n", err)
		return
	}

	searcher := search.New(search.MockSource{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Scan the rest of the library in the background and stream tracks into the TUI.
	scanCh := make(chan lib.ScanEvent)
	go lib.ScanForMedia(ctx, lib.BackgroundLibraryDirs(), scanCh)

	ui := tea.NewProgram(
		tui.NewModel(tracks, searcher, scanCh),
		tea.WithAltScreen(),
	)

	if _, err := ui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
