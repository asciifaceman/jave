---
title: highschool/Embellishments Reference
permalink: /reference/carryons/highschool-embellishments/
category: carryon-reference
source_kind: source
---

# highschool/Embellishments

```jave
install Embellishments from highschool/Embellishments;;
```

## Table of Contents

- [`Banner`](#banner)
- [`Divider`](#divider)
- [`KeyVal`](#keyval)
- [`TwoCol`](#twocol)

## Banner

Framed banner line for console sections.

```jave
outy seq Banner<strang Title> --> <<strang>>
```

Produces a readable section banner for CLI and runtime logs.

### Parameters

- `strang Title`: Primary section label.

### Returns

- `strang`: Banner strang wrapped in governance-approved framing.

## Divider

Horizontal divider line.

```jave
outy seq Divider<> --> <<strang>>
```

Returns a static divider suitable for separating table blocks.

### Returns

- `strang`: Divider strang.

## KeyVal

Two-column key/value display line.

```jave
outy seq KeyVal<strang Key, strang Value> --> <<strang>>
```

Formats one key/value row for structured text output.

### Parameters

- `strang Key`: Row label.
- `strang Value`: Row value.

### Returns

- `strang`: Structured row strang.

## TwoCol

Simple two-column row for summary tables.

```jave
outy seq TwoCol<strang Left, strang Right> --> <<strang>>
```

Formats left/right values with visual separator for compact table-like output.

### Parameters

- `strang Left`: Left column value.
- `strang Right`: Right column value.

### Returns

- `strang`: Rendered row strang.

