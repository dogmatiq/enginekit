# Engine Acceptance Test Suite — Implementation Plan

This document describes how the engine behavior requirements in [enginetest.md] are
implemented as Go tests. Read that document first — it defines the
requirements and their sources.

## Design Principles

The test suite is a "very thorough smoke test." It verifies every
engine behavior requirement that is observable from outside the engine. It
does not inject engine-internal failures or test implementation details such
as caching versus replay.

## Package Layout

The public entry point lives in `enginetest/blackbox`. The name `blackbox`
signals intent and makes the eventual move to `dogmatiq/dogma` easy — only
the directory moves, not its name.

Each handler group lives in its own package under
`enginetest/blackbox/internal/`. The `internal/` placement prevents external
code from importing the sub-groups directly. Each sub-package:

- exports a single `Run` function
- defines its own message types with domain-relevant names (no generic
  `TypeA`/`TypeB` stubs)
- registers those types with hardcoded UUIDs
- depends only on `github.com/dogmatiq/dogma` and the standard library
  (the stubs package is a temporary dependency during development in
  `enginekit`; it must be removed before the package moves to `dogma`)

The top-level `blackbox` package exports a single function:

```go
package blackbox

func Run(t *testing.T, setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor)
```

`Run` delegates to per-handler-type packages. Each group lives in its
own file:

| Package                                         | Handler / Feature |
| ----------------------------------------------- | ----------------- |
| `enginetest/blackbox`                           | entry point       |
| `enginetest/blackbox/internal/commandexecutor`  | CommandExecutor   |
| `enginetest/blackbox/internal/integration`      | Integration       |
| `enginetest/blackbox/internal/aggregate`        | Aggregate         |
| `enginetest/blackbox/internal/process`          | Process           |
| `enginetest/blackbox/internal/projection`       | Projection        |
| `enginetest/blackbox/internal/messageownership` | Message Ownership |
| `enginetest/blackbox/internal/messageorder`     | Message Order     |

Each sub-package's `Run` function receives the same `setup` function and
calls `t.Run` once per test. Every `t.Run` gets its own fresh `setup`
call, its own `dogma.Application`, and its own handler implementations.
There is no shared state between subtests.

No registry, scenario struct, or indirection — just direct function calls.

## Implementation Order

Groups are implemented in order of increasing complexity. Each round nails
down style and terminology before moving to the next.

1. [x] CommandExecutor
2. [x] Integration
3. [x] Aggregate
4. [x] Process
5. [x] Projection
6. [ ] Message Ownership
7. [ ] Message Order

## Message Types

Each sub-package defines its own concrete message types with domain-relevant
names. Generic letter-suffix types (`TypeA`, `TypeB`) are not used — names
should reflect the obligation being tested. For example, the integration
package might define `ChargeCard` (command) and `CardCharged` (event) for
the atomicity test, or simply `TriggerCommand` / `CommandTriggered` when no
domain metaphor is natural.

Each type is registered with a hardcoded UUID. Since types are registered
globally, UUIDs must be unique across all sub-packages. By convention each
sub-package owns a dedicated UUID namespace.

No cross-package message type sharing. If two sub-packages need the same
logical message shape, they each define their own type.

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

### Message Content Preservation

Each handler-type test suite includes at least one scenario where the
command and/or event types carry a real field (e.g., a string ID). The
handler asserts that the value it receives matches what the caller sent.
This verifies the engine's marshal/unmarshal cycle preserves content.

Other scenarios may use empty structs when content fidelity is not the
subject under test.

### CommandExecutor — WithIdempotencyKey

The handler records events and increments a call counter. The test executes
the same command with the same key twice. After the second call, it asserts
the handler was invoked exactly once and — if `WithEventObserver` is used —
that the second call returns `ErrEventObserverNotSatisfied`.

### Integration — Exactly-Once Events

**Abandoned.** Testing the "events from a failed invocation are discarded"
obligation requires the engine to retry a failed `HandleCommand` call.
Retry behavior is not uniformly testable at this abstraction level — some
engines (e.g., testkit) surface errors immediately by design, and
production engines may retry after a delay. See the [Untestable Obligations]
section of `enginetest.md` for full reasoning.

The atomicity obligation (all events from a successful invocation are
persisted) is tested instead.

### Aggregate — Event Replay

The test sends two commands to the same aggregate instance. It verifies
that `ApplyEvent` is called and that the root state is correct when the
second `HandleCommand` is called. The test does not distinguish replay
from caching — only the observable state matters.

The root type used in this test implements `MarshalBinary` /
`UnmarshalBinary` with real field values, acting as a snapshot smoke test.
The engine may or may not invoke marshaling; either way, correctness must
hold. The test does not assert that snapshotting occurred — only that the
final state is correct.

### Process — End

Two subtests are implemented:

1. Calling `ExecuteCommand` after `End` on the same scope causes a panic.
   The handler stub uses `defer`/`recover` inside the handler function and
   flags whether the expected panic occurred.
2. Future events targeting the ended instance are ignored. Tested with a
   100 ms wall-clock timeout as a negative assertion.

The third planned subtest — "pending timeouts are canceled after `End`" —
is not implemented. It requires a timeout to have been scheduled before
`End` is called and for the engine to prove it never delivers. That is a
negative assertion with no causal sync point, making it inherently racy
under real engines. It is deferred.

### Process — State Persistence

The test sends two events to the same process instance. The first handler
call modifies a field on the root. The second handler call asserts that
the root reflects those changes — proving the engine round-tripped the
state between invocations.

The root type implements `MarshalBinary` / `UnmarshalBinary` with real
field values. The engine may persist state by snapshotting or by replaying
events through the root; either way correctness must hold. The test does
not assert that snapshotting occurred — only that the final state is
correct.

### Projection — OCC Conflict

**Partially implemented.** Only the normal path is tested:

1. Normal — handler returns `offset + 1`. The engine advances.

The backward and forward conflict paths (returning an offset other than
`offset + 1`) are not testable at this abstraction level. Testkit treats
any non-`(offset + 1)` return value as an error rather than a replay or
skip signal, so the behavior cannot be observed through the
`CommandExecutor` interface.

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

A temporary test file in the `testkit` clone calls `blackbox.Run` against
`testkit`'s engine. The `go.work` file resolves `enginekit` to the local
clone, so changes are tested immediately.

In steady state, `testkit` (and any other engine) imports
`enginekit/enginetest/blackbox` and runs the suite in its own test files.
The `go.work` file is not committed.

<!-- references -->

[enginetest.md]: enginetest.md
[Untestable Obligations]: enginetest.md#untestable-obligations
