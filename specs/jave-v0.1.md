# Jave v0.1 Specification (Locked)

This document captures locked language decisions for v0.1.

## Identity

Jave is intentionally awkward but semantically reliable.

Design constraints:
- Results must be correct and deterministic.
- Syntax and tooling may be ceremonial, verbose, and inconsistent by design.
- Inconsistency must still be rule-based, not random.

## Visibility

- `outy`: exported symbol
- `inny`: internal symbol

## Entrypoints

- `Foreward`: one-shot carryon initialization sequence; runs once per carryon at first load, including transitive imports
- `Foremost`: executable entry sequence

## Primitive Types

- `exact`: integer
- `vag`: floating-point
- `truther`: boolean
- `strang`: string
- `nada`: no value/void
- `naw`: nil/null

Boolean literals:
- `yee`
- `nee`

## Declarations and Statements

- Variable declaration starts with `allow`.
- Assignment operator is `2b=2`.
- Statements terminate with `;;`.

Example:

```jave
allow exact Count 2b=2 5;;
allow truther Ready 2b=2 yee;;
```

## Sequences (Functions)

- Keywords: `sequence` and alias `seq`
- Return syntax:
  - `give X up;;`
  - `give up;;`

Example:

```jave
outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}
```

## Control Flow

- `maybe`: if
- `furthermore`: else-if
- `otherwise`: else

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

## Operators

Comparison operators:
- `samewise`
- `notsamewise`
- `bigly`
- `lessly`
- `biglysame`
- `lesslysame`

Logical operators:
- `plusalso`
- `orelse`
- `notno`

## Output

- `pront(...)` is built-in output.
- `prontulate<...>` is built-in formatted output.
- `Strangs.Combobulate<...>` is canonical text assembly.
- `prontulate<...>` formats by way of `Strangs.Combobulate` behavior.

## Text Assembly

Infix plus-style text joining is not legal in v0.1.

Canonical form:

```jave
Strangs.Combobulate<"template", Arg1, Arg2>
```

Directives:
- `%exact`
- `%vag`
- `%tru`
- `%strang`
- `%glyph` (reserved)
- `%rat`
- `%maybe`
- `%v` (legacy, warning required)

## Imports

Canonical import form:

```jave
install Strangs from highschool/English;;
```

Rules:
- `highschool/...` path is reserved for standard library carryons.
- `Strangs` is canonical.
- `Srangs` legacy alias is supported with warning.
- Import cycles are rejected with an `import cycle detected` diagnostic that includes the resolved cycle path.

## Collections

### `table`

Ordered indexed collection with first-class nested usage.

Type forms:

```jave
table<exact>
table<table<exact>>
```

Literal forms:

```jave
[1, 2, 3]
[[1, 2], [3, 4]]
```

No alternate literal syntax is supported in v0.1.

Indexing:

```jave
Scores[0]
Grid[1][0]
```

### `enumeration`

Dynamic list-like collection.

Type form:

```jave
enumeration<strang>
```

Literal form:

```jave
<"Ada", "Linus", "Grace">
```

No alternate literal syntax is supported in v0.1.

### `lexis`

Map-like collection.

Type form:

```jave
lexis<strang, exact>
```

Literal form:

```jave
{"Ada": 36, "Linus": 55}
```

No alternate literal syntax is supported in v0.1.

## Sizing

Use built-in `girth(...)`.

Examples:

```jave
allow exact Count 2b=2 girth(Scores);;
allow exact Width 2b=2 girth(Grid[0]);;
```

## Looping (Dual Syntax is Required)

While-ish form:

```jave
given (<X lesslysame 5>) again -> {
    X 2b=2 X + 1;;
}
```

For-ish form:

```jave
given (<allow exact I 2b=2 0;; I lessly 5;; I 2b=2 I + 1;;>) -> {
    pront(I);;
}
```

Collection iteration:

```jave
given (<Name within Names>) -> {
    pront(Name);;
}
```

## Toolchain Names (Locked)

- Compiler: `javec`
- Package/build manager: `baggage`
- Package unit: `carryon`
- Runtime: `JaveVM`
- Compiled artifact extension: `.jbin`

## Sponsor Notice Policy

- `javec` emits deterministic sponsor/ecosystem notices by default.
- Supported controls: `--sponsor-notice full|redacted|off`, `--sponsor-redacted`, `--sponsor-quiet`.
- Alias flags are convenience forms and may not conflict with explicit non-default mode values.

## Diagnostics Tone

Diagnostics should sound dry, technical, and absurdly formal.

Examples:
- `legacy module alias 'Srangs' remains supported for ecosystem continuity`
- `generic %v combobulation is tolerated but imprecise`

## v0.1 Non-goals

Do not implement in v0.1:
- classes
- traits
- async
- concurrency
- macros
- advanced generics
- heavy optimization
- rich standard library beyond immediate examples
