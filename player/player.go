package player

import "os/exec"

type Player struct {
	cmd *exec.Cmd
}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
}

func (p *Player) Play(filepath string) error {
	p.Stop()

	p.cmd = exec.Command("ffplay", "-nodisp", "-autoexit", filepath)
	return p.cmd.Start()
}
