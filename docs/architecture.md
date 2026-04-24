# Architecture Overview

## Purpose

This document explains the current architecture of PlayMusic in a way that helps a new contributor navigate the codebase quickly.

## The Short Version

PlayMusic is a Go application with a Bubble Tea TUI on top of a local media library and playback layer.

The default user path today is:

```text
main.go
  -> load local Media/ synchronously
  -> start background scan for more library roots
  -> start Bubble Tea TUI
  -> select a track
  -> player + decoder handle playback
```

## Package Map

| Package | Responsibility |
| --- | --- |
| `main` | Startup flow and wiring |
| `library` | Track discovery, enrichment, deduplication, scan events |
| `tui` | Application state, keyboard handling, rendering, scan-status updates |
| `player` | Playback lifecycle, pause/resume/next, external video handoff |
| `decoder` | Native decoding through `beep`, FFmpeg-backed decode paths, stream decode helpers |
| `ffmpeg` | FFmpeg adapter functions and probing helpers |
| `search` | Search abstractions and sources used by evolving external-search flows |
| `yt_dlp` | `yt-dlp` integration helpers used by external media flows |
| `helpers` | Shared helper functions such as duration formatting |
| `colors` | Styling constants used by the UI |

## Core Data Shape

The main shared data structure is `library.Track`.

It carries:

- identity such as `Path` or `YTVideoURl`;
- display fields such as `Trackname`, `Artist`, and `Title`;
- metadata such as `Duration`, `Album`, `Year`, and `Genre`.

If you are trying to understand how data moves through the app, start with this type.

## Startup Flow

The real startup path lives in `main.go`.

Today it does this:

1. installs or resolves the `yt-dlp` binary;
2. loads the local `Media/` directory synchronously;
3. starts a background scan for additional library directories;
4. passes the initial tracks and scan channel into the TUI model;
5. runs the Bubble Tea program.

This means the app is optimized to show something quickly instead of waiting for a full library scan.

## Local Library Flow

`library/library.go` handles the initial load.

Important ideas:

- `LoadLibrary("Media")` is the fast startup path;
- `DefaultLibraryDirs()` defines the broader library roots;
- `BackgroundLibraryDirs()` excludes `Media/` so startup does not duplicate work;
- deduplication tries to prevent duplicate tracks across startup and background scans.

## Background Scan Flow

`library/scan.go` emits `ScanEvent` values over a channel.

The TUI listens for:

- discovered tracks;
- enriched tracks;
- non-fatal scan errors;
- scan completion.

This lets the UI stay responsive while the rest of the library is still being discovered.

## Playback Flow

`player/player.go` owns playback state.

For audio files:

1. `player.Play()` asks `decoder.Decode()` for a streamer;
2. `decoder` chooses a native `beep` decoder or an FFmpeg-backed path;
3. `player` sends the streamer to `speaker`;
4. `Wait()`, `Pause()`, `Resume()`, and `Next()` control playback lifecycle.

For `.mp4` files:

- PlayMusic does not render video inside the TUI;
- it hands playback off to `ffplay` through `player.playVideo()`.

## TUI Flow

The TUI lives mostly in:

- `tui/model.go`
- `tui/update.go`
- `tui/view.go`

Key ideas:

- the model keeps track list, selected/current index, elapsed time, focus mode, and scan state;
- the UI has two explicit focus modes: list and search;
- background scan events update the model incrementally;
- the player bar and search bar reflect current state instead of driving business logic.

## Where To Start Based On The Change

- Startup or scan roots: `main.go`, `library/library.go`, `library/scan.go`
- Track metadata or dedup: `library/*`
- Keyboard behavior or layout: `tui/*`
- Audio or video playback: `player/player.go`, `decoder/decoder.go`, `ffmpeg/*`
- External-search plumbing: `search/*`, `yt_dlp/*`

## Areas To Treat As Experimental

The default contributor path is local library plus TUI plus playback.

External-search, radio, and online-flow code exists in the repository history and supporting packages, but those areas are still evolving and should be documented carefully.
