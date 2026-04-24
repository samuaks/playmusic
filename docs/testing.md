# Testing Guide

## Goal

This guide explains what to run locally before a PR and how local verification differs from the current GitHub Actions setup.

## Default Local Check

Run the full local test suite:

```bash
go test ./...
```

This is the recommended default before opening a pull request.

Important: the full local suite includes `yt_dlp` integration-style tests. Those tests may require internet access and the ability to install or run external tooling.

## Also Run The App

For behavior changes, also run:

```bash
go run .
```

Manual verification matters for TUI and playback changes.

## Package-Level Test Expectations

| Package area | Typical coverage |
| --- | --- |
| `library` | loading, deduplication, background scanning, enrichment pipeline |
| `tui` | model updates, focus changes, list behavior, scan event handling |
| `player` | playback lifecycle, stopping, pause/resume, external video handoff |
| `decoder` and `ffmpeg` | decoder selection and FFmpeg probing helpers |
| `yt_dlp` | integration-oriented helpers around external tooling |

## CI vs Local

Current GitHub Actions behavior on pull requests:

- installs ALSA-related dependencies on Ubuntu;
- runs Go tests for most packages;
- currently excludes the `yt_dlp` package from the PR workflow.

That means:

- local `go test ./...` is still the most complete default command;
- if you change `yt_dlp`, you should run the relevant tests locally and mention the result in your PR.
- if you are offline or in a restricted environment, explain which subset you ran and why.

## Manual Checks Worth Doing

If your change affects behavior, try to verify the specific path directly:

- local audio playback from `Media/`;
- entering and leaving search mode;
- filter behavior while typing;
- background scan status updates;
- pause/resume and next-track behavior;
- random mode via `Ctrl+R`;
- `.mp4` handoff if you changed video-related behavior and have `ffplay`.

## When To Add Tests

Add or update tests when:

- you fix a bug with a reliable reproduction path;
- you change branching logic in `library`, `tui`, or `player`;
- you change startup or scan behavior;
- a regression would be easy to miss manually.

## When Manual-Only Verification Is Acceptable

Manual-only verification can be acceptable for:

- pure copy-only or metadata-only changes;
- small visual/help-text changes;
- exploratory work that is hard to test immediately.

If you skip tests, say so explicitly in the PR.

## If A Test Is Environment-Sensitive

Document the limitation instead of hiding it.

Good example:

```text
I ran go test ./...
The change touches yt_dlp behavior, so I also ran the package tests locally on Windows.
CI currently does not cover that package in the PR workflow.
```
