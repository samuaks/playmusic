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



## Concurrency

The player uses goroutines and channels to manage concurrent operations and user input

### WaitGroup() - do many things at once and wait for them ALL to finish before continuing.

Used in library package to probe all track durations concurrently, improving startup time significantly when the number of tracks is large. 
Without this the durations would be probed sequentially, in which time would scale linearly with the number of tracks.

Visual example of the concept:

```
|-main: [wg.Wait()... blocking until all done...]
|-goroutine 1: [ProbeDuration(track1)...]
|-goroutine 2: [ProbeDuration(track2)...]
|-goroutine 3: [ProbeDuration(track3)...]
```

### Non-blocking user input handling with goroutines

Used for input handling in main.go while tracks are playing.
Without a goroutine, main thread could either play music OR handle user input, not both.

```go
go handleInput(player) // runs independently in the background while main thread continues to play music

for _, track := range tracks {
    player.Play(track) 
    player.Wait() // unblocks next track when song naturally finishes OR when user presses a key to skip 
}
```