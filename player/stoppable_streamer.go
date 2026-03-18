package player

import (
	"fmt"
	"sync/atomic"

	"github.com/gopxl/beep/v2"
)

type stoppableStreamer struct {
	inner   beep.Streamer
	stopped *atomic.Bool
	err     error
}

func (s *stoppableStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.stopped.Load() {
		s.err = fmt.Errorf("stream stopped")
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
