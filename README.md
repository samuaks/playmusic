# Go Media Player

A modular command-line media player written in Go. Uses [beep](https://github.com/gopxl/beep) for native audio playback with optional FFmpeg support for extended format compatibility.

## Features
- Plays audio files from a `Media` directory automatically
- Native support for mp3, flac, wav, and ogg without any external dependencies
- Extended format support (m4a, aac, opus) when FFmpeg is installed
- Displays track listing with durations on startup
- Modular architecture split across `library`, `player`, and `colors` packages
- Search and stream songs from the web

## Supported Formats

| Format     | Requires FFmpeg |
|------------|----------------|
| .mp3       | No |
| .flac      | No |
| .wav       | No |
| .ogg       | No |
| .m4a       | Yes |
| .aac       | Yes |
| .opus      | Yes |
| streaming  | Yes |


## Roadmap

| # | Description | Type | Status |
|---|-------------|------|--------|
| 1 | Bubble Tea TUI implementation | Feature | Partial |
| 2 | Pause / Resume / Next track support | Feature | Done |
| 3 | Volume control | Feature | Not Planned |
| 4 | Custom media directory / system-wide media scanning | Feature | Partial |
| 5 | Headphone wear detection (auto-pause on removal) | Feature | Research |
| 6 | Sample rate mismatch on some tracks | Bug | Known/Fixed |
| 7 | CI/CD — GitHub Actions for running tests and building release executables | Feature | Planned |
| 8 | Album / playlist support | Feature | Planned |
| 9 | Artist fetching (background job) | Feature | Planned |
| 10 | Play songs from external sources (YouTube) | Feature | Done |



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
Player will also search for the music in common places in your PC and will add them to the playlist.

## Hotkeys

- `up/down` - navigate list
- `space` - pause/resume
- `enter` in list focus - play selected track
- `q` or `?` - enter search focus
- `enter` in search focus - run external search (if query is not empty), then return to list focus
- `esc` in search focus - clear query and return to list focus

## Requirements

- Go 1.21 or later
- FFmpeg (optional) — install and add to PATH for extended format support. For Windows all can be done by running command: winget install ffmpeg. 
- yt-dlp will be installed to the temporary directory on your PC and will be used for the search and stream music feature.

## Attribution
Music used in demo:
>'Shadows and Dust' by Scott Buckley – released under CC-BY 4.0. www.scottbuckley.com.au
>'Rites of passage' by Scott Buckley – released under CC-BY 4.0. www.scottbuckley.com.au



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

### Bubbletea event loop (main goroutine)

`Update(), View()` Handles messages and updates the view

`Player goroutine` Plays music

`tea.Cmd goroutine` Searcher.Search() blocks here, not in main loop 


### Ideal approach on startup

1. LoadLibrary("Media") - blocking 
2. Start TUI with immediate local tracks
3. `go ScanSystem()` background job, sends new tracks to TUI via channels
    * `newTrackMsg` everytime new track was found 
    * TUI appends it to the list

    