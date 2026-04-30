package enginetest

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

func runCommandExecutor(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("CommandExecutor", func(t *testing.T) {
		t.Run("func ExecuteCommand()", func(t *testing.T) {
			t.Run("it routes the command to the correct handler", func(t *testing.T) {
				called := make(chan struct{}, 1)

				handler := &stubs.IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("handler", "a7665b1b-14d1-46e7-8667-8bd7252fd059")
						c.Routes(
							dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeA]](),
						)
					},
					HandleCommandFunc: func(context.Context, dogma.IntegrationCommandScope, dogma.Command) error {
						called <- struct{}{}
						return nil
					},
				}

				app := &stubs.ApplicationStub{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("app", "4ea05b58-e949-4aaf-8493-74b3521e8906")
						c.Routes(
							dogma.ViaIntegration(handler),
						)
					},
				}

				executor := setup(t, app)

				err := executor.ExecuteCommand(
					t.Context(),
					&stubs.CommandStub[stubs.TypeA]{},
				)
				if err != nil {
					t.Fatal(err)
				}

				select {
				case <-called:
				case <-time.After(5 * time.Second):
					t.Fatal("timed out waiting for handler to be called")
				}
			})

			t.Run("when the command has no route", func(t *testing.T) {
				t.Run("it panics", func(t *testing.T) {
					handler := &stubs.IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "a7665b1b-14d1-46e7-8667-8bd7252fd059")
							c.Routes(
								dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeA]](),
							)
						},
					}

					app := &stubs.ApplicationStub{
						ConfigureFunc: func(c dogma.ApplicationConfigurer) {
							c.Identity("app", "4ea05b58-e949-4aaf-8493-74b3521e8906")
							c.Routes(
								dogma.ViaIntegration(handler),
							)
						},
					}

					executor := setup(t, app)

					// Execute a TypeB command, for which no route exists in the
					// application. The engine must panic — this is a programming
					// error on the caller's part.
					defer func() {
						if recover() == nil {
							t.Fatal("expected panic when executing a command with no route")
						}
					}()

					err := executor.ExecuteCommand(
						t.Context(),
						&stubs.CommandStub[stubs.TypeB]{},
					)
					if err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("WithEventObserver", func(t *testing.T) {
				t.Run("it blocks until the observer returns satisfied", func(t *testing.T) {
					// Build a causal chain:
					//   1. Aggregate handles CommandA, records EventA
					//   2. Process handles EventA, executes CommandB
					//   3. Integration handles CommandB, records EventB
					//   4. Observer watches for EventB (end of chain)

					aggregate := &stubs.AggregateMessageHandlerStub[*stubs.AggregateRootStub]{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity(
								"aggregate",
								"a7665b1b-14d1-46e7-8667-8bd7252fd059",
							)
							c.Routes(
								dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeA]](),
								dogma.RecordsEvent[*stubs.EventStub[stubs.TypeA]](),
							)
						},
						RouteCommandToInstanceFunc: func(dogma.Command) string {
							return "instance"
						},
						HandleCommandFunc: func(
							_ *stubs.AggregateRootStub,
							s dogma.AggregateCommandScope[*stubs.AggregateRootStub],
							_ dogma.Command,
						) {
							s.RecordEvent(&stubs.EventStub[stubs.TypeA]{})
						},
					}

					process := &stubs.ProcessMessageHandlerStub[*stubs.ProcessRootStub]{
						ConfigureFunc: func(c dogma.ProcessConfigurer) {
							c.Identity(
								"process",
								"1f6cf00b-7efa-48ef-ad0f-07a6e4d4c237",
							)
							c.Routes(
								dogma.HandlesEvent[*stubs.EventStub[stubs.TypeA]](),
								dogma.ExecutesCommand[*stubs.CommandStub[stubs.TypeB]](),
							)
						},
						RouteEventToInstanceFunc: func(
							context.Context,
							dogma.Event,
						) (string, bool, error) {
							return "instance", true, nil
						},
						HandleEventFunc: func(
							_ context.Context,
							_ *stubs.ProcessRootStub,
							s dogma.ProcessEventScope[*stubs.ProcessRootStub],
							_ dogma.Event,
						) error {
							s.ExecuteCommand(
								&stubs.CommandStub[stubs.TypeB]{},
							)
							return nil
						},
					}

					integration := &stubs.IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity(
								"integration",
								"4aa4b567-8863-4e4c-b659-a9c6471cb558",
							)
							c.Routes(
								dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeB]](),
								dogma.RecordsEvent[*stubs.EventStub[stubs.TypeB]](),
							)
						},
						HandleCommandFunc: func(
							_ context.Context,
							s dogma.IntegrationCommandScope,
							_ dogma.Command,
						) error {
							s.RecordEvent(&stubs.EventStub[stubs.TypeB]{})
							return nil
						},
					}

					app := &stubs.ApplicationStub{
						ConfigureFunc: func(c dogma.ApplicationConfigurer) {
							c.Identity("app", "4ea05b58-e949-4aaf-8493-74b3521e8906")
							c.Routes(
								dogma.ViaAggregate(aggregate),
								dogma.ViaProcess(process),
								dogma.ViaIntegration(integration),
							)
						},
					}

					executor := setup(t, app)

					observed := false
					err := executor.ExecuteCommand(
						t.Context(),
						&stubs.CommandStub[stubs.TypeA]{},
						dogma.WithEventObserver(
							func(_ context.Context, _ *stubs.EventStub[stubs.TypeB]) (bool, error) {
								observed = true
								return true, nil
							},
						),
					)
					if err != nil {
						t.Fatal(err)
					}

					if !observed {
						t.Fatal("expected observer to be called")
					}
				})

				t.Run("when no further events can occur", func(t *testing.T) {
					t.Run("it returns ErrEventObserverNotSatisfied", func(t *testing.T) {
						handler := &stubs.IntegrationMessageHandlerStub{
							ConfigureFunc: func(c dogma.IntegrationConfigurer) {
								c.Identity("handler", "a7665b1b-14d1-46e7-8667-8bd7252fd059")
								c.Routes(
									dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeA]](),
								)
							},
						}

						app := &stubs.ApplicationStub{
							ConfigureFunc: func(c dogma.ApplicationConfigurer) {
								c.Identity("app", "4ea05b58-e949-4aaf-8493-74b3521e8906")
								c.Routes(
									dogma.ViaIntegration(handler),
								)
							},
						}

						executor := setup(t, app)

						err := executor.ExecuteCommand(
							t.Context(),
							&stubs.CommandStub[stubs.TypeA]{},
							dogma.WithEventObserver(
								func(_ context.Context, _ *stubs.EventStub[stubs.TypeA]) (bool, error) {
									return false, nil
								},
							),
						)
						if !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
							t.Fatalf("got %v, want ErrEventObserverNotSatisfied", err)
						}
					})
				})
			})

			t.Run("WithIdempotencyKey", func(t *testing.T) {
				t.Run("it deduplicates commands", func(t *testing.T) {
					var callCount atomic.Int32

					handler := &stubs.IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "a7665b1b-14d1-46e7-8667-8bd7252fd059")
							c.Routes(
								dogma.HandlesCommand[*stubs.CommandStub[stubs.TypeA]](),
								dogma.RecordsEvent[*stubs.EventStub[stubs.TypeA]](),
							)
						},
						HandleCommandFunc: func(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
							callCount.Add(1)
							s.RecordEvent(&stubs.EventStub[stubs.TypeA]{})
							return nil
						},
					}

					app := &stubs.ApplicationStub{
						ConfigureFunc: func(c dogma.ApplicationConfigurer) {
							c.Identity("app", "4ea05b58-e949-4aaf-8493-74b3521e8906")
							c.Routes(
								dogma.ViaIntegration(handler),
							)
						},
					}

					executor := setup(t, app)

					// First call: the handler is invoked and the observer is satisfied.
					err := executor.ExecuteCommand(
						t.Context(),
						&stubs.CommandStub[stubs.TypeA]{},
						dogma.WithIdempotencyKey("dedup-key"),
						dogma.WithEventObserver(
							func(_ context.Context, _ *stubs.EventStub[stubs.TypeA]) (bool, error) {
								return true, nil
							},
						),
					)
					if err != nil {
						t.Fatalf("first call: %v", err)
					}

					// Second call with the same key: the command is deduplicated, so no
					// events are produced and the observer cannot be satisfied.
					//
					// See: https://github.com/dogmatiq/enginekit/blob/main/docs/plans/enginetest-implementation.md
					err = executor.ExecuteCommand(
						t.Context(),
						&stubs.CommandStub[stubs.TypeA]{},
						dogma.WithIdempotencyKey("dedup-key"),
						dogma.WithEventObserver(
							func(_ context.Context, _ *stubs.EventStub[stubs.TypeA]) (bool, error) {
								return true, nil
							},
						),
					)
					if !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
						t.Fatalf("got %v, want ErrEventObserverNotSatisfied on second call", err)
					}

					if n := callCount.Load(); n != 1 {
						t.Fatalf("got %d handler invocations, want 1", n)
					}
				})
			})
		})
	})
}
