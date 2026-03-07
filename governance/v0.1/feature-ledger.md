# v0.1 Feature Ledger

Status: Ratified and implemented

## Language and Runtime

- Entrypoints: `Foreward`, `Foremost`
- Primitive types: `exact`, `vag`, `truther`, `strang`, `nada`, `naw`
- Control flow: `maybe`/`furthermore`/`otherwise`, `given` loop modes
- Collections: `table`, `enumeration`, `lexis` (including nested tables)
- String assembly: `Strangs.Combobulate`, `Pronts.Prontulate`, `Pront`
- Imports: `install ... from ...`, highschool stdlib paths, cycle detection
- Diagnostics: Multi-file source attribution and deterministic formatting
- Runtime: Deterministic execution for lowered IR and `.jbin` artifacts

## Tooling and UX

- Compiler: `javec`
- Build manager: `baggage`
- Runtime: `javevm`
- Sponsor notice modes: `full`, `redacted`, `off`
- VS Code syntax highlighting extension scaffolding

## Late v0.1 Enablements

- Parameterized sequence calls are now executable in runtime:
- Parser retains sequence parameters in AST.
- Semantic analysis binds parameters into sequence scope.
- Runtime supports user-defined sequence invocation with argument binding.
- Arity mismatch diagnostics for sequence calls.

## Acceptance Trace

- Parser tests include parameter parsing and advanced examples.
- Semantic tests include parameter binding and arity mismatch checks.
- Runtime tests include parameterized sequence execution.
- Full suite passes: `go test ./... -count=1`.
