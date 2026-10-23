# A few numbers and words

The [implementation](./values.go) doubles integers and uppercases strings. The
[executable examples](./values_test.go) run those operations on small slices and
maps, with checked output.

| Scenario | Builder | Input and output |
| --- | --- | --- |
| `sequenceNumbers` | `AsSequenceSliceAction` | `[1 2 3] -> [2 4 6]` in index order |
| `sequenceWords` | `AsSequenceMapAction` | Values `go` and `chain` become `GO` and `CHAIN` under the same keys |
| `parallelNumbers` | `AsParallelSliceAction` | `[1 2 3] -> [2 4 6]`, retaining indexes |
| `parallelNumberMap` | `AsParallelMapAction` | Values under keys `a` and `b` change from `1, 2` to `2, 4` |

The parallel scenarios call a pure double action without shared mutable state.
Map examples print explicit keys, so the expected output does not depend on map
iteration or goroutine scheduling.

## Stopping at an error

`sequenceWithError` refuses to change `2`, returning it with an error. The
`stopOnError` argument determines whether `3` gets processed:

```text
true [2 2 3] true
false [2 2 6] true
```

The executable example prints `true` for the presence of an error in both
cases. When stopping, unprocessed values remain unchanged.

From the `examples` module:

```sh
go test ./collection -v
```
