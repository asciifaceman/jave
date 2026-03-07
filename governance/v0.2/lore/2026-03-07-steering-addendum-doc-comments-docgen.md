# Steering Addendum: Comment/Docstring Standard and Documentation Generation

Date: 2026-03-07
Recorder: Governance Office
Participants referenced: Slavala, Nafta, Catherine, Cheddar

## Addendum Purpose

Capture the approved high-priority candidate for documentation standards and generation, including tool ownership posture.

## Approved Candidate

- Code comments and docstring standard for Jave documentation generation.
- Generator output targets:
  - Markdown
  - Jekyll-compatible Markdown
- Flagged functionality for output-mode control.

## Tooling Placement Decision

Steering favored `baggage` as the operational entrypoint for documentation generation because it already serves as the workflow orchestrator for build/run behavior.

Implementation model:
- `baggage` provides user-facing command routing and flags.
- Compiler/runtime layers provide parse/semantic/doc extraction support as required.
- Output must remain deterministic for CI publication workflows.

## Open Naming Work

Required follow-up record:
- Define naming convention for doc comments/docstrings and generator commands.
- Include at least one rejected naming alternative.

## Governance Constraint

Documentation generation scope remains candidate-approved but must still pass feasibility review and acceptance test definition before ratified implementation starts.
