package ffmpeg

type FFmpegInterface interface {
	Probe(path string) (string, error)
}
