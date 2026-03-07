---
title: highschool/English Reference
permalink: /reference/carryons/highschool-english/
category: carryon-reference
source_kind: source
---

# highschool/English

```jave
install Strangs from highschool/English;;
```

## Table of Contents

- [`Combobulate`](#combobulate)

## Combobulate

Formats template text with sequential directive replacement.

```jave
outy seq Combobulate<strang Template, ...strang Args> --> <<strang>>
```

Combobulate is the canonical standard-library text assembly sequence.
	It replaces template directives in order with the provided argument values.

### Parameters

- `strang Template`: Base strang containing directive slots such as %exact and %strang.
- `...strang Args`: Variadic replacement values consumed left-to-right.

### Returns

- `strang`: Fully combobulated strang output.

### Examples

```jave
allow strang Message 2b=2 Strangs.Combobulate<"Count=%exact", 3>;;
```

