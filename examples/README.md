# Examples

Each directory is a small scenario with its own Go implementation, executable
examples or tests, and README. The library implementation lives in the repository
root; the scenario code lives here.

| Package | Scenario | Implementation and execution |
| --- | --- | --- |
| [numbers](./numbers/README.md) | Add one, double a number, skip a step, handle failure, and stop a path | [numbers.go](./numbers/numbers.go), [numbers_test.go](./numbers/numbers_test.go) |
| [branch](./branch/README.md) | Sign checks, even/odd selection, and one Collatz step | [small.go](./branch/small.go), [small_test.go](./branch/small_test.go), [collatz.go](./branch/collatz.go) |
| [errorcontrol](./errorcontrol/README.md) | Retry an addition, undo it between attempts, or keep an optional result | [counter.go](./errorcontrol/counter.go), [counter_test.go](./errorcontrol/counter_test.go) |
| [collection](./collection/README.md) | Double a few numbers and uppercase short words | [values.go](./collection/values.go), [values_test.go](./collection/values_test.go) |
| [adapter](./adapter/README.md) | Change one struct field while retaining the others | [number.go](./adapter/number.go), [number_test.go](./adapter/number_test.go) |

The small scenarios use the standard library and `chain`, with no external
service, timer, or random input. Their `Example...` functions print results
checked by Go's `Output:` comments. The branch and adapter packages also retain
their existing workflow tests.

## Run the examples

This directory is a separate Go module. Its `go.mod` replaces
`github.com/JSYoo5B/chain` with the repository root, so the scenarios use the
checked-out library.

From the repository root:

```sh
cd examples
go generate ./branch
go test ./...
```

Run one subject at a time from this module:

```sh
go test ./numbers -run Example -v
go test ./errorcontrol -run Example -v
go test ./collection -run Example -v
```

Root `go test ./...` checks the library and generator; it does not include this
module's scenario tests.

## Generated workflows

The branch package commits its [DSL source](./branch/collatz.chain) and
[generated Go](./branch/collatz_chain.gen.go). Its
[`go:generate` directive](./branch/dsl.go) invokes the generator from the root
module through the local replacement. Regenerate after changing the DSL source,
then inspect the generated diff and run the tests.

The [DSL guide](../docs/dsl.md#examples-and-verification) lists both modules'
checks.
