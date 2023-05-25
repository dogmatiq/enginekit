package enginetest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
	"google.golang.org/protobuf/proto"
)

// SetupFunc is a function that sets up an engine for testing.
type SetupFunc func(SetupParams) SetupResult

// SetupParams are the parameters for a [SetupFunc].
type SetupParams struct {
	App dogma.Application
}

// SetupResult is the result of a call to a [SetupFunc].
type SetupResult struct {
	RunEngine func(ctx context.Context) error
	Executor  dogma.CommandExecutor
}

// RunTests runs acceptance tests against a Dogma engine implementation.
func RunTests(t *testing.T, setup SetupFunc) {
	app := &testapp.App{}
	res := setup(SetupParams{
		App: app,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if d, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, d)
		defer cancel()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		t.Log("dogma engine started")

		err := res.RunEngine(ctx)

		if ctx.Err() == context.Canceled && errors.Is(err, context.Canceled) {
			t.Log("dogma engine exited cleanly")
		} else {
			t.Error("dogma engine exited with unexpected error", err)
			cancel()
		}
	}()

	t.Run("acceptance tests", func(t *testing.T) {
		e := &engine{
			ctx:      ctx,
			App:      app,
			Executor: res.Executor,
		}

		testCommandExecutor(ctx, t, e)
		testIntegration(ctx, t, e)
	})

	cancel()
	<-done
}

// engine is the API used by the tests to interact with the engine-under-test.
type engine struct {
	ctx context.Context

	App      *testapp.App
	Executor dogma.CommandExecutor
}

func (e *engine) ExecuteCommand(t *testing.T, c command) {
	if err := e.Executor.ExecuteCommand(e.ctx, c); err != nil {
		t.Fatal("failed to execute command", err)
	}
}

// ExpectEvent waits for the app to produce the given event.
func (e *engine) ExpectEvent(t *testing.T, expect event) {
	e.ExpectEventMatching(
		t,
		func(actual event) bool {
			return proto.Equal(actual, expect)
		},
	)
}

// ExpectEventMatching waits for the app to produce an event for which the given
// predicate function returns true.
func (e *engine) ExpectEventMatching(
	t *testing.T,
	fn func(event) bool,
) {
	ctx, cancel := context.WithTimeout(e.ctx, 10*time.Second)
	defer cancel()

	if err := e.App.Events.Range(
		ctx,
		func(ev dogma.Event) bool {
			return fn(ev.(event))
		},
	); err != nil {
		t.Fatal("failed while waiting for event", err)
	}
}

type (
	command interface {
		proto.Message
		dogma.Command
	}

	event interface {
		proto.Message
		dogma.Event
	}

	timeout interface {
		proto.Message
		dogma.Timeout
	}
)
