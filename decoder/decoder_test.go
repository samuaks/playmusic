package decoder

import "testing"

func TestIsSupported(t *testing.T) {
	cases := []struct {
		filename   string
		want       bool
		needFFmpeg bool
	}{
		{"song.mp3", true, false},
		{"track.wav", true, false},
		{"audio.flac", true, false},
		{"song.MP3", true, false},
		{"document.txt", false, false},
		{"image.jpg", false, false},
		{"song.m4a", true, true},
		{"song.aac", true, true},
		{"song.opus", true, true},
	}
	for _, c := range cases {
		if c.needFFmpeg && !IsFFmpegAvailable() {
			t.Logf("skipping %q test, ffmpeg not available", c.filename)
			continue
		}
		got := IsSupported(c.filename)
		if got != c.want {
			t.Errorf("IsSupported(%q) = %v, want %v", c.filename, got, c.want)
		}
	}
}
