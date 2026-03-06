# Toolchain Guide

Jave v0.1 tool names are locked:
- `javec`: compiler
- `baggage`: package/build manager
- `javevm`: runtime

## Current Bootstrap Status

The CLIs are stubs for now and establish command shape.

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

Run tool stubs:

```bash
mage cmd:javec
mage cmd:baggage
mage cmd:javevm
```

No local Mage install fallback:

```bash
go run github.com/magefile/mage -l
```

## Build

```bash
go build ./cmd/javec ./cmd/baggage ./cmd/javevm
```

## Run Stubs

```bash
go run ./cmd/javec --version
go run ./cmd/baggage --version
go run ./cmd/javevm --version
```
