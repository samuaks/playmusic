package player

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func testFile(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "Media", name)
}

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
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestStopActuallyStops(t *testing.T) {
	p := &Player{}
	err := p.Play(testFile("ShadowsAndDust.mp3"))
	if err != nil {
		t.Skipf("skipping, unable to play file: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	p.Stop()

	if p.ctrl != nil {
		t.Error("expected ctrl to be nil after Stop")
	}
	if p.streamer != nil {
		t.Error("expected streamer to be nil after Stop")
	}
}

func TestSpeakerInitializedOnFirstPlay(t *testing.T) {
	p := &Player{}
	if p.sampleRate != 0 {
		t.Error("expected sampleRate to be 0 before first play")
	}

	err := p.Play(testFile("ShadowsAndDust.mp3"))
	if err != nil {
		t.Skipf("skipping, unable to play file: %v", err)
	}
	defer p.Stop()

	if p.sampleRate == 0 {
		t.Error("expected sampleRate to be set after first play")
	}
}

func TestPlayM4AWithFFmpeg(t *testing.T) {
	if !IsFFmpegAvailable() {
		t.Skip("skipping, ffmpeg not available")
	}

	p := &Player{}
	err := p.Play(testFile("porcelain.m4a"))
	if err != nil {
		t.Errorf("expected no error playing m4a, got: %v", err)
	}
	defer p.Stop()
}
