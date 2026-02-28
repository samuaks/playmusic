package player

import (
	"fmt"
	d "playmusic/decoder"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	ctrl       *beep.Ctrl
	streamer   beep.StreamSeekCloser
	done       chan struct{}
	next       chan struct{}
	sampleRate beep.SampleRate
}

func (p *Player) Play(path string) error {
	p.Stop()

	streamer, format, err := d.Decode(path)
	if err != nil {
		return fmt.Errorf("Decode failed %s: %w", path, err)
	}

	p.done = make(chan struct{})
	p.next = make(chan struct{}, 1)
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
		close(p.done)
	}))}

	speaker.Play(p.ctrl)
	return nil
}

func (p *Player) Wait() {
	if p.done == nil {
		return
	}
	select {
	case <-p.done:
	case <-p.next:
		p.Stop()
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

func (p *Player) Pause() {
	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
	}
}

func (p *Player) Resume() {
	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = false
		speaker.Unlock()
	}
}

func (p *Player) Next() {
	if p.done != nil {
		p.next <- struct{}{}
	}
}
