# Community Carryon Import Model Record

Date: 2026-03-07
Recorder: Governance Office
Context: v0.2 design direction for non-stdlib imports and distribution

## Decision Direction

For non-stdlib imports, Jave should support community carryons using a GitHub/VCS-style source model.

Goal: allow dependency sharing without requiring Orcal-hosted package infrastructure.

## Distribution Posture

1. Stdlib modules are bundled with toolchain distribution.
2. Community carryons are external source dependencies.
3. Resolution and build reproducibility are orchestrated by `baggage`.

## Import Resolution Boundary

- `highschool/...` remains reserved for bundled standard library modules.
- Non-stdlib imports should resolve via explicit VCS/GitHub-style references under a ratified import syntax.
- Resolution policy must be deterministic and diagnosable.

## Build Reproducibility Expectations

- Carryon source resolution must support lock/manifest semantics.
- CI and local builds should produce consistent dependency selections.
- Failure modes (missing repo/tag/revision) must have explicit diagnostics.

## Out Of Scope For v0.2 Baseline

- Orcal-hosted centralized package registry.
- Custom ecosystem service requirements for library publication.

## Open Design Questions

1. Exact syntax for VCS/GitHub import references.
2. Lock file format and lifecycle (`baggage` ownership).
3. Security and trust posture for external source resolution.
4. Cache layout and invalidation behavior across platforms.
