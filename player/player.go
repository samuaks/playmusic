package player

import (
	"fmt"
	"os/exec"
	"path/filepath"
	d "playmusic/decoder"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	ctrl          *beep.Ctrl
	streamer      beep.StreamSeekCloser
	externalCmd   *exec.Cmd
	done          chan struct{}
	doneClosed    atomic.Bool
	next          chan struct{}
	sampleRate    beep.SampleRate
	closeStreamFn func()
}
var launchVideoPlayer = defaultLaunchVideoPlayer

func (p *Player) signalDone() {
	if p.done != nil && p.doneClosed.CompareAndSwap(false, true) {
		close(p.done)
	}
}

func (p *Player) Play(path string) error {
	p.Stop()
	if strings.EqualFold(filepath.Ext(path), ".mp4") {
		return p.playVideo(path)
	}

	streamer, format, err := d.Decode(path)
	if err != nil {
		return fmt.Errorf("Decode failed %s: %w", path, err)
	}

	p.done = make(chan struct{})
	p.doneClosed.Store(false)
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
		p.signalDone()
	}))}

	speaker.Play(p.ctrl)
	return nil
}

func defaultLaunchVideoPlayer(path string) (*exec.Cmd, error) {
	cmd := exec.Command("ffplay", "-autoexit", "-loglevel", "error", path)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (p *Player) playVideo(path string) error {
	cmd, err := launchVideoPlayer(path)
	if err != nil {
		return fmt.Errorf("Video playback failed %s: %w", path, err)
	}

	p.done = make(chan struct{})
	p.doneClosed.Store(false)
	p.next = make(chan struct{}, 1)
	p.externalCmd = cmd

	go func() {
		_ = cmd.Wait()
		p.signalDone()
	}()

	return nil
}

// there is a delay ~1-2 seconds and player counts it as a time of a track most likely
// would want to make the one Play func, was already almost done, but downstream functions was
// reliant on beep.StreamSeekCloser and something went wrong
// DecodeStreamUrl returns beep.Streamer, not beep.StreamSeekCloser as Decode()
func (p *Player) PlayFromSearch(url string) error {
	var stopped atomic.Bool
	stopped.Store(false)

	p.Stop()

	streamer, format, closeStream, err := d.DecodeStreamUrl(url)
	if err != nil {
		return fmt.Errorf("Decode failed %s: %w", url, err)
	}

	p.done = make(chan struct{})
	p.doneClosed.Store(false)
	p.next = make(chan struct{}, 1)

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

	wrappedStreamer := &stoppableStreamer{
		inner:   finalStreamer,
		stopped: &stopped,
	}

	p.ctrl = &beep.Ctrl{Streamer: beep.Seq(wrappedStreamer, beep.Callback(func() {
		p.signalDone()
	}))}

	p.closeStreamFn = func() {
		stopped.Store(true)
		go closeStream() //different goroutine to end the stream
	}

	speaker.Play(p.ctrl)
	return nil
}

// Wait blocks until the track finishes naturally (returns true) or is interrupted (returns false).
func (p *Player) Wait() bool {
	if p.done == nil {
		return false
	}

	select {
	case <-p.done:
		return true
	case <-p.next:
		p.Stop()
		return false
	}
}

func (p *Player) Done() chan struct{} {
	return p.done
}

func (p *Player) Stop() {
	if p.ctrl != nil {
		speaker.Clear()
	}
	if p.externalCmd != nil {
		if p.externalCmd.Process != nil {
			_ = p.externalCmd.Process.Kill()
		}
		p.externalCmd = nil
	}

	if p.closeStreamFn != nil {
		p.closeStreamFn()
		p.closeStreamFn = nil
	}

	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}

	p.ctrl = nil
	p.signalDone()
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

func (p *Player) Skip() {
	if p.done != nil {
		p.Stop()
		p.done = nil
	}
}

func (p *Player) IsPlaying() bool {
	return p.ctrl != nil || p.externalCmd != nil
}
