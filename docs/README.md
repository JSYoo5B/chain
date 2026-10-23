# Documentation

Start with the [repository README](../README.md) for the action and workflow API.
The guides here describe the implementation on main.

| Guide | Contents |
| --- | --- |
| [Examples](../examples/README.md) | Scenario packages with implementation code, executable examples, and their own READMEs |
| [Workflow DSL](./dsl.md) | Generate Go from `.chain` files, Go parameters and prerequisites, connections, diagnostics, and verification |
| [Action patterns](./action-patterns.md) | Branching, failure and abort routing, retries, best-effort actions, collections, and type adaptation |
| [Collatz walkthrough](../examples/branch/README.md) | Handwritten and generated branching workflows with the same behavior |

The public API is described alongside its implementation in
[`doc.go`](../doc.go), [`Action`](../action.go),
[`Workflow`](../workflow_configure.go), and [`RunPlan`](../workflow_plan.go).
The generator's command-line behavior is implemented in
[`cmd/chain-gen`](../cmd/chain-gen/main.go).
