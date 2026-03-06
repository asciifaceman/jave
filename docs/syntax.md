# Jave Syntax Quick Reference (v0.1)

This is a practical reference for writing valid Jave v0.1.

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
    pront("large");;
} furthermore (<X lessly 5>) -> {
    pront("small");;
} otherwise -> {
    pront("middle");;
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
    pront(I);;
}
```

Collection iteration:

```jave
given (<Name within Names>) -> {
    pront(Name);;
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
pront(Message);;
```

`Prontulate` must be called as `Pronts.Prontulate<...>` in v0.1.

## Imports

```jave
install Strangs from highschool/English;;
install Pronts from highschool/Communications;;
```

## Operators

Comparison:
- `samewise`, `notsamewise`
- `bigly`, `lessly`
- `biglysame`, `lesslysame`

Logical:
- `plusalso`, `orelse`, `notno`
