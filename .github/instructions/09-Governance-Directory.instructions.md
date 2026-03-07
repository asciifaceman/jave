---
description: "Use when creating or editing governance records under governance/. Keeps governance materials externally readable and separates implementation process from governance records."
---

# Governance Directory Handling

## Boundary Rules

- `governance/` is for external-facing governance records and release evidence.
- Do not include internal contributor workflow instructions in `governance/`.
- Do not include explicit game framing language (for example: ARG, simulation, player-facing) in `governance/` documents.

## What Belongs In Governance

- Version governance packs under `governance/vX.Y/`.
- Feature ledgers, user stories, steering minutes, policy memos, and release rationale.
- Narrative style is acceptable as long as it remains institutional and internally coherent.

## What Does Not Belong In Governance

- Agent operating instructions.
- Internal implementation workflow templates.
- Tooling-specific contributor instructions unrelated to governance records.

## Where Internal Process Goes

- Place contributor process and agent-operational guidance in `.github/instructions/`.
- Prefer version-agnostic instruction files for reusable process.

## Consistency Expectations

- Keep dates, participants, and decision trails consistent across records.
- Ensure governance records do not contradict specs, release notes, or shipped behavior.
- When introducing a new version folder, include at least: `README.md`, feature slate, user stories, and one lore/policy artifact.
