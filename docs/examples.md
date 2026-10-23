# Small API examples

These examples use integers, short strings, and small collections. They import
only the Go standard library and `chain`; no external service, file, timer, or
random input is needed. Each example includes an `Output:` comment that Go
checks when it runs the example.

From the repository root:

```sh
go test -run '^Example' -v .
```

To run just one, use its full name:

```sh
go test -run '^ExampleAsRetryableAction$' -v .
```

The files use the external `chain_test` package, so they demonstrate the public
API without access to implementation details. They are also Go documentation
examples associated with their named API. All helper functions and state needed
by an example are defined inside it.

## Actions and branches

| API | Example | What to observe |
| --- | --- | --- |
| `NewSimpleAction`, `Action.Name`, `Action.Run` | [ExampleNewSimpleAction](../action_example_test.go) | Doubling `3` returns `6`; the action keeps its name |
| `NewSimpleBranchAction`, `BranchAction.Run`, `BranchAction.NextDirection` | [ExampleNewSimpleBranchAction](../action_example_test.go) | A nil run function passes `-2` through; the selector returns `negative` |
| `AsBranchAction`, `BranchAction.Name`, `BranchAction.Directions` | [ExampleAsBranchAction](../action_example_test.go) | Incrementing `3` gives `4`, so the selector chooses `even` from the output |

The branch examples call `Run` and `NextDirection` separately to show their
contract. A workflow makes these calls for you, and skips `NextDirection` when
the action's `Run` fails. The [Collatz walkthrough](../examples/branch/README.md)
shows complete workflows with branch destinations.

## Workflows and plans

| API | Example | What to observe |
| --- | --- | --- |
| `NewWorkflow`, `Workflow.Name`, `Workflow.Run` | [ExampleNewWorkflow](../workflow_example_test.go) | `2` becomes `3`, then `6`, in constructor order |
| `Workflow.RunAt` | [ExampleWorkflow_RunAt](../workflow_example_test.go) | Start at the double action: `3` becomes `6` without incrementing |
| `Workflow.SetRunPlan`, `RunPlan` | [ExampleWorkflow_SetRunPlan](../workflow_example_test.go) | Skip double: `0` becomes `1`, then `11` |
| `Workflow.ValidateGraph` | [ExampleWorkflow_ValidateGraph](../workflow_example_test.go) | A linear graph is valid; adding a back edge creates a detected cycle |
| `SuccessOnlyPlan` | [ExampleSuccessOnlyPlan](../workflow_example_test.go) | Success selects double; failure and abort terminate |
| `DefaultPlan` | [ExampleDefaultPlan](../workflow_example_test.go) | A failed increment still passes `1` to the add-ten handler; output is `11` and `errors.Is` preserves the failure |
| `DefaultPlanWithAbort` | [ExampleDefaultPlanWithAbort](../workflow_example_test.go) | Inspect the distinct destinations for success, failure, and abort |
| `TerminationPlan` | [ExampleTerminationPlan](../workflow_example_test.go) | `3` becomes `4`, then `8`; the double leaf terminates before add-ten |
| `Terminate` | [ExampleTerminate](../workflow_example_test.go) | Set termination as individual failure and abort destinations in a literal plan |

`TerminationPlan[T]()` applies to a member's entire plan through `SetRunPlan`.
`Terminate[T]()` is a destination for one direction and must not be passed as a
workflow member. The [routing guide](./action-patterns.md#failure-and-abort-routing)
explains error retention and omitted directions.

## Retries and best-effort actions

| API | Example | What to observe |
| --- | --- | --- |
| `AsRetryableAction` | [ExampleAsRetryableAction](../error_control_example_test.go) | The first two attempts fail after adding one; rollback subtracts it, and the third attempt succeeds with `1` |
| `SkipRollback` | [ExampleSkipRollback](../error_control_example_test.go) | Without rollback, two attempts carry the value from `0` to `1` to `2` |
| `AsBestEffortAction` | [ExampleAsBestEffortAction](../error_control_example_test.go) | Fallback sees the original `5` and error; the wrapper returns output `6` with no error |

The retry budget includes the first attempt. These examples use a local attempt
counter to make failures repeatable, with no sleep or network operation. The
best-effort example prints the fallback's input so it can be distinguished from
the action's returned output.

## Slices and maps

| API | Example | What to observe |
| --- | --- | --- |
| `AsSequenceSliceAction` | [ExampleAsSequenceSliceAction](../collection_example_test.go) | `[1 2 3]` becomes `[2 4 6]` in index order |
| `AsSequenceSliceAction` with `stopOnError` | [ExampleAsSequenceSliceAction_stopOnError](../collection_example_test.go) | Refusing `2` gives `[2 2 3]` when stopping, or `[2 2 6]` when continuing; both return an error |
| `AsSequenceMapAction` | [ExampleAsSequenceMapAction](../collection_example_test.go) | Values `go` and `chain` become `GO` and `CHAIN` under the same keys |
| `AsParallelSliceAction` | [ExampleAsParallelSliceAction](../collection_example_test.go) | Concurrent doubling still returns `[2 4 6]` at the original indexes |
| `AsParallelMapAction` | [ExampleAsParallelMapAction](../collection_example_test.go) | Concurrent doubling preserves the keys and produces values `2` and `4` |

The parallel actions use a pure double function with no shared mutable state.
Map examples print explicit keys so the expected output does not depend on map
iteration or goroutine scheduling. See [collection processing](./action-patterns.md#collection-processing)
for error and panic behavior.

## Type adaptation

| API | Example | What to observe |
| --- | --- | --- |
| `AdaptAction` | [ExampleAdaptAction](../adapter_example_test.go) | Increment the `Value` field from `2` to `3` while retaining `Label: "keep"` |

The getter extracts the integer, and the setter writes the result back into the
same aggregate value. The [adapter package](../examples/adapter) shows combining
more than one adapted action in a workflow.
