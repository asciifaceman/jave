---
name: Bug report
about: Report incorrect behavior, crashes, diagnostics issues, or regressions
title: "bug: "
labels: ["bug", "triage"]
assignees: []
---

## Summary

Describe the operational weirdness in one or two sentences.

## Environment

- OS: (Windows/Linux/macOS + version)
- Jave version: (`javec --version`, `baggage --version`, `javevm --version`)
- Install method: (release binary / go install / local build)

## Reproduction

Provide the minimum ritual required to reproduce.

1. 
2. 
3. 

## Minimal Jave source

```jave
// Paste the smallest cursed-but-valid failing example
```

## Command and output

Command run:

```bash
# Example
javec --run path/to/main.jave
```

Observed output/diagnostic:

```text
# Paste output/diagnostic exactly as seen
```

## Expected behavior

Describe what should have happened in a rational universe.

## Regression check

- [ ] I confirmed this still reproduces on the latest `main`
- [ ] I searched existing issues for duplicates
- [ ] I considered whether I am the problem

## Extra context

Add screenshots, links, or related issues if useful to the incident board.
