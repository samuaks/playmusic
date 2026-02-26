package player

import (
	"fmt"
	"os"
	"os/exec"
)

type Player struct {
	cmd *exec.Cmd
}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

func (p *Player) Wait() {
	if p.cmd != nil {
		p.cmd.Wait()
	}
}

func (p *Player) Play(filepath string) error {
	p.Stop()

	if _, err := os.Stat(filepath); err != nil {
		return fmt.Errorf("File not found: %s", filepath)
	}

	p.cmd = exec.Command("ffplay", "-nodisp", "-autoexit", filepath)
	return p.cmd.Start()
}
