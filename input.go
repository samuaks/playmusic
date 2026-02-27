package main

import (
	"bufio"
	"fmt"
	"os"
	c "playmusic/colors"
	play "playmusic/player"
)

func handleInput(p *play.Player) {
	reader := bufio.NewReader(os.Stdin)
	for {
		rune, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		switch rune {
		case 'p':
			p.Pause()
			fmt.Println(c.Colorize("Paused", c.ColorBold+c.ColorYellow))
		case 'r':
			p.Resume()
			fmt.Println(c.Colorize("Resumed", c.ColorBold+c.ColorGreen))
		case 'n':
			p.Next()
			fmt.Println(c.Colorize("Skipping to next track...", c.ColorBold+c.ColorBlue))
		case 'q':
			p.Stop()
			fmt.Println(c.Colorize("Exiting...", c.ColorBold+c.ColorRed))
			os.Exit(0)
		}
	}
}
