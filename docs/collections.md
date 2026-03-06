# Collections Guide

This guide covers `table`, `enumeration`, and `lexis` in Jave v0.1.

## Summary

- `table<T>`: ordered indexed collection using `[...]` literals.
- `enumeration<T>`: dynamic list-like collection using `<...>` literals.
- `lexis<K, V>`: map-like collection using `{...}` literals.

Literal forms are exact in v0.1 with no alternate syntax.

## Examples

```jave
allow table<exact> Scores 2b=2 [1, 2, 3];;
allow enumeration<strang> Names 2b=2 <"Ada", "Linus">;;
allow lexis<strang, exact> Ages 2b=2 {"Ada": 36, "Linus": 55};;
```

See `examples/collections/main.jave` and `examples/multi_dimensional_tables/main.jave`.
