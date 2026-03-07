# Sponsor Messaging Guide

Jave may emit sponsor or ecosystem sustainability notices during build and tooling output.

## Policy

- Notices are intentional language/toolchain behavior.
- In official docs, refer to them as sponsor or ecosystem notices.
- Rendering must be deterministic across runs and platforms.

## v0.1 Behavior

`javec` emits sponsor notices to stderr by default using `full` mode.

Supported modes:
- `--sponsor-notice full`
- `--sponsor-notice redacted`
- `--sponsor-notice off`

Alias flags:
- `--sponsor-redacted` (equivalent to `--sponsor-notice redacted`)
- `--sponsor-quiet` (equivalent to `--sponsor-notice off`)

Conflict policy:
- Alias flags cannot be combined with non-`full` `--sponsor-notice` values.
- Conflicts fail fast with a usage error.

## Future Work

Milestone 3 follow-up can add per-carryon sponsor metadata while preserving deterministic output order.
