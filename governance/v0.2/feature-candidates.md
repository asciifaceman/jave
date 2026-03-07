# v0.2 Feature Candidates

Status: Draft

Decision note (2026-03-06): Cheddar approved the following as major v0.2 candidates:
- I/O Foundation
- File System Operations
- Richer Standard Library Surface
- Runtime Error and Exit Semantics

Decision note (2026-03-07, ad-hoc steering consensus from Slavala, Nafta, Catherine, and Cheddar):
- Environment variable reads
- Structured logging basic pass
- Runtime timing/profiling helpers
- Networking stack basic pass

Decision note (2026-03-07, steering addendum):
- Code comments and docstring standard for generated Jave documentation
- Documentation generator output to Markdown or Jekyll-compatible Markdown (flagged functionality)
- Tooling posture: documentation generation is orchestrated by `baggage`

Decision note (2026-03-07, import/distribution direction):
- Non-stdlib imports should support community carryons via GitHub/VCS-sourced library style.
- Avoid custom hosting infrastructure requirements for community carryons.
- Preserve distinction between bundled stdlib distribution and external dependency acquisition.

## Priority: Required Candidates

1. Runtime I/O Foundation (Approved)
- Standard input handling (`stdin` reads)
- Explicit output channels (`stdout`, `stderr`)
- Process argument access
- Basic flag parsing utilities for CLI-oriented programs

2. File System Operations (Approved)
- Read text and binary files
- Write/append text and binary files
- File existence and metadata checks
- Standardized file-path diagnostics

3. Richer Standard Library Surface (Approved)
- String utilities beyond combobulation
- Numeric helpers and conversion utilities
- Collection helper functions (map/filter-like primitives under governance-approved naming)
- Time/date primitives suitable for logs and scheduling use cases

4. Runtime Error and Exit Semantics (Approved)
- Structured runtime error emission to stderr
- Explicit non-zero exit signaling from program code paths
- Improved call-stack context for runtime failures

5. Path and Working Directory Utilities (Proposed)
- Canonical path join/normalize helpers across Windows/Linux/macOS
- Current working directory query helper
- Predictable relative-vs-absolute path behavior in diagnostics

6. Deterministic Program Exit Contract (Proposed)
- Runtime-level mapping of failure classes to stable exit codes
- Explicit contract for compiler/runtime/toolchain exit code meanings
- Docs and tests for automation-facing behavior

7. Environment Variable Reads (Approved)
- Runtime read access for environment variables
- Missing-variable handling with explicit diagnostics guidance
- Cross-platform expectations for variable name behavior

8. Structured Logging Basic Pass (Approved)
- Structured key/value log emission helpers
- Standardized minimal fields for operational traces
- Logging output behavior aligned with stdout/stderr channel policy

9. Runtime Timing/Profiling Helpers (Approved)
- Baseline elapsed-time measurement helpers
- Program-step timing instrumentation primitives
- Low-overhead usage suitable for production troubleshooting

10. Networking Stack Basic Pass (Approved)
- Foundational network request/response primitives
- Minimal, deterministic error surface for network failures
- Initial scope constrained to practical baseline operations

11. Comment and Docstring Documentation Standard + Generator (Approved)
- Define comment/docstring conventions for sequences, modules, and exported surfaces
- Add documentation generation pipeline producing Markdown and Jekyll-compatible Markdown
- Support output-mode flags for generator formatting behavior
- Documentation generation entrypoint is delivered through `baggage` (with compiler/runtime support as needed)

12. Community Carryon Import and Distribution Model (Proposed, High Priority)
- Support non-stdlib imports using VCS/GitHub-style source references.
- Resolve and cache carryons locally without requiring Orcal-hosted package infrastructure.
- Define deterministic import resolution precedence: stdlib (toolchain-bundled) vs community carryons (external source).
- Provide lock/manifest expectations through `baggage` for reproducible builds.
- Keep custom package manager service out of scope for v0.2 baseline.

## Priority: Optional Candidates

1. Configuration and Environment Access
- Typed config loading helpers

2. Observability Helpers


3. Data Exchange Utilities
- JSON-like encode/decode helpers (scope subject to naming review)

## Deferred Candidates

- Concurrency model primitives
- Custom package manager service and centralized remote registry infrastructure
- Generics and advanced type features

## Decision Criteria

A candidate is promoted to ratified v0.2 scope when:
- It has implementation feasibility review approval.
- It has governance rationale and user-story traceability.
- It has acceptance tests and diagnostics requirements defined.

## Acceptance Artifacts (Per Candidate)

Before ratification, each candidate should include:
- API and naming proposal with at least one rejected alternative.
- Diagnostics and error behavior expectations.
- Cross-platform notes (Windows/Linux/macOS) when file/process behavior is involved.
- At least one acceptance scenario in `examples/` or equivalent test coverage.
