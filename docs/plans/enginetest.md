# Engine Acceptance Test Suite

This document describes the `enginetest` package — a reusable acceptance
test suite for Dogma engine implementations.

Every test must be justified by a specific statement in the Dogma API
documentation or a Dogma ADR. Obligations that are justified only by an ADR
are marked with **[ADR only]** — this indicates the obligation is not yet
reflected in the API documentation and should be.

## Constraints

These constraints shaped the design and should not be revisited without good
reason:

- **Type-based routing.** Dogma routes messages to handlers based on
  message type. Each command type maps to exactly one handler. This means
  you cannot vary handler behavior per-test by sending different message
  values of the same type — the handler is fixed at configuration time.
- **Global message registry.** Message types are registered globally via
  `RegisterCommand`, `RegisterEvent`, and `RegisterTimeout`. Registration
  is idempotent, so parallel tests registering the same types are safe.
- **In-process execution.** The engine runs in the same process as the
  test. Handler stubs have access to the test's goroutine-local state,
  channels, and closures.
- **Messages must round-trip.** The engine is expected to marshal and
  unmarshal messages ([ADR-28]). Test message types must implement
  `MarshalBinary` and `UnmarshalBinary` correctly. At least one test
  scenario per handler type uses messages and roots with real field values
  to verify that content is preserved through the marshal/unmarshal cycle.
  Empty structs are acceptable for other scenarios where content fidelity
  is not the subject of the test.

## Public API

The package exports a single function:

```go
enginetest.Run(t *testing.T, setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor)
```

The engine author's test file looks like:

```go
func TestEngine(t *testing.T) {
    enginetest.Run(t, func(t *testing.T, app dogma.Application) dogma.CommandExecutor {
        // Start the engine with app.
        // Register t.Cleanup to stop it.
        // Return a CommandExecutor ready for use.
    })
}
```

The setup function is called once per scenario. Each scenario provides its
own `dogma.Application`. The setup function starts the engine, registers
cleanup, and returns an executor. If the engine fails asynchronously,
subsequent `ExecuteCommand` calls should return errors — `t.Fatal` cannot
be called from a non-test goroutine.

## How Scenarios Work

Each test scenario:

1. Constructs purpose-built handler implementations using the `stubs`
   package or plain structs. Handler behavior is written as inline Go code
   that calls scope methods directly.
2. Builds a `dogma.Application` containing those handlers.
3. Calls the setup function to start the engine and get an executor.
4. Executes commands and observes outcomes.

Tests may observe outcomes through any means available to the handler stubs
— not only via Dogma events. Handlers can record calls, capture arguments,
inspect root state, write to channels, or set flags. This is essential for
negative assertions and for verifying engine-to-handler obligations.

Where a test needs to observe Dogma events produced by a causal chain, it
uses `WithEventObserver` on `ExecuteCommand`. `WithEventObserver` is
transitive — it observes events anywhere in the causal chain rooted at the
executed command, not just events recorded by the immediate handler. There
is no polling, no projection-as-infrastructure, no hardcoded timeouts.

Projection-specific tests use a minimal projection handler that writes
received events to a channel. The projection is the test subject, not test
infrastructure.

## Message Types

Message types use a base-type embedding pattern. A base type provides
`MarshalBinary`, `UnmarshalBinary`, `MessageDescription`, and `Validate`.
Each concrete type is a small wrapper:

```go
type CreditAccount struct{ CommandBase }
```

Types are globally registered. Adding a new type requires the struct
definition plus a `RegisterCommand`/`RegisterEvent`/`RegisterTimeout` call
with a UUID. The overhead is small enough that types can be added freely per
scenario.

## Obligations

### CommandExecutor

The `CommandExecutor` is the external entry point into the engine.

- The engine delivers a command to the correct handler based on message
  type routing.
  _Source: `CommandExecutor.ExecuteCommand` doc comment._
- The engine panics when asked to execute a command with no matching
  handler.
  _Source: `HandlesCommand` doc comment — "the engine panics if the
  application has multiple handlers that handle T." The converse — no
  handler — is implied by the routing contract._
- `WithEventObserver` blocks until the observer returns `satisfied == true`
  for an event anywhere in the transitive causal chain rooted at the
  executed command.
  _Source: `WithEventObserver` doc comment; [ADR-30]._
- `WithEventObserver` returns `ErrEventObserverNotSatisfied` when the
  engine determines that no further relevant events can occur and no
  observer has been satisfied.
  _Source: `ErrEventObserverNotSatisfied` doc comment; [ADR-30]._
- `WithIdempotencyKey` deduplicates commands — executing the same command
  with the same key a second time does not produce duplicate side-effects.
  _Source: `WithIdempotencyKey` doc comment; [ADR-29]._

### Aggregate

An aggregate handler enforces invariant business rules by handling commands
and recording events.

- The engine routes each command type to exactly one aggregate handler.
  _Source: `HandlesCommand` doc comment — "the engine panics if the
  application has multiple handlers that handle T."_
- The engine calls `RouteCommandToInstance` to determine which instance the
  command targets.
  _Source: `AggregateMessageHandler.RouteCommandToInstance` doc comment._
- The engine calls `New` to create a blank root, then replays the
  instance's historical events via `ApplyEvent` before calling
  `HandleCommand`.
  _Source: `AggregateMessageHandler.New` and
  `AggregateMessageHandler.HandleCommand` doc comments; [ADR-14]._
- Events recorded within `HandleCommand` are persisted atomically — all or
  none.
  _Source: `AggregateCommandScope.RecordEvent` doc comment — "the engine
  persists all events recorded within this scope in a single atomic
  operation."_
- A second command targeting the same instance observes state that reflects
  events recorded by the first command.
  _Source: `AggregateMessageHandler.HandleCommand` doc comment — "r
  reflects the state of the targeted instance after applying its historical
  events."_
- If the engine supports snapshotting, it must not affect the correctness
  of `HandleCommand`. A root that implements `MarshalBinary` /
  `UnmarshalBinary` with real state must receive the correct state
  regardless of whether the engine restores it from a snapshot or from
  event replay. This is tested as a smoke test — the test does not force
  a snapshot/restore cycle, but the root's marshaling code is reachable
  if the engine chooses to use it.

### Integration

An integration handler connects the application to external systems by
handling commands and optionally recording events.

- The engine routes each command type to exactly one integration handler.
  _Source: `HandlesCommand` doc comment._
- Events recorded within `HandleCommand` are persisted atomically.
  _Source: `IntegrationCommandScope.RecordEvent` doc comment — "the engine
  persists all events recorded within this scope in a single atomic
  operation."_
- The engine calls `HandleCommand` at least once per command. The handler's
  side-effects occur exactly once — repeated invocations of
  `HandleCommand` for the same command do not produce duplicate events.
  _Source: `IntegrationMessageHandler.HandleCommand` doc comment — "the
  engine atomically persists the events recorded by exactly one successful
  invocation."_
  **Untestable at this abstraction level — see [Untestable Obligations].**

### Process

A process handler orchestrates workflows by handling events and executing
commands.

- The engine routes events to process instances via `RouteEventToInstance`.
  _Source: `ProcessMessageHandler.RouteEventToInstance` doc comment._
- When `RouteEventToInstance` returns `ok == false`, the engine ignores the
  event.
  _Source: `ProcessMessageHandler.RouteEventToInstance` doc comment — "if
  ok is false, the handler ignores the event."_
- A process handler can execute commands via `ProcessScope.ExecuteCommand`.
  The engine delivers those commands to the appropriate handler.
  _Source: `ProcessScope.ExecuteCommand` doc comment._
- A process handler can schedule timeouts via
  `ProcessScope.ScheduleTimeout`. The engine delivers the timeout at or
  after the scheduled time.
  _Source: `ProcessScope.ScheduleTimeout` doc comment;
  `ProcessMessageHandler.HandleTimeout` doc comment._
- When a process handler calls `ProcessScope.End`, the engine discards the
  instance's state, cancels its pending timeouts, and ignores future events
  that target the ended instance.
  _Source: `ProcessScope.End` doc comment; [ADR-24]; [ADR-25]._
- The engine panics if the handler calls `ExecuteCommand` or
  `ScheduleTimeout` after calling `End` on the same scope.
  _Source: `ProcessScope.ExecuteCommand` doc comment — "this method panics
  if the process instance has ended"; [ADR-25]._
- The engine persists process state via `ProcessRoot.MarshalBinary` and
  restores it via `ProcessRoot.UnmarshalBinary` between invocations. A
  second event targeting the same instance observes state that reflects
  changes made while handling the first event.
  _Source: `ProcessMessageHandler.HandleEvent` doc comment — "r reflects
  the state of the targeted instance after handling any prior Event or
  Timeout messages"; [ADR-28]._
- A root that implements `MarshalBinary` / `UnmarshalBinary` with real
  state must survive round-tripping without loss of correctness. This is
  tested as a smoke test alongside the state persistence test — the same
  test that verifies state across invocations uses a root that implements
  real marshaling, ensuring the marshal/unmarshal path is exercised if
  the engine chooses to use it.

### Projection

A projection handler builds a read-optimized view of application state by
consuming events.

- The engine delivers events to projection handlers that have a matching
  `HandlesEvent` route.
  _Source: `ProjectionMessageHandler.HandleEvent` doc comment._
- Each event is delivered with a stream ID, an offset, and a checkpoint
  offset.
  _Source: `ProjectionEventScope` doc comment — `StreamID`, `Offset`,
  `CheckpointOffset`; [ADR-26]._
- Events within a single stream are delivered in order.
  _Source: `ProjectionMessageHandler.HandleEvent` doc comment — "the engine
  arranges events on streams such that it delivers all events within a
  single scope in the order they occurred"; [ADR-23]._
- Events from a single aggregate instance are delivered in order, even
  across streams.
  _Source: `ProjectionMessageHandler.HandleEvent` doc comment — "it also
  preserves the order of events from a single aggregate instance, even
  across scopes"; [ADR-23]._
- When the handler returns a checkpoint offset that differs from
  `offset + 1`, the engine resumes delivering events from the returned
  checkpoint offset — the OCC conflict protocol.
  _Source: `ProjectionMessageHandler.HandleEvent` doc comment — "otherwise,
  an OCC conflict has occurred, and the engine resumes delivering events
  starting at cp"; [ADR-26]._

### Message Ownership

Whenever a message crosses the boundary between the engine and the
application, the application owns the message. **[ADR only]**

_Source: [ADR-32]. This obligation is not yet documented on the Dogma
interfaces._

- If the caller mutates a command after passing it to `ExecuteCommand`, the
  engine's copy is unaffected — the handler observes the pre-mutation
  values.
- If a handler mutates a message after passing it to a scope method — such
  as `RecordEvent` or `ExecuteCommand` — the engine's copy is unaffected.
- Messages received by handlers and event observers are independent of any
  value the engine retains. Mutating a received message does not affect
  subsequent deliveries of the same message.

### Message Order

The engine exhibits specific ordering guarantees for events and timeouts.

_Source: [ADR-23]. The guarantees are restated in the
`ProcessMessageHandler.HandleEvent` and `ProjectionMessageHandler.HandleEvent`
doc comments._

- Events recorded within a single scope are delivered in the order they
  were recorded.
- Events from a single aggregate instance are delivered in recorded order,
  even across scopes.
- The relative delivery order of events from different handlers or
  aggregate instances is undefined.
- Timeouts for the same process instance follow a weak total order by
  scheduled time. **[ADR only]**
  _Source: [ADR-23]. This specific guarantee is not restated in the
  `HandleTimeout` doc comment._

## Untestable Obligations

Some obligations are real and must be satisfied by conforming engines, but
cannot be verified by this black-box acceptance suite. They are documented
here so they are not lost, and so future work can revisit them.

### Future: engine capabilities

Some of these obligations may become conditionally testable if engines can
declare which optional behaviors they support. The suite would then enable
the corresponding tests only when the capability is declared — allowing
stronger testing of engines that make stronger guarantees without
penalizing engines that don't.

### Integration — events from a failed invocation are discarded

The `IntegrationMessageHandler.HandleCommand` doc comment states:

> The engine atomically persists the events recorded by exactly one
> successful invocation of this method for each command message.

The "exactly one successful" clause implies that events recorded during a
failed invocation (one that returns a non-nil error) must be discarded.

Testing this requires the engine to retry a failed invocation. Retry is not
uniformly testable at this abstraction level:

- Some engines (e.g., testkit) surface handler errors immediately by
  design. Silently retrying would suppress failures during testing and
  undermine the engine's core value.
- Production engines may retry after a delay, making any black-box test
  either slow, timing-sensitive, or flaky.

The "all events or none" direction of atomicity (a successful invocation
persists all of its events) is tested directly in the Integration suite.
Only the rollback-on-failure direction is untestable here.

<!-- references -->

[ADR-14]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0014-apply-historical-events-to-aggregates.md
[ADR-23]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0023-message-order-guarantees.md
[ADR-24]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0024-permanently-end-processes.md
[ADR-25]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0025-prevent-reverting-ended-processes.md
[ADR-26]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0026-event-stream-based-projection-occ.md
[ADR-28]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0028-binary-marshaling.md
[ADR-29]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0029-retain-command-idempotency-keys.md
[ADR-30]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0030-observe-command-outcomes-via-events.md
[ADR-32]: https://github.com/dogmatiq/dogma/blob/main/docs/adr/0032-message-ownership.md
