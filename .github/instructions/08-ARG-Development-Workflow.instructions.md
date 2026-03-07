---
description: "Internal contributor workflow for planning post-v0.1 releases with feature-first implementation and corroborating ARG artifacts. Not player-facing."
---

# ARG Development Workflow (Internal)

Use this file for contributor process, not in-world lore.

## Version Planning Contract

1. Define feature slate first.
- Draft version scope (syntax, semantics, runtime, tooling, diagnostics, docs).
- Mark each feature as `required`, `optional`, or `deferred`.

2. Produce corroborating ARG evidence.
- Create supporting artifacts aligned to feature clusters.
- Required artifact classes:
- Steering minutes
- Product/policy memo
- User stories
- Feature ledger

3. Implement and validate.
- Add/adjust tests for each non-trivial behavior change.
- Validate diagnostics quality and cross-platform behavior.
- Keep docs/spec and lore consistent with shipped behavior.

4. Freeze and provenance.
- Run freeze checklist.
- Archive artifacts under `governance/vX.Y/`.

## Artifact Paths

- `governance/vX.Y/lore/YYYY-MM-DD-steering-minutes.md`
- `governance/vX.Y/lore/YYYY-MM-DD-policy-memo.md`
- `governance/vX.Y/user-stories.md`
- `governance/vX.Y/feature-ledger.md`

## Quality Bar

A version is complete when:
- `go test ./...` is green.
- Feature ledger is complete.
- Lore evidence exists for major decisions.
- Docs/spec and implemented behavior do not contradict.
