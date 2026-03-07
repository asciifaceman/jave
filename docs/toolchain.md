# Toolchain Guide

Jave v0.1 tool names are locked:
- `javec`: compiler
- `baggage`: package/build manager
- `javevm`: runtime

## Current Bootstrap Status

Core compiler/runtime pipeline is functional.
- `javec` compiles `.jave` source and can emit `.jbin` artifacts.
- `javevm` executes `.jbin` artifacts.
- `baggage` currently supports `build`, `run`, `check`, and `test` workflows.

## Build System (Go Side)

Use Mage as the primary task runner.

Install:

```bash
go install github.com/magefile/mage@latest
```

List targets:

```bash
mage -l
```

Core targets:

```bash
mage build
mage test
mage check
mage bootstrap
```

Run tool workflows:

```bash
mage runJavec
mage runBaggage
mage runJavevm
```

No local Mage install fallback:

```bash
go run github.com/magefile/mage -l
```

## Build

```bash
go build ./cmd/javec ./cmd/baggage ./cmd/javevm
```

## Run CLIs

```bash
go run ./cmd/javec --version
go run ./cmd/baggage --version
go run ./cmd/javevm --version
```

Trace import resolution and carryon load order:

```bash
go run ./cmd/javec --trace-imports --run examples/imports/main.jave
```

Force a specific project root for highschool carryon resolution:

```bash
go run ./cmd/javec --project-root . --trace-imports --run examples/imports/main.jave
```

Compile a source file into a `.jbin` artifact:

```bash
go run ./cmd/baggage build examples/hello_world/main.jave
```

If no input is provided, `baggage build` and `baggage run` resolve input in this order:
- `JAVE_FILE` environment variable
- `entry` from local `baggage.jave` manifest (if present)
- fallback `examples/hello_world/main.jave`

Run a source file (compile then execute):

```bash
go run ./cmd/baggage run examples/conditions/main.jave
```

Pass through import tracing from baggage to javec:

```bash
go run ./cmd/baggage run --trace-imports examples/imports/main.jave
```

Run with explicit project root (useful outside repo root):

```bash
go run ./cmd/baggage run --project-root . --trace-imports examples/imports/main.jave
```

Sponsor notice behavior is deterministic and controlled in `javec` with:
- `--sponsor-notice full|redacted|off`
- `--sponsor-redacted` (alias for redacted)
- `--sponsor-quiet` (alias for off)

Examples:

```bash
go run ./cmd/javec --sponsor-notice redacted examples/hello_world/main.jave
go run ./cmd/javec --sponsor-quiet examples/hello_world/main.jave
```

`baggage build` and `baggage run` pass these sponsor flags through to `javec`.

Scaffold a new project:

```bash
go run ./cmd/baggage new hello-jave
```

Run an existing artifact directly:

```bash
go run ./cmd/baggage run examples/hello_world/main.jbin
```

Add a carryon dependency to a project manifest:

```bash
go run ./cmd/baggage add some/carryon
```

Use a custom manifest path:

```bash
go run ./cmd/baggage add --manifest path/to/baggage.jave some/carryon
```
