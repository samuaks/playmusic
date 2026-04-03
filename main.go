package main

import (
	"context"
	"fmt"
	"log"

	"os"
	"playmusic/ffmpeg"
	lib "playmusic/library"
	"playmusic/search"
	"playmusic/tui"
	"playmusic/yt_dlp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	const localMediaDir = "Media"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ui *tea.Program

	var ffmpegAvailable = ffmpeg.InitFFmpeg()

	if len(os.Args) > 1 && os.Args[1] == "--radio" {
		if !ffmpegAvailable {
			log.Fatal("Can't run radio without ffmpeg installed.")
		}

		resInst, err := ytdlp.Install(context.TODO(), nil)
		if err != nil {
			fmt.Printf("Warning: failed to install yt-dlp: %v\n", err)
		} else {
			yt_dlp.SetBinaryPath(resInst.Executable)
		}

		searcher := search.New(search.YTRadioSource{})
		var tracks []lib.Track

		newOnlineModel := tui.NewOnlineModel(tracks, searcher)

		ui = tea.NewProgram(
			&newOnlineModel,
			tea.WithAltScreen(),
		)
	} else {
		tracks, err := lib.LoadLibrary(localMediaDir)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Warning: failed to load local library: %v\n", err)
			}
			tracks = nil
		}

		scanCh := make(chan lib.ScanEvent)
		go lib.ScanForMedia(ctx, lib.BackgroundLibraryDirs(), scanCh)
		ui = tea.NewProgram(
			tui.NewModel(tracks, scanCh),
			tea.WithAltScreen(),
		)
	}

	if _, err := ui.Run(); err != nil {
		log.Fatal(fmt.Printf("Error: %v\n", err))
	}
}
