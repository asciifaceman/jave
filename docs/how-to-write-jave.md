# How To Write Jave

Jave looks odd on purpose, but it follows rules and can be learned quickly.

## 1. Start with `Foremost`

```jave
outy seq Foremost<> --> <<nada>> {
    Pront("hello, jave");;
    give up;;
}
```

## 2. Declare values with `allow`

```jave
allow exact Count 2b=2 5;;
allow strang Name 2b=2 "Ada";;
```

## 3. Use `maybe` for branching

```jave
maybe (<Count bigly 3>) -> {
    Pront("big enough");;
} otherwise -> {
    Pront("too small");;
}
```

## 4. Use `given` for loops

```jave
given (<allow exact I 2b=2 0;; I lessly 3;; I 2b=2 I + 1;;>) -> {
    Pront(I);;
}
```

## 5. Build text with combobulation

Infix plus-style text joining (`"a" + "b"`) is not available in v0.1.
Use `Strangs.Combobulate<...>` for direct text assembly, or builtin `Prontulate<...>` for formatted print.

```jave
install Strangs from highschool/English;;
Pront(Strangs.Combobulate<"Hello, %strang", Name>);;
```

`Strangs.Combobulate` is a core v0.1 feature, not a post-v0.1 addition.

## 6. Work with collections

```jave
allow table<exact> Scores 2b=2 [1, 2, 3];;
allow enumeration<strang> Names 2b=2 <"Ada", "Linus">;;
allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;
```

## 7. Learn the tone

Jave diagnostics are formal and slightly absurd. The language is comedic, but behavior should be clear and reliable.

## Next Reading

- `docs/syntax.md`
- `specs/jave-v0.1.md`
- `examples/`
