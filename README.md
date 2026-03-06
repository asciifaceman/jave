# Jave

Jave is a deliberately committee-damaged programming language.

It is a joke language with real behavior: the syntax is cursed on purpose, but programs should still produce correct and predictable results.

## Project Status

This repository is in bootstrap mode for Jave v0.1.

Current focus:
- Lock language syntax and semantics.
- Build a minimal Go 1.26 implementation.
- Keep behavior consistent on Windows, Linux, and macOS.

## Quick Start (Current)

Today this repo includes docs, examples, and Go CLI stubs for:
- `javec` (compiler)
- `baggage` (package/build tool)
- `javevm` (runtime)

Build the stubs:

```bash
go build ./cmd/javec ./cmd/baggage ./cmd/javevm
```

Run one:

```bash
go run ./cmd/javec --help
```

## Go Workflow With Mage

We use Mage to keep local and CI command flows simple.

Install Mage (one-time):

```bash
go install github.com/magefile/mage@latest
```

Common tasks:

```bash
mage -l
mage build
mage test
mage check
mage cmd:javec
mage cmd:baggage
mage cmd:javevm
```

If Mage is not installed yet, run it via Go:

```bash
go run github.com/magefile/mage -l
```

## Hello, Jave

```jave
outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
    give up;;
}
```

## Read These First

- `specs/jave-v0.1.md`: Locked v0.1 syntax and semantics.
- `docs/syntax.md`: Quick reference.
- `docs/how-to-write-jave.md`: Friendly beginner guide.
- `docs/agent-milestones.md`: Engineering plan and acceptance criteria.

## Repo Layout

```text
cmd/            CLI tools: javec, baggage, javevm
docs/           Human docs and guides
examples/       Runnable language samples
specs/          Language specs and behavior definitions
.github/        Agent and instruction configuration
```

## Contributing

Early contributions should prioritize parser/runtime correctness, diagnostics quality, and cross-platform behavior.

Tone can be funny. Behavior must stay trustworthy.
