package tui

import "fmt"

func setTerminalTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
