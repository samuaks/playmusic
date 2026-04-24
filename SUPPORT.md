# Support

## Where To Ask For Help

At the moment, the official support path for PlayMusic is GitHub-first.

- Use GitHub issues for bugs, task proposals, and contributor-flow questions.
- Use pull request comments for implementation-specific questions on an active change.
- If a future Discord or Notion space becomes part of the contributor flow, this file should be updated to link it.

## Which Channel To Use

- Bug report: open a bug issue.
- Feature or improvement proposal: open a task or enhancement issue.
- Contributor-flow gap or unclear project information: open a regular issue with the missing context.
- Question about your current PR: comment in the PR.
- Security concern: follow [SECURITY.md](SECURITY.md).

## Response Expectations

These are targets, not guarantees:

- first response on new issues: within 3 working days;
- first review pass on small PRs: within 5 working days;
- clarification on blocked student-first tasks: as soon as a maintainer is available.

If a task is time-sensitive for onboarding or coursework, say that clearly in the issue or PR.

## How To Ask A Good Question

Good questions usually include:

- what you are trying to do;
- what file, issue, or command you are working with;
- what you expected to happen;
- what actually happened;
- what you already tried.

Example:

```text
I am working on issue #123 in tui/update.go.
I expected Enter in search mode to keep the filter and return to list mode.
Instead the filter resets after Esc.
I ran go test ./tui and checked the focus-related tests, but I am not sure where the state reset should live.
```

## Project Language

English is the project language for repository docs, issues, pull requests, and contributor support requests.

## If You Are Completely New

Start here:

1. [README.md](README.md)
2. [docs/onboarding.md](docs/onboarding.md)
3. [CONTRIBUTING.md](CONTRIBUTING.md)
