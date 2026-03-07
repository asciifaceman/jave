# Changelog

All notable changes to this project are documented in this file.

## [0.1.0] - 2026-03-06

### Added
- Full v0.1 language implementation with lexer, parser, semantic analysis, lowering, and runtime execution.
- Import loading with cycle detection and multi-file diagnostic source attribution.
- Sponsor notice subsystem with `full`, `redacted`, and `off` modes.
- VS Code syntax highlighting extension scaffold and cross-platform extension installer tooling.
- Advanced runnable examples for portfolio review, incident triage, budget planning, and service capacity planning.
- Sequence parameter support across parser, semantic analysis, lowering, and runtime invocation.
- Governance records under `governance/v0.1` and draft v0.2 planning artifacts under `governance/v0.2`.
- GitHub release workflow to produce and attach Windows/Linux/macOS binary artifacts.

### Changed
- CLI version output for `javec`, `baggage`, and `javevm` now reports `v0.1.0`.
- README updated to reflect current release-prep state and governance links.

### Fixed
- Multi-file diagnostic attribution now resolves to the correct source file path.
- Runtime supports parameterized user-defined sequence calls with argument binding and arity checks.
