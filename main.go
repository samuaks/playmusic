package main

import (
	"fmt"
	c "playmusic/colors"
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
		index := fmt.Sprintf("%s", "["+c.Colorize(fmt.Sprintf("%d", i+1), c.ColorBold+c.ColorCyan)+"]")
		title := c.Colorize(track.Title, c.ColorWhite)
		duration := fmt.Sprintf("%s", "("+c.Colorize(fmt.Sprintf("%s", track.FormatDuration()), c.ColorBold+c.ColorCyan)+")")
		fmt.Printf(" %s %s %s\n", index, title, duration)
	}

	p := &play.Player{}
	go handleInput(p)

	for _, track := range tracks {
		fmt.Printf("\nPlaying: %s\n", track.Title)
		if err := p.Play(track.Path); err != nil {
			fmt.Printf("Error playing track: %v\n", err)
		}
		p.Wait()
	}
	fmt.Println("\nPlayList finished.")
}
