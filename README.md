# Chain: A Flexible Action-Pipeline Library

The chain package provides a flexible framework for building and executing sequential workflows (pipelines) of actions. It allows you to define a series of steps, each represented by an Action, which processes inputs and produces outputs. The package supports conditional branching, customizable execution paths, and error handling, enabling complex workflows with minimal boilerplate.

## Features

- Action and Pipeline Composition  
  Create modular, reusable units of work (`Action`) and orchestrate them into robust workflows using `Pipeline`.

- DAG-based Execution Plans  
  Ensure predictable execution and robust validation with built-in cycle detection to enforce acyclic workflows.

- Nested Pipelines  
  Use pipelines as actions within other pipelines, enabling modular and hierarchical workflow designs.

- Conditional Branching  
  Support for `Success`, `Failure`, `Abort` and custom direction-based branching within your execution flows.

- AggregateAction Support  
  Simplify the orchestration of complex workflows by combining `Action`s or `Pipeline`s of different types into a unified control flow using AggregateAction.

## Getting Started

### Installation

To install the chain package, use the following command:

```bash
go get github.com/JSYoo5B/chain
```

### Key Concepts

#### Action

An `Action` represents a single task in the pipeline. Each action can process input data and return output or an error.

```go
type Action[T any] interface {
    Name() string
    Run(ctx context.Context, input T) (output T, err error)
}
```

#### BranchAction

A `BranchAction` extends `Action` and supports conditional branching. It can change the execution flow based on the results of the action, allowing for multiple execution paths.

```go
type BranchAction[T any] interface {
    Name() string
    Run(ctx context.Context, input T) (output T, err error)
    Directions() []string
    NextDirection(ctx context.Context, output T) (direction string, err error)
}
```

#### Pipeline

A `Pipeline` is a sequence of `Action`s executed in order. It orchestrates the flow of data between actions and handles branching, success, error, and abort conditions.

#### ActionPlan

An `ActionPlan` is a map that associates a direction (e.g., success, error, abort) with the next Action to execute, defining the flow of a pipeline.

```go
type ActionPlan[T any] map[string]Action[T]
```

## Examples

Practical examples for using the `chain` package will be added in future updates. Stay tuned!