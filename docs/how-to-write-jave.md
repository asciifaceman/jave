# How To Write Jave

Jave looks odd on purpose, but it follows rules and can be learned quickly.

## 1. Start with `Foremost`

```jave
outy seq Foremost<> --> <<nada>> {
    pront("hello, jave");;
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
    pront("big enough");;
} otherwise -> {
    pront("too small");;
}
```

## 4. Use `given` for loops

```jave
given (<allow exact I 2b=2 0;; I lessly 3;; I 2b=2 I + 1;;>) -> {
    pront(I);;
}
```

## 5. Build text with combobulation

String concatenation is not available in v0.1.

```jave
install Strangs from highschool/English;;
pront(Strangs.Combobulate<"Hello, %strang", Name>);;
```

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
