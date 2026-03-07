# Jave v0.1 Repo Bootstrap

## What Jave is

Jave is a deliberately committee-damaged programming language parodying open source governance, corporate-managed FOSS, and over-rationalized language design.

Jave must remain:

* readable after a short adjustment period
* internally consistent even when it is stupid
* semantically trustworthy in input/output behavior
* operationally embarrassing in tooling, naming, and runtime posture

Jave should not be random chaos. It should feel like a real language that was harmed by governance.

## Core design principles

1. **Correct results, weird path**
   Programs should produce the expected result. The implementation may be slower, more ceremonial, or more overengineered than necessary.

2. **Verbose where it hurts**
   Common tasks should often require a ratified wrapper, helper, or awkward syntax.

3. **Dual syntax where it feels cursed**
   Jave supports multiple forms for some core constructs when that inconsistency adds committee-scarring without making the language unusable.

4. **Naming is part of the joke**
   Built-in types are lowercase. Library/module names are PascalCase. Keywords are lowercase. Legacy aliases may survive for backwards compatibility.

5. **Tooling is part of the language**
   `javec`, `baggage`, `carryon`, JaveVM, partner messages, and awful CLI flags are all first-class parts of the joke.

## Locked v0.1 language decisions

### Visibility

* `outy` = exported for sequences/modules
* `inny` = internal/private for sequences/modules

### Entrypoints

* `Foreward` = one-shot pre-main carryon init sequence, runs once when carryon is first loaded
* `Foremost` = main executable entry sequence

### Primitive built-in types

* `exact` = integer
* `vag` = floating-point number
* `truther` = boolean
* `strang` = string
* `nada` = void/no return
* `naw` = nil/null

### Boolean literals

* `yee`
* `nee`

### Declaration and assignment

* Declaration starts with `allow`
* Assignment operator is `2b=2`
* Statements terminate with `;;`

Examples:

```jave
allow exact Count 2b=2 5;;
allow vag Ratio 2b=2 0.6;;
allow truther Ready 2b=2 yee;;
allow strang Name 2b=2 "Jave";;
```

### Sequences

* `sequence` / `seq`
* return syntax: `give X up;;` or bare `give up;;`

Example:

```jave
outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}
```

### Control flow

* `maybe` = if
* `furthermore` = else if
* `otherwise` = else

Example:

```jave
maybe (<X bigly 5>) -> {
    pront("large");;
} furthermore (<X lessly 5>) -> {
    pront("small");;
} otherwise -> {
    pront("middle");;
}
```

### Comparison operators

* `samewise`
* `notsamewise`
* `bigly`
* `lessly`
* `biglysame`
* `lesslysame`

### Logical operators

* `plusalso`
* `orelse`
* `notno`

### Built-in output

* `pront(...)` is built-in
* advanced formatted output may use `Strangs.Combobulate` or `Pronts.Prontulate`

## Text assembly

### Strang rules

* built-in type is `strang`
* standard library namespace is `Strangs`
* legacy alias `Srangs` is supported with warning for backwards compatibility

### Combobulate

String concatenation is not legal as an infix operator. Text assembly must proceed through ratified combobulation.

Canonical form:

```jave
Strangs.Combobulate<"template", Arg1, Arg2>
```

Examples:

```jave
allow strang Message 2b=2 Strangs.Combobulate<"Hello, %strang", Name>;;
pront(Strangs.Combobulate<"Count: %exact", Count>);;
```

### Combobulation directives

* `%exact`
* `%vag`
* `%tru`
* `%strang`
* `%glyph` (reserved for later if glyph lands)
* `%rat`
* `%maybe`
* legacy `%v` supported with warning

### Pronts

* `pront(...)` = built-in simple print
* `Pronts.Prontulate<...>` = wrapper around `Strangs.Combobulate(...)` + `pront(...)`
* richer output comes later through separate helpers like `PrettyPront` or `BetterPront`

Imports:

```jave
install Strangs from highschool/English;;
install Pronts from highschool/Communications;;
```

## Collections

### table

Ordered indexed collection with first-class multi-dimensional support.

Type forms:

```jave
table<exact>
table<table<exact>>
```

Literal form:

```jave
[1, 2, 3]
[[1, 2], [3, 4]]
```

Examples:

```jave
allow table<exact> Scores 2b=2 [1, 2, 3];;
allow table<table<exact>> Grid 2b=2 [[1, 2], [3, 4]];;
```

Indexing:

```jave
Scores[0]
Grid[1][0]
```

v0.1 requirement: multi-dimensional tables must work as nested tables and be documented as first-class supported usage.

### enumeration

Dynamic list-like collection.

Type form:

```jave
enumeration<strang>
```

Literal form:

```jave
<"Ada", "Linus", "Grace">
```

Example:

```jave
allow enumeration<strang> Names 2b=2 <"Ada", "Linus", "Grace">;;
```

### lexis

Map/dictionary-like collection.

Type form:

```jave
lexis<strang, exact>
```

Literal form:

```jave
{"Ada": 36, "Linus": 55}
```

Example:

```jave
allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;
```

## Collection sizing

Use built-in `girth(...)`.

Examples:

```jave
allow exact Count 2b=2 girth(Scores);;
allow exact NameCount 2b=2 girth(Names);;
allow exact Width 2b=2 girth(Grid[0]);;
```

## Looping

Dual syntax is intentional and core to Jave.

### while-ish loop

```jave
given (<X lesslysame 5>) again -> {
    X 2b=2 X + 1;;
}
```

### for-ish loop

```jave
given (<allow exact I 2b=2 0;; I lessly 5;; I 2b=2 I + 1;;>) -> {
    pront(I);;
}
```

### collection iteration

```jave
given (<Name within Names>) -> {
    pront(Name);;
}
```

Notes:

* The language intentionally has both `given (<cond>) again` and `given (<init;; cond;; step;;>)`.
* This should be treated as a first-order grammar requirement, not a later extension.

## Imports and standard library

Import syntax:

```jave
install Strangs from highschool/English;;
install Algebra from highschool/Algebra;;
install Pronts from highschool/Communications;;
```

Rules:

* `highschool/...` is reserved for standard library packages
* `Strangs` is canonical
* `Srangs` redirects with warning

## Toolchain

### Compiler

* `javec`

### Package/build manager

* `baggage`

### Package unit

* `carryon`

### Runtime

* JaveVM

### Compiled artifact

* `.jbin`

### Suggested commands

```bash
baggage new hello-jave
baggage build
baggage run
baggage check
baggage test
baggage add some/carryon
```

## Sponsor messaging / invasive ad ops

Compilation output may include sponsor or ecosystem sustainability notices.
These must not be called ads in official docs, though examples can acknowledge the joke.

Example suppression flags:

```bash
baggage build -JC:-UseEcosystemNotice -JC:+HideSponsorMessage -JCXxM:+DPartnerMessage=false
```

Behavior:

* sponsor lines are not fully removable in normal community mode
* flags can partially obscure them
* some flags yield `[partner message hidden by local preference policy]`
* other flags produce low-quality redaction / white-out / partial censorship

This is intentional.

## Diagnostics tone

Diagnostics must sound dry, technical, and absurdly formal.

Example language:

* `legacy module alias 'Srangs' remains supported for ecosystem continuity`
* `major transition exceeded declared compatibility comfort`
* `generic %v combobulation is tolerated but imprecise`

## Implementation target

Implement Jave in Go.

Why:

* easy parser/compiler prototyping
* fast iteration
* amusing self-reference given the project tone

## Repo bootstrap recommendation

```text
jave/
  README.md
  LICENSE
  go.mod
  cmd/
    javec/
      main.go
    baggage/
      main.go
    javevm/
      main.go
  internal/
    lexer/
    parser/
    ast/
    types/
    lowering/
    ir/
    vm/
    diagnostics/
    runtime/
    baggage/
    sponsor/
  stdlib/
    highschool/
      English/
      Algebra/
      Communications/
  examples/
    hello_world/
    conditions/
    loops/
    collections/
    multi_dimensional_tables/
    imports/
    foreward_foremost/
    combobulate/
  docs/
    syntax.md
    collections.md
    loops.md
    toolchain.md
    sponsor-messaging.md
    how-to-write-jave.md
    contributing.md
  specs/
    jave-v0.1.md
    diagnostics-style.md
    grammar-notes.md
```

## Example programs to include immediately

### 1. hello world

```jave
outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
    give up;;
}
```

### 2. foreward and foremost

```jave
outy seq Foreward<> --> <<nada>> {
    pront("warming carryon");;
    give up;;
}

outy seq Foremost<> --> <<nada>> {
    pront("running foremost");;
    give up;;
}
```

### 3. conditional logic

```jave
outy seq Foremost<> --> <<nada>> {
    allow vag Foo 2b=2 0.6;;

    maybe (<Foo bigly 0.5>) -> {
        pront("Over half");;
    } furthermore (<Foo lessly 0.5>) -> {
        pront("Under half");;
    } otherwise -> {
        pront("Exactly half");;
    }

    give up;;
}
```

### 4. combobulate and prontulate

```jave
install Strangs from highschool/English;;
install Pronts from highschool/Communications;;

outy seq Foremost<> --> <<nada>> {
    allow strang Name 2b=2 "Jave";;
    pront(Strangs.Combobulate<"Hello, %strang", Name>);;
    Pronts.Prontulate<"Still hello, %strang", Name>;;
    give up;;
}
```

### 5. tables, enumeration, lexis

```jave
install Strangs from highschool/English;;

outy seq Foremost<> --> <<nada>> {
    allow table<exact> Scores 2b=2 [1, 2, 3];;
    allow table<table<exact>> Grid 2b=2 [[1, 2], [3, 4]];;
    allow enumeration<strang> Names 2b=2 <"Ada", "Linus">;;
    allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;

    pront(Strangs.Combobulate<"Scores girth: %exact", girth(Scores)>);;
    pront(Strangs.Combobulate<"Grid[1][0]: %exact", Grid[1][0]>);;
    pront(Strangs.Combobulate<"First name: %strang", Names[0]>);;
    pront(Strangs.Combobulate<"Ada age: %exact", Ages["Ada"]>);;

    give up;;
}
```

### 6. dual given loops

```jave
install Strangs from highschool/English;;

outy seq Foremost<> --> <<nada>> {
    allow exact X 2b=2 0;;

    given (<X lesslysame 3>) again -> {
        pront(Strangs.Combobulate<"while-ish X: %exact", X>);;
        X 2b=2 X + 1;;
    }

    given (<allow exact I 2b=2 0;; I lessly 3;; I 2b=2 I + 1;;>) -> {
        pront(Strangs.Combobulate<"for-ish I: %exact", I>);;
    }

    give up;;
}
```

### 7. collection iteration

```jave
outy seq Foremost<> --> <<nada>> {
    allow enumeration<strang> Names 2b=2 <"Ada", "Linus", "Grace">;;

    given (<Name within Names>) -> {
        pront(Name);;
    }

    give up;;
}
```

## Agent implementation guidance

### Phase 1

Build the smallest coherent vertical slice:

* lexer
* parser
* AST
* diagnostics
* interpreter or direct VM execution path sufficient to run examples
* built-ins: `pront`, `girth`
* control flow: `maybe`, `furthermore`, `otherwise`, both `given` loop forms
* sequences with `Foremost`
* basic types: `exact`, `vag`, `truther`, `strang`, `nada`, `naw`
* collections: `table`, `enumeration`, `lexis`
* indexing and nested table indexing

### Phase 2

Add:

* imports
* `Foreward`
* `Strangs.Combobulate`
* `Pronts.Prontulate`
* `baggage` basic commands
* `.jbin` emission and JaveVM execution model

### Phase 3

Add:

* sponsor messaging subsystem
* suppression flags with partial censorship behavior
* diagnostics polish
* docs/how-tos
* more examples

## Documentation writing guidance

Every doc should:

* explain the actual behavior clearly
* preserve the Jave tone without obscuring real usage
* include at least one runnable example
* distinguish canonical syntax from legacy compatibility syntax

Required docs:

* getting started
* installation/build from source
* syntax quick reference
* collections guide
* loops guide
* formatting and combobulation guide
* carryons, baggage, and javec guide
* sponsor messaging guide
* contributing guide for implementers

Docs should be funny, but still genuinely useful.

## Immediate deliverables for repo start

1. `README.md` with project pitch and minimal hello world
2. `specs/jave-v0.1.md` with locked syntax and semantics
3. `docs/syntax.md` quick reference
4. `docs/how-to-write-jave.md` beginner guide
5. `examples/` programs matching the canonical samples above
6. `cmd/javec`, `cmd/baggage`, and `cmd/javevm` stubs in Go
7. parser/lexer milestone issues or task list for agents

## Non-goals for v0.1

Do not implement yet:

* classes
* traits
* async
* concurrency
* macros
* advanced generic systems
* optimizer beyond basic sanity
* rich standard library beyond what examples require

The goal is a working, cursed, teachable v0.1.

## Further instruction from the original design agent

 What I’d do next, in order:

First, tighten the few syntax choices that still need one final pass:

whether table, enumeration, and lexis literal forms are exactly [...], <...>, and {...}

whether Pronts.Prontulate must always be fully namespaced or can be imported into local scope

whether Foreward runs per carryon load only, or also for transitive imports

Second, turn the bootstrap into three first files for the repo:

README.md

specs/jave-v0.1.md

docs/syntax.md

Third, break implementation into agent-sized milestones:

lexer

parser

AST

evaluator/VM path

builtins and collections

imports and stdlib

javec / baggage / javevm

sponsor messaging

A couple of recommendations I’d make before agents start coding:

For multi-dimensional table, keep v0.1 simple and define it as nested tables, not a special matrix type. That keeps it first-class in practice without forcing complicated type machinery on day one.

For girth, make it work on strang, table, enumeration, and maybe lexis. That will make examples and docs much nicer.

For the first implementation, I’d actually have javec compile to an internal IR or serialized AST and let javevm execute that. You can still call it .jbin, keep the joke intact, and avoid overcommitting too early to a lower-level bytecode design.

The next strongest move is for me to turn the canvas content into:

a clean README.md

a proper specs/jave-v0.1.md

an agent task plan with milestones and acceptance criteria