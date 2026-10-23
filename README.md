# Chain

`chain` is a small Go workflow library built around composable `Action`
values. An action receives one value, returns the next value, and can be
connected to other actions through a `Workflow`.

The package focuses on in-process workflow control. It does not try to be a
distributed scheduler or a full parallel DAG engine.

## Installation

```bash
go get github.com/JSYoo5B/chain
```

## Core Concepts

### Action

An `Action` is the basic execution unit.

```go
type Action[T any] interface {
    Name() string
    Run(ctx context.Context, input T) (output T, err error)
}
```

Use `NewSimpleAction` when a function is enough:

```go
double := chain.NewSimpleAction(
    "double",
    func(ctx context.Context, input int) (int, error) {
        return input * 2, nil
    },
)
```

### Workflow

A `Workflow` connects actions and also implements `Action`, so workflows can be
nested inside other workflows.

```go
workflow := chain.NewWorkflow("calculation", action1, action2, action3)

output, err := workflow.Run(context.Background(), input)
```

By default, actions are connected through the `Success` direction in constructor
order. `Failure` and `Abort` terminate unless a custom `RunPlan` is set.

```go
workflow.SetRunPlan(action1, chain.DefaultPlan(action2, errorHandler))
```

### Directions and Run Plans

A `RunPlan` maps a direction to the next action.

```go
type RunPlan[T any] map[string]Action[T]
```

Built-in directions are:

- `Success`
- `Failure`
- `Abort`

`BranchAction` can add custom directions.

## Action Builders

### Basic Actions

- `NewSimpleAction` creates an action from a function.
- `NewSimpleBranchAction` creates a branch action from functions.
- `AsBranchAction` wraps an existing action and adds branching logic.

### Error Control

- `AsRetryableAction` retries a main action and optionally runs rollback before
  the next attempt.
- `SkipRollback` explicitly disables rollback for a retryable action.
- `AsBestEffortAction` suppresses a non-critical action error and optionally
  calls a fallback hook.

### Collection Processing

- `AsSequenceSliceAction` runs an action over a slice sequentially.
- `AsSequenceMapAction` runs an action over a map sequentially.
- `AsParallelSliceAction` runs an action over a slice concurrently.
- `AsParallelMapAction` runs an action over a map concurrently.

Parallel collection actions preserve output positions or keys, but joined error
order follows completion order.

### Type Adaptation

- `AdaptAction` runs an action against a field or sub-value inside a larger
  value.

This is useful when a workflow carries one aggregate state type, while some
actions only operate on one part of that state.

## Validation

Call `Workflow.ValidateGraph` explicitly to check a handwritten workflow before
execution; `Run` does not call it automatically.

It rejects:

- directed cycles
- disconnected workflow graphs

Branching and merge points are allowed. Connectivity is weak: every node must
belong to one component when edge directions are ignored. This does not require
every node to be reachable from the entry node.

## Examples

The [small API examples](./docs/examples.md) cover actions, workflow methods,
plan helpers, retries, collection processing, and type adaptation using simple
arithmetic and strings. They use only the standard library and `chain`, and
their expected outputs are checked by Go tests. From the repository root:

```sh
go test -run '^Example' -v .
```

The branching and type adaptation walkthroughs are Go tests in the separate
[`examples` module](./examples/README.md):

- [`Branching and Collatz`](./examples/branch/README.md) compares handwritten
  workflows with constructors generated from the DSL.
- [`Type adaptation`](./examples/adapter) combines integer and string actions
  in an aggregate state.

From the repository root:

```sh
cd examples
go generate ./branch
go test ./...
```

The root module's `go test ./...` does not include the examples module.

## Workflow DSL

`chain-gen` generates Go workflow constructors from `.chain` files:

```go
package example

import "github.com/JSYoo5B/chain"

workflow Build(first, second chain.Action[int]) Linear[int] {
    nodes: first, second
    success: first --> second --> end
}
```

The declaration generates a Go constructor named `Build` returning
`*chain.Workflow[int]`. Constructor parameters and `prerequisite` statements use
Go syntax. A return inside a closure is allowed; returning from the prerequisite
itself is rejected. Nodes require explicit connections, and omitted supported
directions terminate instead of following declaration order.

From a checkout of this repository:

```sh
go run ./cmd/chain-gen -o path/workflow_chain.gen.go path/workflow.chain
```

See the [DSL guide](./docs/dsl.md) for a complete generate-and-test example,
branch syntax, diagnostics, and input-file protection. The
[documentation index](./docs/README.md) also links the action patterns and
runnable examples.
