package player

import (
	"sync/atomic"

	"github.com/gopxl/beep/v2"
)

// wrapper around beep.Streamer to allow stopping the stream from another goroutine
type stoppableStreamer struct {
	inner   beep.Streamer
	stopped *atomic.Bool // indicates whether the stream has been stopped already assuming that the stream can be stopped from another goroutine
	err     error
}

func (s *stoppableStreamer) Stream(samples [][2]float64) (n int, ok bool) { //pointed to specific stoppableStreamer to change it's state
	if s.stopped.Load() {
		return 0, false
	}
	n, ok = s.inner.Stream(samples)
	if !ok {
		s.err = s.inner.Err()
	}
	return n, ok
}

func (s *stoppableStreamer) Err() error {
	return s.err
}
