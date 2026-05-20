# Action Patterns

This document describes patterns that are useful once a workflow needs more
than a simple linear chain.

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
validate run plans.

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

The wrapped action's output is still returned. Only the error is suppressed.
Panics are not suppressed by best-effort handling.

## Retries and Rollback

Use `AsRetryableAction` when an action should retry before the workflow handles
the failure.

```go
save := chain.AsRetryableAction("save", saveOnce, rollbackSave, 3)
```

Rollback runs between attempts and does not run after the final failed attempt.
If rollback is not needed, pass `SkipRollback[T]()`.

## Collection Processing

Use sequence actions when order and deterministic execution are more important
than throughput.

Use parallel actions when each item can be processed independently.

Parallel actions preserve slice indexes or map keys in their outputs. Error
order is based on completion order and should not be used for control logic.

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
adapted action returns.
