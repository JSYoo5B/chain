# Small arithmetic workflows

The [implementation](./numbers.go) uses only addition and multiplication. The
[executable examples](./numbers_test.go) call those scenarios and check their
printed results.

The starting workflow is `increment -> double`: input `2` becomes `3`, then `6`.
`NewSimpleAction` defines the two operations, and `NewWorkflow` connects them in
order. The example also reads the workflow's name and runs it through `Run`.

| Scenario | What changes | Result |
| --- | --- | --- |
| `double` | Run just one named action | `3 -> 6` |
| `linear` | Run increment, then double | `2 -> 3 -> 6` |
| `runFromDouble` | Use `RunAt` to start at the second member | `3 -> 6` |
| `skipDouble` | Use `SetRunPlan` and a literal `RunPlan` to jump to add-ten | `0 -> 1 -> 11` |
| `stopAtDouble` | Use `DefaultPlan` for the first member and `TerminationPlan` at double | `3 -> 4 -> 8`, then stop |
| `handleFailure` | The failed increment passes its output to add-ten through `DefaultPlan` | `0 -> 1 -> 11`, with the original error retained |
| `resetNegative` | Use `DefaultPlanWithAbort` to reset a negative result to zero | `-3 -> -2 -> 0`; input `3` takes the success path to `8` |
| `cycle` | Add a backward connection using `SuccessOnlyPlan` | `ValidateGraph` reports a cycle; the example does not run it |

`skipDouble` uses `Terminate` as the failure and abort destinations. That is
different from `TerminationPlan`, which stops every supported direction at a
member. A plan changes connections; declaring a member does not force execution
to visit it.

The reset scenario selects `Abort` explicitly from a branch action, so it needs
no panic. The failure scenario uses `errors.Is` to show that running an error
handler does not discard the original failure.

From the `examples` module:

```sh
go test ./numbers -v
```
