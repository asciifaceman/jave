# Jave Syntax Quick Reference

This is a practical reference for writing modern Jave syntax.

## v0.2 Naming Policy

- Exported sequence names (`outy seq`) use PascalCase (for example: `Foremost`, `AddNumbers`, `ReadConfig`).
- Builtin/runtime surface names are PascalCase (for example: `Pront`, `Prontulate`, `Girth`, `Slotify`, `ProntOops`, `FeudGirth`).
- Lowercase builtin spellings and legacy `Pronts.*` forms are not part of the canonical v0.2 surface.

## Declarations

```jave
allow exact Count 2b=2 5;;
allow vag Ratio 2b=2 0.6;;
allow truther Ready 2b=2 yee;;
allow strang Name 2b=2 "Jave";;
```

## Sequences

```jave
outy seq Add<exact A, exact B> --> <<exact>> {
    give A + B up;;
}

outy seq Foremost<> --> <<nada>> {
    give up;;
}
```

`Foreward` runs once per carryon at first load, including transitive imports.

## Conditionals

```jave
maybe (<X bigly 5>) -> {
    Pront("large");;
} furthermore (<X lessly 5>) -> {
    Pront("small");;
} otherwise -> {
    Pront("middle");;
}
```

## Loops

While-ish:

```jave
given (<X lesslysame 3>) again -> {
    X 2b=2 X + 1;;
}
```

For-ish:

```jave
given (<allow exact I 2b=2 0;; I lessly 3;; I 2b=2 I + 1;;>) -> {
    Pront(I);;
}
```

Collection iteration:

```jave
given (<Name within Names>) -> {
    Pront(Name);;
}
```

## Collections

```jave
allow table<exact> Scores 2b=2 [1, 2, 3];;
allow table<table<exact>> Grid 2b=2 [[1, 2], [3, 4]];;
allow enumeration<strang> Names 2b=2 <"Ada", "Linus">;;
allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;
```

These collection literal forms are exact in v0.1 with no alternates.

## Text Assembly

```jave
install Strangs from highschool/English;;

allow strang Message 2b=2 Strangs.Combobulate<"Hello, %strang", Name>;;
Pront(Message);;
```

Formatted output uses builtin `Prontulate<Template, Args...>`.

## Imports

```jave
install Strangs from highschool/English;;
```

Import graphs must be acyclic in v0.1. Cycles are rejected during load with an `import cycle detected` diagnostic that includes the resolved path chain.

## Operators

Comparison:
- `samewise`, `notsamewise`
- `bigly`, `lessly`
- `biglysame`, `lesslysame`

Logical:
- `plusalso`, `orelse`, `notno`
