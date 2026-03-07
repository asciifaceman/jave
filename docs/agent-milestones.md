# Agent Milestones and Acceptance Criteria

This is the implementation plan for Jave v0.1.

## Milestone 1: Vertical Slice

Scope:
- lexer
- parser
- AST
- diagnostics
- execution path sufficient for core examples
- built-ins: `pront`, `girth`
- control flow and both `given` loop forms
- base types and core collections

Acceptance criteria:
- Parse and run `examples/hello_world`.
- Parse and run conditional and loop examples.
- Nested `table<table<exact>>` indexing works.
- Diagnostics include line/column with stable wording.

## Milestone 2: Imports and Toolchain Shape

Scope:
- import resolution for `highschool/...`
- `Foreward` semantics
- `Strangs.Combobulate`
- `Pronts.Prontulate`
- `javec`, `baggage`, `javevm` CLI command growth
- `.jbin` format defined
- syntax highlighting support with VS Code priority and GitHub compatibility

Acceptance criteria:
- Examples in `examples/imports` and `examples/combobulate` run.
- Legacy alias `Srangs` emits warning and still resolves.
- VS Code highlighting is available via a Jave grammar extension in-repo.
- GitHub highlighting path is documented and tracked (Linguist strategy and fallback docs).

## Milestone 2.1: Syntax Highlighting (Priority Track)

Scope:
- implement a VS Code TextMate grammar for `.jave`
- ship a minimal VS Code extension scaffold under repo tooling
- document install/test workflow for local highlighting development
- define GitHub highlighting plan (Linguist grammar upstream path and interim fallback)

Acceptance criteria:
- Opening `.jave` files in VS Code highlights keywords, types, literals, operators, and comments.
- A contributor can install and test the highlighting extension from this repo in under 5 minutes.
- GitHub strategy is recorded with next actions and ownership.

## Milestone 3: Sponsor Messaging and Polish

Scope:
- sponsor notice subsystem
- suppression flags and partial redaction behavior
- docs and diagnostic polish

Acceptance criteria:
- Sponsor output behavior is deterministic.
- Suppression flags match spec behavior.
- Required docs are present and coherent.

## Role Split (Sub-agent Friendly)

- `01-Compiler-Architect`: phase boundaries, IR/jbin shape
- `02-Parser-and-Diagnostics`: grammar, recovery, errors
- `03-Runtime-and-Execution`: evaluator/vm execution behavior
- `04-Tooling-and-Release`: CLI, test matrix, CI
