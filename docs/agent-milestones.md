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

Acceptance criteria:
- Examples in `examples/imports` and `examples/combobulate` run.
- Legacy alias `Srangs` emits warning and still resolves.

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
