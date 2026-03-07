---
title: Language Feature Reference
permalink: /reference/language/
category: language-reference
source_kind: manifest
---

# Language Feature Reference

## Table of Contents

- [`2b=2`](#2b-2)
- [`given`](#given)
- [`maybe/furthermore/otherwise`](#maybe-furthermore-otherwise)

## 2b=2

Assignment operator

```jave
Name 2b=2 Expr;;
```

Assigns a new value to an existing identifier, or initializes one during allow declarations.

### Examples

```jave
allow exact Count 2b=2 1;;
Count 2b=2 Count + 1;;

```

### Status

stable

## given

Loop family keyword

```jave
given (<...>) -> { ... }
```

Supports while-ish, for-ish, and within iteration forms.

### Examples

```jave
given (<allow exact I 2b=2 0;; I lessly 3;; I 2b=2 I + 1;;>) -> {
    pront(I);;
}

```

### Status

stable

## maybe/furthermore/otherwise

Conditional control flow chain

```jave
maybe (<Condition>) -> { ... } furthermore (...) -> { ... } otherwise -> { ... }
```

Evaluates branch conditions top-to-bottom and executes the first matching branch.

### Examples

```jave
maybe (<X bigly 5>) -> {
    pront("large");;
} otherwise -> {
    pront("small-or-equal");;
}

```

### Status

stable

