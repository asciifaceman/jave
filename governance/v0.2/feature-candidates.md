# v0.2 Feature Candidates

Status: Draft

## Priority: Required Candidates

1. Runtime I/O Foundation
- Standard input handling (`stdin` reads)
- Explicit output channels (`stdout`, `stderr`)
- Process argument access
- Basic flag parsing utilities for CLI-oriented programs

2. File System Operations
- Read text and binary files
- Write/append text and binary files
- File existence and metadata checks
- Standardized file-path diagnostics

3. Richer Standard Library Surface
- String utilities beyond combobulation
- Numeric helpers and conversion utilities
- Collection helper functions (map/filter-like primitives under governance-approved naming)
- Time/date primitives suitable for logs and scheduling use cases

4. Runtime Error and Exit Semantics
- Structured runtime error emission to stderr
- Explicit non-zero exit signaling from program code paths
- Improved call-stack context for runtime failures

## Priority: Optional Candidates

1. Configuration and Environment Access
- Environment variable reads
- Typed config loading helpers

2. Observability Helpers
- Basic structured log formatting
- Runtime timing helpers

3. Data Exchange Utilities
- JSON-like encode/decode helpers (scope subject to naming review)

## Deferred Candidates

- Concurrency model primitives
- Networking stack
- Package manager and remote dependency resolution
- Generics and advanced type features

## Decision Criteria

A candidate is promoted to ratified v0.2 scope when:
- It has implementation feasibility review approval.
- It has governance rationale and user-story traceability.
- It has acceptance tests and diagnostics requirements defined.
