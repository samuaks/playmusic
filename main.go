package main

import (
	"context"
	"fmt"

	"os"
	lib "playmusic/library"
	"playmusic/search"
	"playmusic/tui"
	"playmusic/yt_dlp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	const localMediaDir = "Media"

	tracks, err := lib.LoadLibrary(localMediaDir)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to load local library: %v\n", err)
		}
		tracks = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ui *tea.Program

	if len(os.Args) > 1 && os.Args[1] == "--radio" {
		resInst, err := ytdlp.Install(context.TODO(), nil)
		if err != nil {
			fmt.Printf("Warning: failed to install yt-dlp: %v\n", err)
		} else {
			yt_dlp.SetBinaryPath(resInst.Executable)
		}
		searcher := search.New(search.YTRadioSource{})
		ui = tea.NewProgram(
			tui.NewOnlineModel(tracks, searcher),
			tea.WithAltScreen(),
		)
	} else {
		scanCh := make(chan lib.ScanEvent)
		go lib.ScanForMedia(ctx, lib.BackgroundLibraryDirs(), scanCh)
		ui = tea.NewProgram(
			tui.NewModel(tracks, scanCh),
			tea.WithAltScreen(),
		)
	}

	if _, err := ui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
