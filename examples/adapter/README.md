# Changing one field

The small [number example](./number.go) carries two fields:

```text
{Value: 2, Label: "keep"} -> {Value: 3, Label: "keep"}
```

`incrementNumber` uses `AdaptAction` to apply an `Action[int]` to `Value`. The
getter extracts that integer, and the setter writes back the incremented value
while retaining `Label`. Its [executable example](./number_test.go) checks the
printed result `{3 keep} <nil>`.

## Combining adapted actions

The existing [Pair example](./actions.go) combines integer increment and string
append actions in one workflow. The [tests](./adapter_test.go) show a value
`{10, "f"}` becoming `{13, "foo"}`, and also adapt nested workflows to the same
aggregate state.

From the `examples` module:

```sh
go test ./adapter -v
```
