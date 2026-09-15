# `messaginggrpc.CommandExecutorAPI` — design decision record

A gRPC API for executing arbitrary commands within a Dogma application,
parallel to `eventstreamgrpc` and `configgrpc`. First resident of the new
_messaging plane_ proto package.

## Taxonomy & naming

- Two API planes: **messaging** (moving messages in/out of running apps:
  command execution, event consumption) and **config** (introspection,
  `configgrpc` — untouched).
- New Go package: `grpc/messaginggrpc`. Proto package: `dogma.messaging.v1`
  (the `grpc` suffix is a Go directory convention only).
- File-per-service, named after the service: `commandexecutor.proto` defines
  `CommandExecutorAPI` — the remote projection of `dogma.CommandExecutor`.
- `eventstreamgrpc` was unused in reality, so it has been moved into this
  package as `eventstreamconsumer.proto` / `EventStreamConsumerAPI` (same
  methods, renamed: `ListEventStreams`, `ConsumeEvents`); the old package is
  deleted (**[BC]**, but acceptable since it was never used). Its
  `UnrecognizedEventType` error detail is shared with `CommandExecutorAPI`;
  consolidation happens with the move.
- Edition 2024, `features.field_presence = IMPLICIT` for scalars that don't
  need presence, imports from `uuidpb`/`envelopepb`, per existing protos.

## Service contract

```proto
service CommandExecutorAPI {
  rpc ExecuteCommand(ExecuteCommandRequest)
      returns (stream ExecuteCommandResponse);
}

message ExecuteCommandRequest {
  // The application that is the target of the command.
  dogma.protobuf.UUID application_key = 1;

  // The command's message type ID, per the Dogma message type registry.
  dogma.protobuf.UUID command_type_id = 2;

  // The command's binary representation (dogma.Message.MarshalBinary).
  bytes data = 3 [features.field_presence = IMPLICIT];

  // Optional; maps to dogma.WithIdempotencyKey(). Empty = absent.
  string idempotency_key = 4 [features.field_presence = IMPLICIT];

  // Type IDs of events the client wishes to observe; may be empty.
  repeated dogma.protobuf.UUID observed_event_type_ids = 5;
}

message ExecuteCommandResponse {
  // The engine has taken ownership of the command — the wire equivalent
  // of a nil return from ExecuteCommand. Always the first frame.
  message CommandAccepted {}

  // An event matching a requested type ID was recorded as a result of
  // executing the command.
  message EventRecorded { dogma.protobuf.Envelope envelope = 1; }

  oneof operation {
    CommandAccepted command_accepted = 1;
    EventRecorded   event_recorded   = 2;
  }
}
```

### Semantics

- **Streaming from day one** so `dogma.WithEventObserver` has a faithful wire
  mapping and implementors never build the rpc twice.
- First frame is always `CommandAccepted` (= durably accepted, _not_
  handled/complete). If `observed_event_type_ids` is empty the server sends it
  and closes immediately — a non-observing client is one `Recv()`.
- The observer _predicate_ stays client-side: the client cancels the stream
  when satisfied. There is deliberately no way to ship a predicate over the
  wire.
- **Clean close (OK)** = the engine determined no further relevant events can
  occur (an unsatisfied client-side observer synthesizes
  `dogma.ErrEventObserverNotSatisfied`). There is no explicit `NoMoreEvents`
  frame — gRPC status codes already distinguish clean close from failure.
- **Non-OK close** = outcome unknown; observation cannot resume. Re-execute
  with the same idempotency key if delivery must be confirmed.
- Duplicate idempotency key → `CommandAccepted` (success), never an error;
  idempotent resubmission is the feature working. The server SHOULD deliver
  the relevant events of the logical execution regardless of deduplication,
  but MAY close (OK) immediately if it cannot determine them — best-effort,
  engine-defined. (The local dogma API is silent on this interplay; see
  parked items.)
- No `message_id` in the request: the engine mints envelope identity. No
  extensions/baggage: trace context belongs in gRPC metadata.
- `dogma.WithIdempotencyKey` panics on empty key — server adapters must only
  apply it when `idempotency_key` is non-empty.

### Errors

Typed error details with helper constructors in a plain `.go` file, per the
`eventstreamgrpc/consume.go` pattern:

| Condition                                                         | Code               | Detail message                                          |
| ----------------------------------------------------------------- | ------------------ | ------------------------------------------------------- |
| Unknown application key                                           | `NOT_FOUND`        | `UnrecognizedApplication{application_key}`              |
| Type ID unknown, not a command, or not handled by the target app  | `INVALID_ARGUMENT` | `UnrecognizedCommandType{command_type_id}`              |
| `data` fails to unmarshal                                         | `INVALID_ARGUMENT` | `MalformedCommand{command_type_id}`                                    |
| Decoded command fails validation                                  | `INVALID_ARGUMENT` | `InvalidCommand{command_type_id}` (validation error in status message) |
| Requested `observed_event_type_ids` entry unknown or not an event | `INVALID_ARGUMENT` | `UnrecognizedEventType{event_type_id}`                  |

`MalformedCommand` vs `InvalidCommand` are deliberately distinct: wrong bytes
vs wrong values are different client bugs.

## Scope decisions

- **No listing/self-description RPC.** Discovery is ConfigAPI's job
  (application identities + command Go type names). The Go-type-name ↔
  type-ID bridge is a known gap, parked (see below). Explicit sidestep to
  keep `CommandExecutorAPI` a pure verb.
- **Explicit application targeting** (`application_key` required): a server
  may host multiple apps; routing by type ID alone is a runtime landmine.
- **Bespoke payload fields** (`command_type_id` + `data`), not
  `envelopepb.Message` — avoids a dead `description` field on submission.
- **Protos + generated code + error helpers only.** No client- or server-side
  Go adapters in this package for now. When a client adapter is built
  (e.g. for clikit's remote mode), it must map `EventObserverOption` onto
  `observed_event_type_ids` + a local predicate — and any option it cannot
  honor must be rejected with a clear error, never silently dropped.

## Implementation context

- Registry APIs: `dogma.RegisteredMessageTypeByID(id)` → unmarshal `data`;
  registration is mandatory (an app using unregistered commands is invalid).
- Generation pipeline: protos build via `make` (protoc toolchain under
  `artifacts/protobuf/`, plus `protoc-gen-go`, `protoc-gen-go-grpc` and
  primo — see existing `*_primo.pb.go`). Match `eventstreamgrpc` file layout:
  `commandexecutor.proto`, generated `*.pb.go`/`*_grpc.pb.go`/`*_primo.pb.go`,
  hand-written error helpers + `doc.go`.
- This package IS tagged with enginekit releases (unlike clikit): record the
  addition under `[Unreleased]` in `CHANGELOG.md`.
- Test fixtures: `enginetest/stubs` registered `CommandStub[T]`/`EventStub[T]`
  types and `MessageTypeID[T]()`.

## Parked items (out of scope, recorded so they aren't lost)

- **configpb type-ID bridge**: nothing on the wire maps Go type names
  (configpb vocabulary) to message type IDs (envelope/registry vocabulary).
  Candidate fix: add type IDs to `configpb.Application` (BC-safe field
  addition), benefiting all APIs including eventstream.
- **eventstreamgrpc application gap**: streams still don't identify their
  application; not fixed by the move, remains open.
- **Upstream dogma clarification**: the local API is silent on
  `WithEventObserver` × idempotency-key-duplicate interplay; consider a doc
  clarification in `dogmatiq/dogma`.
- **clikit remote mode**: a `dogma.CommandExecutor` client adapter over this
  API (see `docs/.plans/clikit.md`).
