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

Compile a source file into a `.jbin` artifact:

```bash
go run ./cmd/baggage build examples/hello_world/main.jave
```

Run a source file (compile then execute):

```bash
go run ./cmd/baggage run examples/conditions/main.jave
```

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
