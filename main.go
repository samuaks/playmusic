package main

import (
	"fmt"
	lib "playmusic/library"
	play "playmusic/player"
)

func main() {
	tracks, err := lib.LoadLibrary("Media")
	if err != nil {
		fmt.Printf("Error loading library: %v\n", err)
		return
	}

	if len(tracks) == 0 {
		fmt.Println("No tracks found in Media/")
		return
	}

	fmt.Printf("Loaded %d tracks:\n", len(tracks))
	for i, track := range tracks {
		fmt.Printf(" [%d] %s\n", i+1, track.Title)
	}

	p := &play.Player{}
	fmt.Printf("\nPlaying: %s\n", tracks[0].Title)
	if err := p.Play(tracks[0].Path); err != nil {
		fmt.Printf("Error playing track: %v\n", err)
	}
}
