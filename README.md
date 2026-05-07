# PlayMusic

PlayMusic is a terminal media player written in Go and a student-friendly open-source playground for learning collaborative development.

The project is meant to be useful in two ways:
- as a real application you can run, inspect, and improve;
- as a safe place for Hive students/sprinters and first-time contributors to practice reading code, discussing scope, writing tests, and opening pull requests.

Project language is English for repository docs, issues, pull requests, and other contributor-facing materials.

## Related Application

This repository is the CLI/TUI version of PlayMusic. It runs in the terminal and is focused on local playback, keyboard-driven navigation, and Go contributor practice.

The related desktop application lives at [samuaks/musical-palm-tree](https://github.com/samuaks/musical-palm-tree). That project opens in a separate native window and is built with Tauri, React, and Rust. Use that repository for GUI application work such as the windowed library view, waveform player UI, video viewport, desktop packaging, and Tauri-specific behavior.

For a longer comparison, read [docs/related-projects.md](docs/related-projects.md).

When opening an issue or pull request, choose the repository based on the user surface:

- terminal behavior, Go packages, CLI/TUI controls, or this contributor workflow: use this repository;
- desktop window behavior, React UI, Rust/Tauri scanner or waveform code, or native app packaging: use `samuaks/musical-palm-tree`.

## Start Here

- Want to run the app quickly? Start with [Quick Start](#quick-start).
- Want to make your first contribution? Read [docs/onboarding.md](docs/onboarding.md).
- Need contribution rules? Read [CONTRIBUTING.md](CONTRIBUTING.md).
- Need help or response expectations? Read [SUPPORT.md](SUPPORT.md).

## What You Can Do In Your First Hour

1. Run the app against the local `Media/` folder.
2. Learn how `main`, `library`, `tui`, `player`, and `decoder` fit together.
3. Pick a small test, bug, or focused UI issue.
4. Open a first pull request without needing private project context.

## Current Status

### Working Today

- Bubble Tea TUI for browsing and playing a local music library.
- Audio playback through `beep` for `.mp3`, `.flac`, `.wav`, and `.ogg`.
- FFmpeg-backed support for additional audio formats such as `.m4a`, `.aac`, and `.opus` when FFmpeg is installed.
- Fast startup from the local `Media/` folder, followed by background scanning of other library directories.
- Live local filtering with an explicit search mode opened by `q` or `?`.
- Random next-track mode toggled with `Ctrl+R`.
- External video handoff for `.mp4` files through `ffplay`.
- Go test coverage across the core packages plus GitHub Actions for CI and releases.

### Still Evolving

- Online, radio, and external-search flows are actively evolving and should not be treated as the default contributor path.
- If something in the contributor flow is unclear, open a regular issue and describe the gap in context.

## Quick Start

### Requirements

- Go `1.25.5` or newer.
- FFmpeg if you want extended audio-format support and external `.mp4` playback through `ffplay`.
- Local media files in the repository `Media/` folder if you want immediate content on startup.
- Internet access on the first run may be needed so the app can resolve or install its `yt-dlp` helper binary.

### Run The Full Test Suite

```bash
go test ./...
```

Note: the full local suite includes `yt_dlp` integration-style tests, so it may require internet access and external tooling on your machine.

### Start The App

```bash
go run .
```

You can also build the binary first:

```bash
go build
./playmusic
```

On Windows the binary name will be `playmusic.exe`.

## Controls

| Key | Action |
| --- | --- |
| `Up` / `Down` | Move through the list |
| `Enter` | Play the selected track |
| `Space` | Pause or resume audio playback |
| `q` or `?` | Enter search mode |
| `Enter` in search mode | Keep the current filter and return to list mode |
| `Esc` in search mode | Clear the filter and return to list mode |
| `Ctrl+R` | Toggle random next-track mode |
| `Ctrl+Q` or `Ctrl+C` | Quit |

For `.mp4` files, PlayMusic hands playback to `ffplay`, so playback controls happen in the external player instead of inside the TUI.

## Project Map

- [docs/onboarding.md](docs/onboarding.md): guided path from clone to first PR in 30-60 minutes.
- [docs/workflow.md](docs/workflow.md): how to choose a task, keep scope small, and know what "done" means.
- [docs/architecture.md](docs/architecture.md): package map and main application flows.
- [docs/testing.md](docs/testing.md): local testing, CI expectations, and manual verification tips.
- [docs/related-projects.md](docs/related-projects.md): how this CLI/TUI repository relates to the separate desktop application.

## Community

- [CONTRIBUTING.md](CONTRIBUTING.md)
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- [SUPPORT.md](SUPPORT.md)
- [SECURITY.md](SECURITY.md)
- [LICENSE](LICENSE)

## Known Notes For Contributors

- The primary contributor path is local library playback and TUI behavior.
- The current GitHub Actions PR workflow excludes the `yt_dlp` package, so local `go test ./...` remains the best pre-PR verification command.
- `config.go` contains path-resolution helpers that are not yet wired into the current `main.go` flow, so do not document them as a supported user-facing configuration system.

## Attribution

Demo music in the repository is attributed in the existing project materials and should remain credited when reused for demos.
