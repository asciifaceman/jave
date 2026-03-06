# Grammar Notes (v0.1)

This file tracks grammar notes for parser implementation planning.

## Core v0.1 Notes

- Statement terminator is always `;;`.
- Declarations begin with `allow`.
- Sequences support `sequence` and alias `seq`.
- Conditional chain supports `maybe` -> zero or more `furthermore` -> optional `otherwise`.
- Looping requires both `given (<cond>) again -> { ... }` and `given (<init;; cond;; step;;>) -> { ... }`.
- Collection iteration uses `given (<Name within Names>) -> { ... }`.

## Collections and Literals

- `table<T>` uses `[...]`.
- `enumeration<T>` uses `<...>`.
- `lexis<K, V>` uses `{...}`.
- Literal forms are exact in v0.1 with no alternate syntax.
