# Contributing To PlayMusic

Thanks for contributing to PlayMusic.

This project is intentionally beginner-friendly. Small pull requests, focused tests, and scoped bug fixes are all valuable contributions.

Project language is English. Please use English in issues, pull requests, review comments, and repository documentation.

## Before You Start

- Read [README.md](README.md) for the project goal and quick start.
- Read [docs/onboarding.md](docs/onboarding.md) if this is your first contribution here.
- Read [SUPPORT.md](SUPPORT.md) if you are unsure where to ask for help.

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

## Local Setup

Clone the repository and verify that the code builds and tests on your machine:

```bash
go test ./...
go run .
```

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

Prefer short imperative messages such as:

- `fix tui search hint after focus change`
- `test library scan skips missing dirs`
- `fix player status text after track switch`

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
