---
title: Builtin Reference
permalink: /reference/builtins/
category: builtin-reference
source_kind: manifest
---

# Builtin Reference

## Table of Contents

- [`girth`](#girth)
- [`pront`](#pront)
- [`prontulate`](#prontulate)
- [`slotify`](#slotify)

## girth

Measure collection/text length

```jave
girth(Value) --> <<exact>>
```

Returns item count for table/enumeration/lexis values and rune count for strang values.

### Examples

```jave
allow exact N 2b=2 girth([1, 2, 3]);;

```

### Status

stable

## pront

Print a single value to stdout

```jave
pront(Value)
```

Emits one display-converted value and appends a newline.

### Examples

```jave
pront("hello, jave");;

```

### Status

stable

## prontulate

Builtin formatted output

```jave
prontulate<Template, Args...>
```

Formats the template by replacing directives with provided values and prints the result.

### Notes

- Builtin behavior is self-sufficient and does not require importing Strangs.

### Examples

```jave
prontulate<"Count=%exact", 2>;;

```

### Status

stable

## slotify

Replace first directive in template

```jave
slotify(Template, Value) --> <<strang>>
```

Replaces the first available formatting directive in a template and returns the updated strang.

### Examples

```jave
allow strang Next 2b=2 slotify("A=%exact B=%strang", 7);;

```

### Status

stable

