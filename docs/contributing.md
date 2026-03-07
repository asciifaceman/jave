# Contributing

Thanks for helping build Jave.

## Priorities

- Correct semantics over optimization.
- Stable diagnostics with clear line/column info.
- Cross-platform behavior across Windows, Linux, and macOS.
- Tests for non-trivial behavior changes.

## Workflow

- Read `specs/jave-v0.1.md` first.
- Keep examples in sync with syntax and behavior.
- Propose architecture changes with short ADR-style notes.
- New and reopened GitHub issues are auto-labeled and acknowledged by `.github/workflows/issue-triage.yml`.
- New issues are marked for Copilot/agent first-pass by `.github/workflows/issue-copilot-routing.yml`.
- Set repository variable `COPILOT_ASSIGNEE` to auto-assign a specific account for first-pass issue handling.
- Linked pull requests automatically post/update status back on referenced issues via `.github/workflows/issue-pr-linkback.yml`.
- Include `Fixes #<issue>` or `Refs #<issue>` in PR title/body so linkback can find the triggering issue.

## Decision Authority

Lead engineering/design execution is delegated to the agent workflow, and final language/design decisions remain with the project owner.
