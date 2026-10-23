# Chain DSL

The DSL adds workflow declarations and direction-based connections to Go.
`chain-gen` reads a `.chain` file and emits ordinary Go constructors using
`Workflow[T]` and `RunPlan[T]`.

```go
package example

import "github.com/JSYoo5B/chain"

workflow Build(first, second chain.Action[int]) Linear[int] {
    nodes:
        first, second

    success:
        first --> second --> end
}
```

This declares `func Build(first, second chain.Action[int]) *chain.Workflow[int]`
with the runtime workflow name `"Linear"`. The first listed node is the entry
point. A workflow must declare a nonempty `nodes:` list and at least one
connection. A single node and a connection directly to `end` are valid.

## Go regions

Constructor parameters and the bracketed state type use Go syntax. An optional
`prerequisite { ... }` block before `nodes:` contains ordinary Go statements
executed when the constructor is called. Its variables remain available to the
node list. Action implementations and shared types can live in ordinary `.go`
files in the same package.

Returns that exit the prerequisite block are rejected, including returns inside
its nested `if`, loop, switch, or select statements. Returns inside anonymous
functions and closures are allowed.

```go
prerequisite {
    callback := func(input int) int {
        return input + 1
    }
    first := makeAction(callback)
    second := makeOtherAction()
}
```

DSL words such as `nodes`, `success`, and `end` are ordinary Go identifiers
inside Go regions. The frontend uses `go/scanner` to extract balanced regions,
then `go/parser` to build their AST. No ANTLR runtime is required. Diagnostics
refer to the original `.chain` file, including when Go code contains `//line`
directives.

## Connections

`success:`, `failure:`, and `abort:` select built-in directions. Both `a --> b`
and `b <-- a` describe a connection from `a` to `b`. Chains and repeated
connection sections are allowed. Statements are separated by newlines or
semicolons; arrow characters must be adjacent.

Custom branch directions use double-quoted or raw Go string literals:

```go
branches:
    decide -"even"-> even
    decide -`odd`-> odd
    decide -"done"-> end
```

These strings name directions returned by `BranchAction`; they are not boolean
expressions. Each node can have one destination per direction. Duplicate
directions are errors, including collisions between built-in sections and
equivalent branch strings.

`end` is a termination target and cannot be listed as a node or used as a source.
Node references must be declared. Graphs must be acyclic and weakly connected,
following `Workflow.ValidateGraph`. The compiler keeps the declared entry point
and node order. Every node receives an explicit run plan, so unlisted directions
terminate without inheriting the constructor's default sequential connections.

## Generate and verify

From this repository's root:

```sh
go run ./cmd/chain-gen -o path/workflow_chain.gen.go path/workflow.chain
```

Omit `-o` to print generated Go to stdout. The command replaces an output file
only after successful parsing, graph analysis, and Go formatting. It does not
execute prerequisite code or automatically build the generated package.

Go imports are preserved. The generated file must build as ordinary Go:
undefined symbols, unused imports, and incompatible action types are reported
by the Go toolchain. Runtime action configuration checks still apply.

The runnable Collatz example is in `examples/branch/collatz.chain`. From the
`examples` module:

```sh
go generate ./branch
go test ./...
```

Its generated constructors are compared with the handwritten Collatz workflows.
A separate example verifies that a node without outgoing connections terminates
instead of falling through to the next declared node. The root test suite also
checks that the committed generated file matches its `.chain` source.
