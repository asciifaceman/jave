---
description: "Use when editing the GitHub Pages player-facing surface so governance records, correspondence, and language docs stay coherent with shipped behavior."
---

# GitHub Pages Lore Surface Instructions

## Purpose

The GitHub Pages site is a player-facing Orcal communications surface.

Treat it as in-world publication, not internal process documentation.

## Source Of Truth

When updating site content, align with:
- `governance/` records for version decisions and narrative timeline.
- `specs/` and `docs/` for actual language/runtime behavior.
- `examples/` for executable usage patterns.

Do not publish claims on the site that contradict shipped behavior or ratified governance records.

## Tone And Style

- Voice should read as official Orcal communications.
- Humor and corporate absurdity are allowed and encouraged.
- Keep records coherent, traceable, and tied to real features.
- Avoid meta commentary about "running an ARG" on player-facing pages.

## Character Dynamics

- Orcal staff should not read as uniformly cooperative.
- Reflect that leaders have competing goals, incentives, and territory concerns.
- Show strategic disagreement, political maneuvering, and attempts to redirect scope without turning characters into cartoon villains.
- Keep conflict grounded in real candidate tradeoffs, implementation constraints, governance timelines, and release risk.
- Preserve plausibility: they work together because program pressure requires it, not because they are personally aligned.

## Content Expectations

The site should surface:
- Governance updates by version.
- Record artifacts: minutes, mailing-list snippets, executive commentary, transition briefs.
- Tutorials, examples, and spec links.
- Forward-looking notes for planned doc systems (for example standard library docs pipeline).

## Update Workflow

- Changes to governance/docs/specs/examples should flow into Pages output on merges to `main`.
- Preserve generated feed/index behavior and avoid hand-editing generated outputs unless updating the generator itself.
- If new lore artifact types are introduced, update site navigation and generator output accordingly.

## Quality Checks

Before finalizing site changes:
- Verify links resolve within the site structure.
- Verify record entries map to actual version records.
- Verify language claims match current specification and examples.
