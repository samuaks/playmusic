# Onboarding

## Why Read This

This guide is for a new contributor who wants to go from clone to first pull request in about 30-60 minutes.

You do not need deep Go or open-source experience to get value from this project.

## Outcome

By the end of this guide, you should be able to:

- run the project locally;
- understand the main files worth reading first;
- choose a small, realistic task;
- open a first PR with the right level of scope.

## Step 1: Clone And Verify

Run the basic checks first:

```bash
go test ./...
go run .
```

If the app starts, you already have a working local environment.

Two practical notes:

- the app may need internet access on first run so it can resolve or install `yt-dlp`;
- the full test suite includes `yt_dlp` integration-style tests, so offline environments may need extra care.

## Step 2: Put Something In `Media/`

The default startup flow loads the local `Media/` directory first.

- Put one or more audio files into `Media/`.
- Start the app again if needed.
- Confirm that you can see the track list and play a file.

## Step 3: Learn The Five Files That Explain The Project

Read these files in order:

| File | Why it matters |
| --- | --- |
| `main.go` | Shows the real startup flow and what the default app path is today |
| `library/library.go` | Explains local library loading and default scan roots |
| `library/scan.go` | Explains background scan events and incremental updates |
| `tui/model.go` and `tui/update.go` | Explain the UI state and keyboard behavior |
| `player/player.go` | Explains playback lifecycle, pause/resume, and `.mp4` handoff |

If you want one extra file after that, read `decoder/decoder.go`.

## Step 4: Pick A Beginner-Friendly Task

Good first tasks are usually:

- a focused test addition;
- a small TUI text or UX fix;
- a narrow bug with clear reproduction;
- an issue rewrite that adds context and acceptance criteria.

If labels are incomplete or inconsistent, use common sense:

- prefer one-package changes;
- avoid vague refactors;
- ask in the issue before taking on something large.

## Step 5: Keep Scope Intentionally Small

A strong first PR usually changes one of these:

- one test file;
- one behavior inside one package;
- one user-facing string or help message.

You do not need to "prove yourself" with a large change.

## Step 6: Open The PR

Before opening the PR:

```bash
go test ./...
```

Then explain:

- what you changed;
- why it matters;
- how you verified it;
- what you deliberately did not change.

Use [CONTRIBUTING.md](../CONTRIBUTING.md) and the PR template for the detailed workflow.

## If You Get Stuck

- Ask in the issue if the task scope is still unclear.
- Ask in the PR if you already started work.
- Use [SUPPORT.md](../SUPPORT.md) when you need the right channel or response expectations.

## Suggested First-Week Ladder

Try to move through tasks like this:

1. issue cleanup
2. focused test
3. small TUI fix
4. small bug fix
5. one-package feature improvement

That progression gives you context without making your first week stressful.
