package main

import (
	"fmt"
	"os"
	"os/exec"
)

func playMP3(filepath string) error {
	// ffplay comes with ffmpeg
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", filepath)
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <mp3file>")
		return
	}

	fmt.Println("Playing audio...")
	if err := playMP3(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
