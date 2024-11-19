# Railway: A Flexible Action Pipeline Library

Railway is a Go package designed for building and managing **directed acyclic graph (DAG)**-based pipelines and actions. Inspired by the principles of [**Railway Oriented Programming**](https://fsharpforfunandprofit.com/rop/), it allows developers to define workflows with clear control flow, branching, and error handling.

## Features

- **Action and Pipeline Composition**  
  Create modular, reusable units of work (`Action`) and orchestrate them into robust workflows using `Pipeline`.

- **DAG-based Execution Plans**  
  Ensure predictable execution and robust validation with built-in cycle detection to enforce acyclic workflows.

- **Nested Pipelines**  
  Use `Pipeline`s as `Action`s within other pipelines, enabling modular and hierarchical workflow designs.

- **Conditional Branching**  
  Support for `Success`, `Failure`, `Abort` and custom direction-based branching within your execution flows.

- **AggregateAction Support**  
  Simplify the orchestration of complex workflows by combining `Action`s or `Pipeline`s of different types into a unified control flow using AggregateAction.

## Getting Started

### Installation

Railway package is part of dago module. Install the package using:

```bash
go get github.com/JSYoo5B/dago
```
