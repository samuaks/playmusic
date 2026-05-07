# Contributing To PlayMusic

Thanks for contributing to PlayMusic.

This project is intentionally beginner-friendly. Small pull requests, focused tests, and scoped bug fixes are all valuable contributions.

Project language is English. Please use English in issues, pull requests, review comments, and repository documentation.

## Before You Start

- Read [README.md](README.md) for the project goal and quick start.
- Read [docs/onboarding.md](docs/onboarding.md) if this is your first contribution here.
- Read [SUPPORT.md](SUPPORT.md) if you are unsure where to ask for help.

## CLI Repository vs Desktop Application

This repository is the CLI/TUI version of PlayMusic. It runs in the terminal and is mostly Go code.

The related desktop application is [samuaks/musical-palm-tree](https://github.com/samuaks/musical-palm-tree). It opens in a separate native window and is built with Tauri, React, and Rust.

For a longer comparison, read [docs/related-projects.md](docs/related-projects.md).

Use this repository for:

- terminal controls, TUI behavior, and local CLI playback;
- Go packages such as `library`, `tui`, `player`, `decoder`, `ffmpeg`, and `yt_dlp`;
- CLI contributor documentation, tests, and GitHub workflow improvements.

Use `samuaks/musical-palm-tree` for:

- desktop-window application behavior;
- React UI, waveform player UI, video viewport, and library view changes;
- Tauri/Rust scanner, waveform generation, IPC, and native app packaging.

If a task mentions the "application" part of PlayMusic, check whether it means the separate desktop app before opening an issue or PR here.

## Good First Contributions

These are great first PRs:
- add or improve a focused test in `library`, `tui`, or `player`;
- improve error messages, help text, or TUI copy;
- fix a well-scoped bug with a clear reproduction path;
- rewrite or clarify an issue so a new contributor can understand it.

These are usually not ideal first PRs:
- wide refactors across several packages;
- release workflow changes without prior discussion;
- exploratory work around unstable online or radio flows;
- "clean up everything" changes with no clear acceptance criteria.

If you are looking for a first task, prefer issues that are already labeled `good first issue`, `bug`, or have a clearly small scope. A good first issue should be understandable from the issue body without private context.

Before starting, check that the issue has:

- a clear problem or desired outcome;
- small scope, ideally one package or one user flow;
- acceptance criteria or obvious verification steps;
- starting points such as files, packages, or tests to read first.

If an issue is missing those details, improving the issue description is a good first contribution too. Ask a short clarifying question or propose a rewritten issue body before coding.

## Local Setup

Clone the repository and verify that the code builds and tests on your machine:

```bash
go test ./...
go run .
```

For more detailed test output, add `-v`:

```bash
go test -v ./...
go test -v ./library
go test -run TestName -v ./tui
```

Use `go test ./...` as the default pre-PR check. Use `-v` when you want to see each test name, debug a failure, or show clearer verification notes in a pull request.

If you want extended format support or `.mp4` handoff, install FFmpeg so `ffplay` is available in your `PATH`.

On the first app run, PlayMusic may also need internet access to resolve or install its `yt-dlp` helper binary.

## How To Pick A Task

- Start with issues labeled `good first issue` or clearly scoped `bug` tasks.
- If the issue is vague, comment first and ask for scope confirmation before coding.
- If you are new, prefer tasks that touch one package and have obvious verification steps.

For more guidance, read [docs/workflow.md](docs/workflow.md).

## How To Create An Issue

If the work does not already have an issue, create one before writing code.

Recommended flow:

1. search open issues first to avoid duplicates;
2. choose the matching GitHub template:
   `Bug report` for broken behavior or `Development task` for planned work;
3. write a short, specific title;
4. fill in the template so a new contributor could understand the task without private context;
5. only then use `Create a branch` from that issue.

When you open an issue, keep it to one problem or one outcome.

## Issue Size And Estimates

Issue estimates are rough planning hints, not promises.

Use small estimates for first-time contributor work:

| Estimate | Meaning |
| --- | --- |
| `1 point` | copy, docs, one focused test, or a tiny behavior fix |
| `2 points` | small change in one package with clear verification |
| `3 points` | one-package behavior change that needs tests and manual checking |
| `5 points` | larger work, multiple files, or unclear edge cases |
| `8 points` | split this before treating it as a first contribution |

If a task feels bigger than its estimate, stop and comment with what you found. Splitting the issue is better than growing a first PR until it becomes hard to review.

## Issue-First Workflow

For implementation work, the expected flow is:

1. create or confirm the GitHub issue;
2. make sure the issue title is short and descriptive;
3. use GitHub's `Create a branch` action from that issue;
4. open the pull request back to the same issue context.

Because the branch is created from the issue, issue naming matters.

Good issue titles are:

- short enough to scan quickly in the issue list;
- specific about the problem or missing behavior;
- written in clear English without internal shorthand.

Examples:

- `Fix search hint after leaving search mode`
- `Keep background scan items in arrival order`
- `Handle missing ffplay more clearly`

Avoid issue titles like:

- `Fix stuff`
- `UI bug`
- `Refactor`

If GitHub generates the branch from the issue, keep that default unless there is a strong reason to change it manually.

## Fork Workflow

If you do not have write access to this repository, use a fork.

Recommended fork flow:

1. fork `samuaks/playmusic` on GitHub;
2. clone your fork locally;
3. create a branch for one issue;
4. make the change and verify it;
5. push the branch to your fork;
6. open a pull request back to `samuaks/playmusic`.

If you are using GitHub's issue page, you may still use the issue title as your branch name. Keep the branch short and specific, for example `fix-search-hint` or `test-library-scan-missing-dir`.

## Coding Expectations

- Keep changes focused and easy to review.
- Match the existing style of the package you touch.
- Run `gofmt` on modified Go files.
- Add or update tests when behavior changes.
- Update docs when user-facing behavior changes.
- Do not document experimental features as stable.

## Pull Request Checklist

Before opening a PR, make sure you can say "yes" to most of these:

- I can explain the problem in one or two sentences.
- I kept the scope intentionally small.
- I ran `go test ./...` locally, or I explained why I could not.
- I manually checked the relevant user flow if the change affects behavior.
- I updated docs or help text if the change is user-facing.
- I described what is intentionally out of scope.

## Commit Messages

There is no strict format requirement, but clear commit messages help reviewers.

Clear means that a reviewer can understand the change without opening the diff first. Prefer a short imperative sentence that names the area and the outcome.

Prefer short imperative messages such as:

- `fix tui search hint after focus change`
- `test library scan skips missing dirs`
- `fix player status text after track switch`

Good commit messages usually:

- start with the action: `fix`, `test`, `docs`, `add`, `update`, or `remove`;
- mention the package or user-facing area when helpful;
- describe one meaningful change, not a whole work session;
- avoid private shorthand that only the author understands.

Avoid vague messages such as:

- `fix`
- `changes`
- `update`
- `work in progress`
- `stuff`

## Reviews And Merge Flow

- Open the PR early if you want feedback on direction.
- Keep discussion in the PR so future contributors can learn from it.
- Expect review comments to improve clarity, scope, and maintainability, not to gatekeep.
- If something is unclear, ask directly in the PR instead of guessing.

## Reporting Blockers

If you are blocked:

1. Explain what you tried.
2. Mention the file or issue you were working on.
3. Say what kind of help you need: scope, code direction, test help, or review.

Use [SUPPORT.md](SUPPORT.md) for the right channel.
