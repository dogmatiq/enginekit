# Engine Acceptance Test Suite — Implementation Plan

This document describes how the obligations in [enginetest.md] are
implemented as Go tests. Read that document first — it defines the
obligations and their sources.

## Design Principles

The test suite is a "very thorough smoke test." It verifies every
obligation that is observable from outside the engine. It does not inject
engine-internal failures or test implementation details such as caching
versus replay.

## Package Layout

The `enginetest` package exports a single function:

```go
func Run(t *testing.T, setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor)
```

`Run` delegates to unexported per-group functions. Each group lives in its
own file:

| File                  | Function              | Obligation Group  |
| --------------------- | --------------------- | ----------------- |
| `run.go`              | `Run`                 | entry point       |
| `commandexecutor.go`  | `runCommandExecutor`  | CommandExecutor   |
| `aggregate.go`        | `runAggregate`        | Aggregate         |
| `integration.go`      | `runIntegration`      | Integration       |
| `process.go`          | `runProcess`          | Process           |
| `projection.go`       | `runProjection`       | Projection        |
| `messageownership.go` | `runMessageOwnership` | Message Ownership |
| `messageorder.go`     | `runMessageOrder`     | Message Order     |

Each per-group function receives the same `setup` function and calls
`t.Run` once per obligation. Every `t.Run` gets its own fresh `setup` call,
its own `dogma.Application`, and its own handler stubs. There is no shared
state between subtests.

No registry, scenario struct, or indirection — just direct function calls.

## Implementation Order

Groups are implemented in order of increasing complexity. Each round nails
down style and terminology before moving to the next.

1. CommandExecutor
2. Integration
3. Aggregate
4. Process
5. Projection
6. Message Ownership
7. Message Order

## Message Types

Tests reuse the pre-registered letter-based stubs — `CommandStub[TypeA]`,
`EventStub[TypeB]`, and so on. Each scenario starts from `TypeA` and uses
as many letters as it needs. Letter choice is arbitrary per scenario but
should be mnemonic where practical — `TypeC` for "credit", `TypeW` for
"withdrawal."

Custom type parameters (e.g., `CommandStub[*mutableData]`) are used only
when the letter types are insufficient — for example, message ownership
tests that require mutable content. These custom types are defined in the
file that uses them and registered with a hardcoded UUID.

## Handler Identities

Identity names and keys are inline string literals related to the scenario.
For example, a routing test might use `routes.Identity("target",
"target-key")` and `routes.Identity("other", "other-key")`.

## Assertions

Standard library only — `t.Fatalf("got %v, want %v", got, want)`. Use
`cmp.Diff` only when comparing slices or structures where a plain
`%v`-based failure message would be unreadable.

## Synchronization

### Event chains

When the test needs to wait for a causal chain (command produces event,
event triggers process, process executes command, etc.), it uses
`WithEventObserver`. The observer function accepts the specific event type
it's waiting for and returns `true` when satisfied.

### Handler invocation

When the test needs to confirm a handler was called — and there are no
events to observe — the handler stub writes to a buffered channel. The test
reads from the channel with a generous `select` timeout (e.g., 5 seconds)
as a deadlock guard. `ExecuteCommand` is not guaranteed to block until the
handler returns, so a bare channel read without a timeout is unsafe.

### Negative assertions

Proving that something did _not_ happen — such as a handler not being
called — requires a short wall-clock timeout. There is no causal sync point
that works under unbounded engine concurrency. Where such tests are
necessary, they use a short timeout (e.g., 100 ms). Where the flake risk
is unacceptable, the test is deferred or omitted with a comment explaining
why.

## Timing and Timeouts

Tests that involve `ScheduleTimeout` schedule the timeout at or near
`time.Now()` to make it immediately eligible. Assertions are weak — for
example, "a timeout scheduled in the future is not delivered immediately" —
rather than asserting precise delivery times. If the Dogma docs or ADRs are
ambiguous about a timing guarantee, the test is marked with a `// TODO`.

## Atomicity

Tests verify the "all" direction — record multiple events in one scope,
observe that all appear downstream. The "none" direction (rollback on
failure) is not testable at this abstraction level. Tests include a comment
noting the untested direction.

## Obligation-Specific Notes

### CommandExecutor — WithIdempotencyKey

The handler records events and increments a call counter. The test executes
the same command with the same key twice. After the second call, it asserts
the handler was invoked exactly once and — if `WithEventObserver` is used —
that the second call returns `ErrEventObserverNotSatisfied`.

### Integration — Exactly-Once Events

The handler calls `RecordEvent` and then returns an error on the first
invocation, forcing the engine to retry. On the second invocation, it calls
`RecordEvent` with a different event and returns nil. The test asserts that
only the second invocation's events are persisted.

### Aggregate — Event Replay

The test sends two commands to the same aggregate instance. It verifies
that `ApplyEvent` is called (via the `AggregateRootStub`'s
`AppliedEvents` field) and that the root state is correct when the second
`HandleCommand` is called. The test does not distinguish replay from
caching — only the observable state matters.

### Process — End

Three subtests:

1. After `End`, the instance's pending timeouts are canceled and its state
   is discarded.
2. Calling `ExecuteCommand` or `ScheduleTimeout` after `End` on the same
   scope causes a panic. The handler stub uses `defer`/`recover` inside the
   handler function and flags whether the expected panic occurred.
3. Future events targeting the ended instance are ignored. This is a
   negative assertion — decision on whether to test it with a short timeout
   or omit it entirely is deferred.

### Process — State Persistence

The test sends two events to the same process instance. The first handler
call modifies the `ProcessRootStub.Value` field. The second handler call
asserts that the root's `Value` reflects the first call's changes — proving
the engine round-tripped the state through `MarshalBinary` /
`UnmarshalBinary`.

### Projection — OCC Conflict

Three subtests covering the OCC checkpoint protocol:

1. Normal — handler returns `offset + 1`. The engine advances.
2. Backward — handler returns an offset less than `offset + 1`. The engine
   replays events from the returned offset.
3. Forward — handler returns an offset greater than `offset + 1`. The
   engine skips ahead to the returned offset.

The projection stub tracks call counts and conditionally returns abnormal
checkpoints. It records the offsets and events it receives so the test can
assert replay or skip behavior.

### Message Ownership — Mutable Content

A named type with mutable fields (e.g., `type mutableData struct{ Values
[]string }`) is defined in `messageownership.go` and registered with a
hardcoded UUID.

Tests cover three directions:

1. **Caller to engine.** Execute a command containing a pointer, then
   mutate the pointer after `ExecuteCommand` returns. Send a second command
   to the same aggregate instance — the aggregate's state reflects the
   original value.
2. **Handler to engine.** A handler records an event containing a pointer,
   then mutates the pointer after `RecordEvent`. A downstream consumer
   (projection or process) receives the event with the original value.
3. **Engine to handler.** Two downstream consumers receive the same event.
   The first mutates its copy. The second observes the original value.

## Development Testing

The test suite cannot be tested within this repository — any engine
implementation (such as `testkit`) would create a cyclic dependency.

During development, a `go.work` file in the parent directory links the
local `enginekit` and `testkit` clones:

```
go 1.25.0

use (
    ./enginekit
    ./testkit
)
```

A temporary test file in the `testkit` clone calls `enginetest.Run` against
`testkit`'s engine. The `go.work` file resolves `enginekit` to the local
clone, so changes are tested immediately.

In steady state, `testkit` (and any other engine) imports
`enginekit/enginetest` and runs the suite in its own test files. The
`go.work` file is not committed.

<!-- references -->

[enginetest.md]: enginetest.md
