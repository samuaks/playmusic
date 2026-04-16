package player

import (
	"os"
	"os/exec"
	"path/filepath"
	. "playmusic/decoder"
	"runtime"
	"strconv"
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

func TestPauseWithoutPlay(t *testing.T) {
	p := &Player{}
	p.Pause()
}

func TestResumeWithoutPlay(t *testing.T) {
	p := &Player{}
	p.Resume()
}

func TestNextWithoutPlay(t *testing.T) {
	p := &Player{}
	p.Next()
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
	err := p.Play(testFile("sample3.m4a"))
	if err != nil {
		t.Errorf("expected no error playing m4a, got: %v", err)
	}
	defer p.Stop()
}

func TestPlayMP4UsesExternalVideoPlayer(t *testing.T) {
	originalLauncher := launchVideoPlayer
	defer func() { launchVideoPlayer = originalLauncher }()

	var called bool
	var gotPath string
	launchVideoPlayer = func(path string) (*exec.Cmd, error) {
		called = true
		gotPath = path
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperExternalPlayerProcess")
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_EXTERNAL_PLAYER=1",
			"GO_HELPER_SLEEP_MS=20",
		)
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		return cmd, nil
	}

	p := &Player{}
	if err := p.Play("sample.mp4"); err != nil {
		t.Fatalf("expected mp4 playback to start, got error: %v", err)
	}
	defer p.Stop()

	if !called {
		t.Fatal("expected external video launcher to be used for mp4")
	}
	if gotPath != "sample.mp4" {
		t.Fatalf("expected launcher to receive sample.mp4, got %q", gotPath)
	}
	if p.ctrl != nil {
		t.Fatal("expected beep ctrl to stay nil for external video playback")
	}

	select {
	case <-p.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("expected external mp4 playback to eventually complete")
	}
}

func TestStopEndsExternalVideoPlayback(t *testing.T) {
	originalLauncher := launchVideoPlayer
	defer func() { launchVideoPlayer = originalLauncher }()

	launchVideoPlayer = func(path string) (*exec.Cmd, error) {
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperExternalPlayerProcess")
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_EXTERNAL_PLAYER=1",
			"GO_HELPER_SLEEP_MS=5000",
		)
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		return cmd, nil
	}

	p := &Player{}
	if err := p.Play("sample.mp4"); err != nil {
		t.Fatalf("expected mp4 playback to start, got error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	p.Stop()

	select {
	case <-p.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("expected stop to signal playback completion for external player")
	}
	if p.externalCmd != nil {
		t.Fatal("expected external command to be cleared after stop")
	}
}

func TestHelperExternalPlayerProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_EXTERNAL_PLAYER") != "1" {
		return
	}

	sleepMs, _ := strconv.Atoi(os.Getenv("GO_HELPER_SLEEP_MS"))
	if sleepMs <= 0 {
		sleepMs = 1000
	}
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	os.Exit(0)
}

func TestPauseResume(t *testing.T) {
	p := &Player{}
	err := p.Play(testFile("ShadowsAndDust.mp3"))
	if err != nil {
		t.Skipf("skipping, unable to play file: %v", err)
	}
	defer p.Stop()

	p.Pause()
	if !p.ctrl.Paused {
		t.Error("Expected to be paused")
	}
	p.Resume()
	if p.ctrl.Paused {
		t.Error("Expected to be resumed")
	}
}

func TestNext(t *testing.T) {
	p := &Player{}
	err := p.Play(testFile("ShadowsAndDust.mp3"))
	if err != nil {
		t.Skipf("skipping, unable to play file: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	p.Next()

	// anon func to wait for done / next
	done := make(chan struct{})
	go func() {
		p.Wait()
		close(done)
	}()

	// check channel if done
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("Expeceted Wait() to return after Next()")
	}
}
