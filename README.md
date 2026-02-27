# Go Media Player

A modular command-line media player written in Go. Uses [beep](https://github.com/gopxl/beep) for native audio playback with optional FFmpeg support for extended format compatibility.

## Features
- Plays audio files from a `Media` directory automatically
- Native support for mp3, flac, wav, and ogg without any external dependencies
- Extended format support (m4a, aac, opus) when FFmpeg is installed
- Displays track listing with durations on startup
- Modular architecture split across `library`, `player`, and `colors` packages

## Supported Formats

| Format | Requires FFmpeg |
|--------|----------------|
| .mp3   | No |
| .flac  | No |
| .wav   | No |
| .ogg   | No |
| .m4a   | Yes |
| .aac   | Yes |
| .opus  | Yes |


## Roadmap

| # | Description | Type | Status |
|---|-------------|------|--------|
| 1 | Bubble Tea TUI implementation | Feature | Planned |
| 2 | Pause / Resume / Next track support | Feature | In Progress |
| 3 | Volume control | Feature | Planned |
| 4 | Custom media directory / system-wide media scanning | Feature | Planned |
| 5 | Headphone wear detection (auto-pause on removal) | Feature | Research |
| 6 | Sample rate mismatch on some tracks | Bug | Known/Fixed |
| 7 | CI/CD — GitHub Actions for running tests and building release executables | Feature | Planned |


## Usage
```bash
go run .
```

or build and run the executable:
```bash
go build
./playmusic
```

Drop your media files into the `Media` directory and they will be played in order.

## Requirements

- Go 1.21 or later
- FFmpeg (optional) — install and add to PATH for extended format support

## Attribution
Music used in demo:
>'Shadows and Dust' by Scott Buckley – released under CC-BY 4.0. www.scottbuckley.com.au