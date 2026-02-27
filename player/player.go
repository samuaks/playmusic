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

	p.done = make(chan struct{})
	p.streamer = streamer

	if p.sampleRate == 0 {
		p.sampleRate = format.SampleRate
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	}

	var finalStreamer beep.Streamer
	if format.SampleRate != p.sampleRate {
		finalStreamer = beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
	} else {
		finalStreamer = streamer
	}

	p.ctrl = &beep.Ctrl{Streamer: beep.Seq(finalStreamer, beep.Callback(func() {
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
