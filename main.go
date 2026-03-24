package main

import (
	"context"
	"fmt"

	"os"
	lib "playmusic/library"
	"playmusic/tui"
	"playmusic/yt_dlp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	resInst, err := ytdlp.Install(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	yt_dlp.SetBinaryPath(resInst.Executable)

	const localMediaDir = "Media"

	// Load the local Media directory synchronously so the TUI can start fast
	// with an immediate playlist even on cold startup.
	tracks, err := lib.LoadLibrary(localMediaDir)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to load local library: %v\n", err)
		}
		tracks = nil
	}

	/* 	searcher := search.New(search.YTSource{}) */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Scan the rest of the library in the background and stream tracks into the TUI.
	scanCh := make(chan lib.ScanEvent)
	go lib.ScanForMedia(ctx, lib.BackgroundLibraryDirs(), scanCh)

	ui := tea.NewProgram(
		tui.NewModel(tracks, scanCh),
		tea.WithAltScreen(),
	)

	if _, err := ui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
