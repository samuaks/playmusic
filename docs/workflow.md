# Contributor Workflow

## Purpose

This guide explains how to choose a task, keep the scope healthy, and know what "done" means in PlayMusic.

## How To Choose A Task

Pick work that is:

- understandable from the issue alone, or after one short clarification;
- limited to one package or one clear user flow;
- easy to verify with a test or a short manual check.

Strong first choices:

- `good first issue`
- small `bug`
- focused `enhancement`

If labels are missing, look for issues with:

- a clear problem statement;
- concrete acceptance criteria;
- hints about which files to read first.

## Suggested Difficulty Scale

| Level | Typical scope | Typical time |
| --- | --- | --- |
| `easy` | tests, copy, one small behavior fix | 30 min to 2 h |
| `medium` | one-package logic change with tests | half day to 1 day |
| `hard` | cross-package behavior, concurrency, release, or architecture work | more than 1 day |

## What A Good Issue Should Contain

A contributor-friendly issue should include:

- a short, descriptive title that still reads well when GitHub turns it into a branch name;
- context: why the work matters;
- problem statement: what is wrong today;
- scope: what is in and out;
- acceptance criteria: how we know it is done;
- starting points: files or packages to read first;
- level and rough time estimate.

If an issue is missing these pieces, improving the issue first is a valid contribution.

## How To Open A Good Issue

Use GitHub issues as the starting point for implementation work.

Recommended flow:

1. search existing issues first so you do not create a duplicate;
2. choose the right template:
   `Bug report` for broken behavior or `Development task` for planned work;
3. write one issue for one concrete problem or outcome;
4. keep the title short and specific;
5. fill the body so someone outside the core team could understand the task;
6. once the issue is ready, use `Create a branch` from that issue.

The issue should be good enough that another contributor could pick it up later without needing hidden context.

## Definition Of Done

A task is usually done when:

- the intended behavior works;
- relevant tests pass, or the missing coverage is explained honestly;
- user-facing changes have matching docs or help text updates;
- the PR explains verification steps;
- the change stays within the agreed scope;
- known limitations are called out instead of hidden.

## Good First Contribution Examples

- add a regression test in `library` for scan or dedup logic;
- improve a small `tui` help string or status message;
- fix a narrow playback state bug with a reliable reproduction path;
- rewrite a vague issue using the project template.

## Usually Not A First Contribution

- major refactors with no user-facing acceptance criteria;
- release or workflow changes without maintainer agreement;
- experimental online or radio flows that are still moving quickly;
- cross-package concurrency changes unless the issue is tightly scoped.

## Preferred Task Flow

1. Confirm the issue still matches the current code.
2. Use `Create a branch` from the GitHub issue so the work stays linked to that task.
3. Make the smallest change that solves the stated problem.
4. Verify it locally.
5. Open a PR with a short summary and verification notes.

## Issue Title Guidance

Since branches are created from issues in the current workflow, issue titles should be clean and specific.

Prefer titles that:

- describe one concrete problem or outcome;
- avoid vague words like `stuff`, `cleanup`, or `refactor` on their own;
- stay short enough to read in lists and pull request references.

Good examples:

- `Fix panic when pausing during radio track switch`
- `Show clearer empty-library message on startup`
- `Keep selected track when background scan finishes`

Weak examples:

- `Bug in player`
- `TUI issue`
- `Refactor code`

## If The Task Is Bigger Than Expected

Do not force the original scope.

Instead:

- stop and describe what you learned;
- split the task into smaller follow-ups;
- ask for scope confirmation before continuing.

That is good collaboration, not failure.
