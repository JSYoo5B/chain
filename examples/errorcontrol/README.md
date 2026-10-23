# Retrying a small counter

The [implementation](./counter.go) adds one to an integer. The first two calls
return an error, and the third succeeds. A local counter makes this repeatable
without a network call or sleep. The [executable examples](./counter_test.go)
check the final value and attempt counts.

## Retry and rollback

`retryingIncrement(true)` wraps the addition with `AsRetryableAction`. Between
failed attempts, its rollback action subtracts one:

```text
0 -> 1 (failed) -> 0 (rollback)
0 -> 1 (failed) -> 0 (rollback)
0 -> 1 (success)
```

The three-attempt budget includes the first call. The printed result is
`1 3 2 <nil>`: value one, three attempts, two rollbacks, and no final error.

`retryingIncrement(false)` passes `SkipRollback[int]()` instead. Each next
attempt receives the previous output, giving `0 -> 1 -> 2 -> 3`. Its result is
`3 3 0 <nil>`.

## Optional addition

`optionalIncrement` adds one but returns an error. `AsBestEffortAction` records
that error in a fallback, keeps the action's output, and returns a nil error:

```text
fallback: 5 optional step failed
result: 6 <nil>
```

The fallback receives the original input `5`, while the action result is `6`.

From the `examples` module:

```sh
go test ./errorcontrol -v
```
