# Agent Guidelines

This repository is building a new programming language implementation in Go 1.26.

## Priority Order
1. Follow instructions in `.github/instructions/00-Bootstrap.instructions.md` when that file is populated.
2. Follow applicable `*.instructions.md` files in `.github/instructions/`.
3. Follow this file for repo-wide defaults.

## Repo-Wide Defaults
- Target platforms: Windows, Linux, and macOS.
- Prefer cross-platform path handling and avoid OS-specific assumptions.
- Keep behavior deterministic across operating systems.
- Favor readable, testable compiler and runtime code over clever shortcuts.
- For larger work, propose decomposition into sub-agents with clear output contracts.

## Go Constraints
- Language implementation: Go 1.26.
- Keep modules tidy and avoid unnecessary dependencies early in bootstrap.
- Add tests alongside parser, type system, and runtime changes.
