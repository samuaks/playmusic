package player

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestStopWithoutPlay(t *testing.T) {
	p := &Player{}
	p.Stop()
}

func TestWaitWithoutPlay(t *testing.T) {
	p := &Player{}
	p.Wait()
}

func TestPlayInvalidFile(t *testing.T) {
	p := &Player{}
	err := p.Play("/does/not/exist.mp3")
	if err == nil {
		t.Error("Expected an error for non-existent file, got nil")
	}
}

func TestStopActuallyStops(t *testing.T) {
	p := &Player{}
	_, filename, _, _ := runtime.Caller(0)
	testfile := filepath.Join(filepath.Dir(filename), "..", "Media", "ShadowsAndDust.mp3")
	err := p.Play(testfile)
	if err != nil {
		t.Skip("Skipping test, unable to play file: ", err)
	}

	time.Sleep(500 * time.Millisecond)
	p.Stop()

	if p.cmd != nil {
		t.Error("Expected cmd to be nil after Stop")
	}
}
