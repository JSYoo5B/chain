# Action Patterns

This document describes patterns that are useful once a workflow needs more
than a simple linear chain.

The [small API examples](./examples.md) provide complete, executable versions
using arithmetic and strings, with expected output for every example.

## Branching

Use `AsBranchAction` when an existing action should also decide the next
direction.

```go
check := chain.AsBranchAction(
    baseAction,
    func(ctx context.Context, output State) (string, error) {
        if output.Ready {
            return "ready", nil
        }
        return "pending", nil
    },
    "ready",
    "pending",
)
```

Use `NewSimpleBranchAction` when the branch node does not need a separately
defined base action.

Custom directions must be declared by the branch action so the workflow can
configure their run plans. `NextDirection` runs after `Run` succeeds and receives
the action's output. When `Run` returns an error, the workflow uses `Failure`
without calling `NextDirection`. A recovered panic uses `Abort`.

[Action and branch examples](./examples.md#actions-and-branches) show both a
pass-through sign check and an increment followed by an even/odd check.

The DSL uses the same direction strings in its `branches:` section. See the
[DSL guide](./dsl.md#connections) and the
[Collatz comparison](../examples/branch/README.md#dsl-version).

## Failure and Abort Routing

For a workflow with multiple members, `Run` starts at the first member and uses
the same routing loop as `RunAt`. Built-in directions are available on every
node, and `BranchAction` can declare additional directions. `SetRunPlan` fills
omitted supported directions with termination; it replaces the previous plan.

| Plan helper | Success | Failure | Abort |
| --- | --- | --- | --- |
| `SuccessOnlyPlan(next)` | `next` | terminate | terminate |
| `DefaultPlan(next, handler)` | `next` | `handler` | terminate |
| `DefaultPlanWithAbort(next, handler, abortHandler)` | `next` | `handler` | `abortHandler` |
| `TerminationPlan[T]()` | terminate | terminate | terminate |

An abort handler can run when configured. Routing through a failure or abort
handler does not clear the original error: routed execution joins action and
selection errors, so callers can use `errors.Is` to inspect their causes. A
selector may also return a direction together with an error; that direction is
followed and the error is retained.

If a selector returns a direction absent from the plan, routed execution stops
with a selection error. `Run` on a single-member workflow runs that one action
once and returns its result without looking up a next action.

Use `ValidateGraph` explicitly for handwritten workflows. It checks directed
cycles and weak connectivity across all configured directions. The DSL generator
performs these structural checks before emitting Go.

## Best-Effort Work

Use `AsBestEffortAction` for non-critical work that should not move the workflow
through the `Failure` direction.

```go
notify := chain.AsBestEffortAction(
    sendNotification,
    func(ctx context.Context, input State, err error) {
        logger.Warnf("notification failed: %v", err)
    },
)
```

The wrapped action's output is still returned. The fallback receives the original
input, and only the returned error is suppressed. This wrapper does not recover
a panic raised by the wrapped action or fallback.

## Retries and Rollback

Use `AsRetryableAction` when an action should retry before the workflow handles
the failure.

```go
save := chain.AsRetryableAction("save", saveOnce, rollbackSave, 3)
```

The numeric argument is the total attempt budget, including the first call.
Here, the action gets at most three attempts and rollback runs at most twice.
Each attempt uses the output of the preceding action or rollback. Rollback does
not run after the final failed attempt. A rollback error stops retries and is
joined with the main action's error. If rollback is not needed, pass
`SkipRollback[T]()`.

The [retry examples](./examples.md#retries-and-best-effort-actions) use an add-one
action and subtract-one rollback to show exactly which value reaches each
attempt. The same section shows a fallback receiving the original input.

## Collection Processing

Sequential slices visit items by index. `AsSequenceSliceAction` accepts a
`stopOnError` argument: `true` stops at the first error, while `false` continues
and joins errors. Sequential maps process one value at a time, but Go map
iteration order and the resulting error order are unspecified. Both sequential
wrappers stop on panic and keep unprocessed values unchanged.

Parallel actions process independent items using concurrent calls to the same
wrapped action. The wrapped action must support those concurrent calls. Slice
indexes and map keys are preserved in the output, and other items continue when
one item fails or panics. Joined errors follow result arrival order and should
not be used for control logic.

The [slice and map examples](./examples.md#slices-and-maps) double numbers and
uppercase short strings. A separate example compares both `stopOnError` values.

## Type Adaptation

Use `AdaptAction` when the workflow carries one aggregate value but an action
only operates on one part of it.

```go
workflowAction := chain.AdaptAction(
    stringAction,
    func(state State) string {
        return state.Name
    },
    func(state State, name string) State {
        state.Name = name
        return state
    },
)
```

The adapted action's output is written back through the setter before the
adapted action returns, including when that action returns an error. A runnable
example is in [`examples/adapter`](../examples/adapter).
