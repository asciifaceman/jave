---
title: Builtin Reference
permalink: /reference/builtins/
category: builtin-reference
source_kind: manifest
---

# Builtin Reference

## Table of Contents

- [`DossierAffixStrang`](#dossieraffixstrang)
- [`DossierJotStrang`](#dossierjotstrang)
- [`DossierPeruseStrang`](#dossierperusestrang)
- [`DossierPresent`](#dossierpresent)
- [`Exeunt`](#exeunt)
- [`FeudAt`](#feudat)
- [`FeudGirth`](#feudgirth)
- [`Girth`](#girth)
- [`HomeStead`](#homestead)
- [`Pront`](#pront)
- [`ProntOops`](#prontoops)
- [`Prontulate`](#prontulate)
- [`Slotify`](#slotify)
- [`TrailJunction`](#trailjunction)
- [`TrailNormify`](#trailnormify)

## DossierAffixStrang

Append dossier text as strang

```jave
DossierAffixStrang(Trail, Content)
```

Appends Content to dossier Trail and creates the dossier when absent.

## DossierJotStrang

Write dossier text as strang

```jave
DossierJotStrang(Trail, Content)
```

Writes Content to dossier Trail, replacing existing text when present.

## DossierPeruseStrang

Read dossier text as strang

```jave
DossierPeruseStrang(Trail)
```

Reads a dossier at Trail and returns full text as strang.

## DossierPresent

Check dossier presence

```jave
DossierPresent(Trail)
```

Returns truther indicating whether dossier Trail exists.

## Exeunt

Explicit runtime program exit

```jave
Exeunt(Code)
```

Ends program execution with the requested exit code. Code must be exact in the range 0 through 255.

## FeudAt

Read CLI argument by index

```jave
FeudAt(Index)
```

Returns the runtime program argument at zero-based Index. Fails with deterministic runtime diagnostics on out-of-range access.

## FeudGirth

Count CLI program arguments

```jave
FeudGirth()
```

Returns the number of runtime program arguments passed after the input source or jbin path.

## Girth

Measure collection/text length

```jave
Girth(Value) --> <<exact>>
```

Returns item count for table/enumeration/lexis values and rune count for strang values.

### Examples

```jave
allow exact N 2b=2 Girth([1, 2, 3]);;

```

### Status

stable

## HomeStead

Current working stead trail

```jave
HomeStead()
```

Returns the runtime current working directory trail.

## Pront

Print a single value to stdout

```jave
Pront(Value)
```

Emits one display-converted value and appends a newline.

### Examples

```jave
Pront("hello, jave");;

```

### Status

stable

## ProntOops

Pront one line to stderr

```jave
ProntOops(Message)
```

Writes a single value to stderr using Jave display conversion and a trailing newline.

## Prontulate

Builtin formatted output

```jave
Prontulate<Template, Args...>
```

Formats the template by replacing directives with provided values and prints the result.

### Notes

- Builtin behavior is self-sufficient and does not require importing Strangs.

### Examples

```jave
Prontulate<"Count=%exact", 2>;;

```

### Status

stable

## Slotify

Replace first directive in template

```jave
Slotify(Template, Value) --> <<strang>>
```

Replaces the first available formatting directive in a template and returns the updated strang.

### Examples

```jave
allow strang Next 2b=2 Slotify("A=%exact B=%strang", 7);;

```

### Status

stable

## TrailJunction

Join path trail segments

```jave
TrailJunction(Parts...)
```

Joins one or more trail segments using host platform separators.

## TrailNormify

Normalize a path trail

```jave
TrailNormify(Trail)
```

Cleans redundant separators and dot segments for deterministic trail formatting on the host platform.

