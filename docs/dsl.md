# Chain DSL

The DSL adds workflow declarations and direction-based connections to Go.
`chain-gen` reads one `.chain` file and emits ordinary Go constructors using
`Workflow[T]` and `RunPlan[T]`.

## Quick start

In an existing Go module, use the current main implementation:

```sh
go get github.com/JSYoo5B/chain@main
```

Save this as `workflow.chain` in a package directory:

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

From that directory, generate its Go constructor:

```sh
go run github.com/JSYoo5B/chain/cmd/chain-gen -o workflow_chain.gen.go workflow.chain
```

This declares `func Build(first, second chain.Action[int]) *chain.Workflow[int]`
with the runtime workflow name `"Linear"`. Use it from ordinary Go. For example,
save this as `workflow_test.go` beside the generated file:

```go
package example

import (
    "context"
    "fmt"

    "github.com/JSYoo5B/chain"
)

func ExampleBuild() {
    first := chain.NewSimpleAction("increment", func(_ context.Context, input int) (int, error) {
        return input + 1, nil
    })
    second := chain.NewSimpleAction("double", func(_ context.Context, input int) (int, error) {
        return input * 2, nil
    })
    workflow := Build(first, second)
    if err := workflow.ValidateGraph(); err != nil {
        panic(err)
    }
    output, err := workflow.Run(context.Background(), 10)
    if err != nil {
        panic(err)
    }
    fmt.Println(output)
    // Output: 22
}
```

Run `go test ./...` in that package or its module root to build the generated
code and execute the example. Generated Go can be committed and used without
running the generator on every build. Regenerate it after changing the `.chain`
source.

## Declarations and Go regions

A file contains a Go package declaration, optional Go imports, and zero or more
workflow declarations. Each workflow has a constructor name, Go parameters,
workflow name, state type, optional prerequisite, node list, and connections.
Top-level helper functions, constants, and types belong in ordinary `.go` files
in the same package.

The first listed node is the entry point. A workflow must declare a nonempty
`nodes:` list and at least one connection. Nodes are comma-separated Go
identifiers. A single node connected directly to `end` is valid.

Constructor parameters and the bracketed state type use Go syntax. Grouped
parameters such as `a, b int`, callbacks such as `callback func(int) error`, and a
final variadic parameter such as `options ...Option` are accepted. The
[parameters fixture](../internal/dsl/compiler/testdata/parameters/parameters.chain)
and its [Go callers](../internal/dsl/compiler/testdata/parameters/parameters_test.go)
demonstrate calls with omitted, explicit, and slice-expanded variadic arguments.

An optional `prerequisite { ... }` block before `nodes:` contains ordinary Go
statements executed when the constructor is called. Its variables remain
available to the node list and generated run plans. Action `Run` methods execute
later when the workflow runs.

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

Here, `makeAction` and `makeOtherAction` would be helpers in the package's Go
files. DSL words such as `nodes`, `success`, and `end` are ordinary Go identifiers
inside Go regions.

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
expressions. Declare custom directions in the branch action's `Directions`
method or the direction list passed to its builder. Each node has at most one
destination per direction. Duplicate directions are errors, including identical
repeated edges, repeated sections, and collisions between built-in sections and
equivalent branch strings.

`end` is a termination target and cannot be listed as a node or used as a source.
All other connection references must appear in `nodes:`. Duplicate node names
and `_` as a node name are rejected.

The graph must be acyclic and weakly connected, following
`Workflow.ValidateGraph`. Weak connectivity ignores edge directions; it does
not require every node to be reachable from the entry. Node declaration order
sets the entry and membership, while arrows set the execution paths. Only the
selected direction is followed at each step.

Every node receives an explicit run plan. Built-in and declared custom
directions omitted from that plan terminate, including `success` on a leaf.
This avoids inheriting `NewWorkflow`'s default sequential connections. An
unrecognized direction returned during routed execution stops with a selection
error. See [failure and abort routing](./action-patterns.md#failure-and-abort-routing)
for error aggregation, panic handling, and single-member execution behavior.

## Command and generation workflow

From this repository's root:

```sh
go run ./cmd/chain-gen -o path/workflow_chain.gen.go path/workflow.chain
```

The command accepts one input file. Omit `-o` to print generated Go to stdout,
or use `-h` for help. The output directory must already exist. It replaces an
output file only after successful parsing, graph analysis, and Go formatting.
Existing output remains intact when those stages fail. Input and output must
be different files, including when linked paths refer to the same file.

For repeated use outside a checkout, an optional installed command is available:

```sh
go install github.com/JSYoo5B/chain/cmd/chain-gen@main
chain-gen -o workflow_chain.gen.go workflow.chain
```

Keep the installed generator and library dependency aligned. The commands above
select main explicitly; a project's ordinary `go run` invocation uses its module
dependency version.

For `go generate`, put a directive in an ordinary Go file in the same directory
as the DSL source. The command runs from that package directory:

```go
//go:generate go run github.com/JSYoo5B/chain/cmd/chain-gen -o workflow_chain.gen.go workflow.chain
```

Generation checks DSL structure and embedded Go syntax. It does not execute the
prerequisite or build the generated package.

Go imports are preserved. The generated file must build as ordinary Go:
undefined symbols, unused imports, and incompatible action types are reported
by the Go toolchain. Runtime constructor checks still apply to the actual action
values, including duplicate instances and direction support.

The generator resolves each state type in a private package-level alias before
emitting its constructor. A parameter named `int`, or a prerequisite variable
that shadows a state type or imported package, cannot change the types used by
the generated run plans.

## Diagnostics

Frontend diagnostics refer to the original `.chain` source. Generated Go carries
`//line` directives so Go build diagnostics also refer to its filename and line.
Paths are relative to the output directory, allowing generated files to be
committed without machine-specific absolute paths. Columns after Go formatting
are not guaranteed to match the original text.

User-written Go line directives are retained as ordinary comments in generated
code so they cannot override the physical source mapping. The frontend extracts
balanced Go regions using `go/scanner`, then uses `go/parser` for their AST and
prerequisite return checks. No ANTLR runtime is required.

## Examples and verification

The [Collatz example](../examples/branch/README.md#dsl-version) compares generated
constructors with handwritten workflows. From the repository root:

```sh
go test ./...
go vet ./...
cd examples
go generate ./branch
go test ./...
go vet ./...
```

The examples directory is a separate Go module; root tests do not run its tests.
A separate example verifies that a node without outgoing connections terminates
instead of falling through to the next declared node. The root test suite also
checks that the committed generated file matches its `.chain` source.

The root tests also build and run temporary Go packages for type and package
shadowing, generated-name collisions, trailing labels, and closures. Negative
build tests verify diagnostics against original source lines. Fixtures under
`internal/dsl/compiler/testdata` compare generated and handwritten workflows
using node execution order, final output, and preserved error causes. They cover
failure and abort routing, quoted and raw branches, selector errors and panics,
and termination for unlisted directions. Separate Go callers exercise grouped
parameters, function callbacks, and both explicit and expanded variadic options.
These fixtures are compiled and executed by `TestGeneratedGoSemantics`; they
are not standalone Go packages until their `.chain` constructors are generated.

The [Linux CI workflow](../.github/workflows/dsl.yml) runs root and example tests,
vet, and regeneration checks using Go 1.26 and stable. The stable job also runs
a bounded compiler fuzz test.
