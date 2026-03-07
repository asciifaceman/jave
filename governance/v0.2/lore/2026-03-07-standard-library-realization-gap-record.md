# Standard Library Realization Gap Record

Date: 2026-03-07
Recorder: Governance Office
Context: Post-discussion review of `Strangs`/`Pronts` behavior versus import/governance expectations

## Current Reality

- Import loading is real for `.jave` files: imported source is parsed, merged, and analyzed.
- `Strangs.Combobulate` and `Pronts.Prontulate` behavior is currently runtime-hardcoded rather than executed from imported Jave module code.
- This creates a capability gap between standard library branding and implementation posture.

## Sufficiency Assessment: Can Current Language Move `Strangs.Combobulate` To Pure In-Language Implementation?

Short answer: not fully, yet.

What is already sufficient:
- Real import loading and merged analysis of `.jave` modules.
- Sequence declarations and sequence calls with parameters.
- Namespace-like call syntax at source level (`Module.Member<...>`).

What is currently missing for true module-backed execution:
- Member-call dispatch to imported module exports (runtime currently hardcodes `Strangs`/`Pronts` behavior).
- Imported-module symbol table binding that survives into runtime execution.
- Clear intrinsic bridge model for operations that cannot yet be expressed ergonomically in core language.

Implication:
- `Strangs.Combobulate` cannot yet be considered a fully Jave-authored standard-library implementation.

## Immediate Steering-Aligned Adjustment

- `Prontulate` is allowed as a first-class builtin identifier for operational consistency.
- `Pronts.Prontulate` remains supported for compatibility.

## What Is Needed To Make `Strangs` A True Jave-Authored Standard Library Module

1. Namespaced sequence dispatch
- Resolve member calls like `Strangs.Combobulate<...>` to imported module exports, not runtime hardcoded branches.

2. Module export contract
- Define explicit exported sequence model for imported modules and stable symbol lookup rules.

3. Import alias binding semantics
- Bind `install Name from path;;` to concrete module symbol tables so member resolution is deterministic.

4. Standard library bootstrap policy
- Define which modules are pure Jave code, which require runtime intrinsics, and how mixed-mode modules are declared.

5. Diagnostics for module/member failures
- Missing module/member diagnostics should be explicit and consistent with existing diagnostic tone.

6. Compatibility path
- Maintain support for existing call sites while migrating hardcoded implementations behind module-facing adapters.

7. External carryon and stdlib resolution boundary
- Import resolution must distinguish bundled stdlib modules from community carryon modules so module-backed dispatch remains deterministic.

## Proposed v0.2 Sequencing

1. Ratify module export/member-resolution semantics.
2. Implement namespaced member dispatch in sema + runtime.
3. Move `Pronts.Prontulate` to module-first implementation with builtin fallback during transition.
4. Move `Strangs.Combobulate` to module-first implementation where feasible.
5. Remove hardcoded runtime paths once module parity and tests are complete.

## `Combobulate` Migration Starter Plan

1. Keep `Combobulate` behavior stable while introducing module-export dispatch plumbing.
2. Add module-first dispatch path with intrinsic fallback behind explicit compatibility guard.
3. Port `Strangs` module logic to Jave-authored source where language/runtime support permits.
4. Remove fallback only after parity tests prove module implementation equivalence.

## Governance Position

The project should avoid pretending imported standard library modules are implemented in Jave when behavior is still intrinsic-only.

Until migration is complete, records and docs should explicitly state intrinsic-backed behavior.
