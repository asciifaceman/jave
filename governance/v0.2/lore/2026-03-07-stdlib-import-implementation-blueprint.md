# Stdlib Import Implementation Blueprint

Date: 2026-03-07
Recorder: Governance Office
Purpose: Define the minimum implementation path for real stdlib execution and eventual Jave-authored `Strangs.Combobulate`.

## Goal

Move from intrinsic-only namespace behavior to true module-backed stdlib execution, starting with `Strangs`.

## Current Constraint Snapshot

- Import loading is real (`install ... from ...` loads source and merges programs).
- Runtime dispatch for `Strangs.Combobulate` and `Pronts.Prontulate` is still hardcoded.
- Program lowering currently stores sequences in a flat map keyed by unqualified sequence name.
- No module-qualified runtime dispatch path exists yet.

## Phase 1: Real Module Dispatch (Minimum Viable)

1. Module-qualified sequence identity
- Introduce qualified sequence identity in IR (`ModuleAlias.SequenceName`) for imported module exports.
- Preserve unqualified names for local/root sequences.

2. Import alias binding map
- During load, keep deterministic map: `alias -> resolved module source` for root program imports.
- Carry this into lowered IR metadata so runtime can resolve member calls.

3. Member call execution path
- In runtime `evalCall`, when encountering `MemberExpr` (`Alias.Member`), try module dispatch first:
  - resolve `Alias.Member` to qualified sequence
  - execute that sequence with current argument evaluation semantics
- Keep intrinsic fallback temporarily for compatibility gates.

4. Sema checks
- Validate that imported alias/member references resolve to exported sequences.
- Emit explicit diagnostics for missing alias/member combinations.

## Phase 2: Strangs Migration (Module-First)

1. Add Jave-authored `highschool/English/main.jave` surface for `Combobulate` and companion helpers.
2. Route `Strangs.Combobulate` through module dispatch first.
3. Keep intrinsic fallback behind compatibility guard while parity tests are built.
4. Remove hardcoded fallback only after parity and diagnostics stability are proven.

## Phase 3: Language/Runtime Capability Gaps For Pure In-Language Combobulate

`Combobulate` can be module-hosted early, but cannot be fully intrinsic-free until these capabilities exist:

1. String traversal and slicing primitives
- Deterministic string length, substring/slice, and index/find behavior.

2. Controlled text replacement primitives
- First-match replacement or equivalent composable operations.

3. Variadic or list-based formatting input model
- Current fixed-arity sequence params make general directive replacement awkward.

4. Stable value-to-text conversion contract
- Explicit conversion semantics across exact/vag/truther/strang and other supported forms.

Without these, a Jave-authored `Combobulate` becomes either incomplete or dependent on hidden intrinsic helpers.

## Compatibility Policy During Migration

- `Prontulate` remains usable both as builtin identifier and namespaced call during transition.
- Existing examples and docs must keep functioning while dispatch changes land.
- Records/docs must explicitly call out any intrinsic-backed fallback still in use.

## Acceptance Criteria

Phase 1 done when:
- Member calls execute imported module sequences without intrinsic hardcoding for at least one pilot module.
- Missing module/member diagnostics are deterministic and tested.

Phase 2 done when:
- `Strangs.Combobulate` resolves through module dispatch in normal operation.
- Runtime fallback path is either removed or clearly feature-flagged with documented rationale.

Phase 3 done when:
- `Combobulate` behavior can be authored and maintained primarily in Jave source with no hidden behavior gaps.
