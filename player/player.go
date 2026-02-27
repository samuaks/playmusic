package player

import (
	"fmt"
	c "playmusic/colors"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	ctrl       *beep.Ctrl
	streamer   beep.StreamSeekCloser
	done       chan struct{}
	sampleRate beep.SampleRate
}

func (p *Player) Play(path string) error {
	p.Stop()

	streamer, format, err := decode(path)
	if err != nil {
		return fmt.Errorf("Decode failed %s: %w", path, err)
	}

	fmt.Printf(c.Colorize("debug: ", c.ColorBold+c.ColorCyan)+"decoded successfully, format: %v\n", format)

	p.done = make(chan struct{})
	p.streamer = streamer

	if format.SampleRate != p.sampleRate {
		p.sampleRate = format.SampleRate
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		fmt.Printf(c.Colorize("debug: ", c.ColorBold+c.ColorCyan)+"speaker initialized with sample rate %v\n", format.SampleRate)
	}

	p.ctrl = &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		fmt.Println(c.Colorize("debug: callback fired, closing done", c.ColorBold+c.ColorCyan))
		close(p.done)
	}))}

	speaker.Play(p.ctrl)
	fmt.Println(c.Colorize("debug: ", c.ColorBold+c.ColorCyan) + "speaker.Play called, waiting...")
	return nil
}

func (p *Player) Wait() {
	if p.done != nil {
		<-p.done
	}
}

func (p *Player) Stop() {
	if p.ctrl != nil {
		speaker.Clear()
		p.ctrl = nil
	}
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
}
