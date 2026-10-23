# Runnable examples

This directory is a separate Go module. Its `go.mod` replaces
`github.com/JSYoo5B/chain` with the repository root, so these examples use the
checked-out implementation.

From the repository root:

```sh
cd examples
go generate ./branch
go test ./...
```

| Package | What it demonstrates | Tests |
| --- | --- | --- |
| [branch](./branch/README.md) | Custom directions, explicit run plans, and handwritten/DSL Collatz workflows | [branch_test.go](./branch/branch_test.go), [dsl_test.go](./branch/dsl_test.go) |
| [adapter](./adapter) | Adapting `Action[int]` and `Action[string]` into one `Action[Pair]` workflow | [adapter_test.go](./adapter/adapter_test.go) |

The branch package commits both its [DSL source](./branch/collatz.chain) and
[generated Go](./branch/collatz_chain.gen.go). Its
[`go:generate` directive](./branch/dsl.go) invokes the generator from the root
module through the local replacement. Regenerate after changing the DSL source,
then inspect the generated diff and run the tests.

To run just the handwritten and generated branching examples, from this module:

```sh
go test ./branch -run 'TestBranchingWorkflow|TestGenerated' -v
```

Root `go test ./...` does not include this module. The
[DSL guide](../docs/dsl.md#examples-and-verification) lists both modules' checks.
