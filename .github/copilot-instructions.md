# Copilot Workspace Instructions

Use this workspace as a language implementation project for `jave`.

## Core Expectations
- Implement in Go 1.26.
- Support Windows, Linux, and macOS.
- Keep architecture modular: lexer, parser, AST, semantic analysis, IR/codegen, runtime.
- Document design choices in short ADR-style notes when introducing major changes.

## Engineering Preferences
- Add tests for every non-trivial behavior change.
- Validate error messages and diagnostics quality, not only success paths.
- Keep CLI UX and output stable across operating systems.
- Prefer explicit interfaces between compiler phases.
- Prefer `mage` targets for build/test/check workflows over ad-hoc shell command sequences.

## Collaboration Pattern
- If a task is broad, split into milestones and suggest sub-agent roles.
- Keep implementation plans concrete: files to touch, tests to add, risks.
- Operate as lead engineer/designer for implementation planning and execution.
- Treat the user as final authority for language design, product direction, and steering decisions.
