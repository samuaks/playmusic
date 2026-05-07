# Related PlayMusic Projects

PlayMusic currently has two related repositories with different user surfaces.

## CLI/TUI

Repository: [samuaks/playmusic](https://github.com/samuaks/playmusic)

This repository is the terminal version of PlayMusic. It is written in Go and focuses on:

- local media library loading and scanning;
- terminal keyboard controls;
- Bubble Tea TUI behavior;
- playback through Go packages and external helpers;
- beginner-friendly Go contribution practice.

Open issues and pull requests here when the change affects terminal behavior, Go packages, CLI testing, or this contributor workflow.

## Desktop Application

Repository: [samuaks/musical-palm-tree](https://github.com/samuaks/musical-palm-tree)

This is the separate desktop application. It opens in a native window and is built with Tauri, React, and Rust. It focuses on:

- windowed library browsing;
- React UI and desktop interaction;
- waveform player UI;
- video playback viewport;
- Tauri IPC between the frontend and Rust backend;
- native desktop builds and packaging.

Open issues and pull requests there when the change affects the windowed application experience, React components, Tauri commands, Rust scanner or waveform code, or desktop packaging.

## Choosing The Right Repository

If the behavior happens in the terminal, use `samuaks/playmusic`.

If the behavior happens in a separate desktop window, use `samuaks/musical-palm-tree`.

If a task says "application part" or "app window", confirm whether it refers to the desktop application before starting implementation in this repository.
