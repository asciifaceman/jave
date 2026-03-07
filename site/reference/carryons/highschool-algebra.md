---
title: highschool/Algebra Reference
permalink: /reference/carryons/highschool-algebra/
category: carryon-reference
source_kind: source
---

# highschool/Algebra

```jave
install Algebra from highschool/Algebra;;
```

## Table of Contents

- [`LeastExact`](#leastexact)
- [`LeastVag`](#leastvag)
- [`MostExact`](#mostexact)
- [`MostVag`](#mostvag)
- [`Nearlydont`](#nearlydont)
- [`PosiExact`](#posiexact)
- [`PosiVag`](#posivag)
- [`PosidirExact`](#posidirexact)
- [`PosidirVag`](#posidirvag)
- [`Stretch`](#stretch)

## LeastExact

Minimum of two exact values.

```jave
outy seq LeastExact<exact Left, exact Right> --> <<exact>>
```

Chooses the smaller exact argument and returns it.

### Parameters

- `exact Left`: First exact candidate.
- `exact Right`: Second exact candidate.

### Returns

- `exact`: Smaller exact value.

## LeastVag

Minimum of two vag values.

```jave
outy seq LeastVag<vag Left, vag Right> --> <<vag>>
```

Chooses the smaller vag argument and returns it.

### Parameters

- `vag Left`: First vag candidate.
- `vag Right`: Second vag candidate.

### Returns

- `vag`: Smaller vag value.

## MostExact

Maximum of two exact values.

```jave
outy seq MostExact<exact Left, exact Right> --> <<exact>>
```

Chooses the larger exact argument and returns it.

### Parameters

- `exact Left`: First exact candidate.
- `exact Right`: Second exact candidate.

### Returns

- `exact`: Larger exact value.

## MostVag

Maximum of two vag values.

```jave
outy seq MostVag<vag Left, vag Right> --> <<vag>>
```

Chooses the larger vag argument and returns it.

### Parameters

- `vag Left`: First vag candidate.
- `vag Right`: Second vag candidate.

### Returns

- `vag`: Larger vag value.

## Nearlydont

Near-zero check for vag values.

```jave
outy seq Nearlydont<vag Value> --> <<truther>>
```

Uses a fixed epsilon threshold to determine whether a vag is nearly zero.

### Parameters

- `vag Value`: Vag value under review.

### Returns

- `truther`: truther indicating near-zero posture.

## PosiExact

Absolute value for exact numbers.

```jave
outy seq PosiExact<exact Value> --> <<exact>>
```

Returns a non-negative exact by reflecting negative inputs across zero.

### Parameters

- `exact Value`: Exact input to normalize.

### Returns

- `exact`: Non-negative exact magnitude.

### Examples

```jave
allow exact A 2b=2 Algebra.PosiExact<0 - 9>;;
```

## PosiVag

Absolute value for vag numbers.

```jave
outy seq PosiVag<vag Value> --> <<vag>>
```

Returns a non-negative vag by reflecting negative inputs across zero.

### Parameters

- `vag Value`: Vag input to normalize.

### Returns

- `vag`: Non-negative vag magnitude.

### Examples

```jave
allow vag A 2b=2 Algebra.PosiVag<0 - 2.5>;;
```

## PosidirExact

Direction of an exact value.

```jave
outy seq PosidirExact<exact Value> --> <<exact>>
```

Converts an exact into sign direction semantics: -1, 0, or 1.

### Parameters

- `exact Value`: Exact value to classify.

### Returns

- `exact`: Exact direction marker in {-1, 0, 1}.

## PosidirVag

Direction of a vag value.

```jave
outy seq PosidirVag<vag Value> --> <<exact>>
```

Converts a vag into sign direction semantics: -1, 0, or 1.

### Parameters

- `vag Value`: Vag value to classify.

### Returns

- `exact`: Exact direction marker in {-1, 0, 1}.

## Stretch

Linear interpolation across two vag endpoints.

```jave
outy seq Stretch<vag Start, vag End, vag Progress> --> <<vag>>
```

Computes Start + ((End - Start) * Progress).

### Parameters

- `vag Start`: Lower or initial endpoint.
- `vag End`: Upper or target endpoint.
- `vag Progress`: Interpolation ratio, typically 0.0 through 1.0.

### Returns

- `vag`: Interpolated vag value.

