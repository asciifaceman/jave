# Algebra Carryon Surface

Current exported sequences:

- `PosiExact`: exact absolute value
- `PosiVag`: vag absolute value
- `MostExact`, `LeastExact`: exact comparisons
- `MostVag`, `LeastVag`: vag comparisons
- `PosidirExact`, `PosidirVag`: sign direction as `exact` (-1, 0, 1)
- `Nearlydont`: near-zero check for `vag`
- `Stretch`: linear interpolation for `vag`

Planned naming map (not implemented yet):

- `sqrt` -> `Squirt`
- `round` -> `Erode`
- `max` -> `Most`
- `min` -> `Least`
- `sign` -> `Posidir`
- `almostzero` -> `Nearlydont`
- `isinf` -> `Overlybigly`
- `lerp` -> `Stretch`

Known blockers for not-yet-implemented names:

- `Squirt` requires a square-root primitive or intrinsic bridge.
- `Erode` requires a ratified rounding contract and implementation primitive.
- `Overlybigly` requires stable infinity detection semantics in language/runtime.

Type-family stub direction:

- Future exact families are expected to include width/signedness variants (for example `exactly8`, `exactly32`, `exactlyposi8`, `exactlyposi32`).
- Future vag families should mirror this pattern with names like `vagly32`.
