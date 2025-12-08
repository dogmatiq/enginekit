package xsync_test

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"

	. "github.com/dogmatiq/enginekit/x/xsync"
	"golang.org/x/sync/errgroup"
)

func TestFuture(t *testing.T) {

	synctest.Test(t, func(t *testing.T) {
		var f Future[int]

		cancelled, cancel := context.WithCancel(t.Context())
		cancel()

		_, ok := f.Load()
		if ok {
			t.Fatal("did not expect load to succeed")
		}

		_, err := f.Wait(cancelled)
		if err != context.Canceled {
			t.Fatalf("unexpected error: got %v, want %v", err, context.Canceled)
		}

		want := 42

		g, ctx := errgroup.WithContext(t.Context())

		g.Go(func() error {
			v, err := f.Wait(ctx)
			if err != nil {
				return fmt.Errorf("unexpected error waiting for future in background: %v", err)
			}

			if v != want {
				return fmt.Errorf("unexpected value after waiting for future in background: got %v, want %v", v, want)
			}

			return nil
		})

		g.Go(func() error {
			synctest.Wait()

			if !f.Store(want) {
				return fmt.Errorf("expected Store() to return true the first time a value is stored")
			}

			if f.Store(want + 1) {
				return fmt.Errorf("expected Store() to return false the second time a value is stored")
			}

			return nil
		})

		select {
		case <-ctx.Done():
			t.Fatal(context.Cause(ctx))
		case <-f.Ready():
		}

		v, ok := f.Load()
		if !ok {
			t.Fatal("expected load to succeed")
		}

		if v != want {
			t.Fatalf("unexpected value: got %v, want %v", v, want)
		}

		// Should not block despite context being cancelled.
		v, err = f.Wait(cancelled)
		if err != nil {
			t.Fatal(err)
		}

		if v != want {
			t.Fatalf("unexpected value: got %v, want %v", v, want)
		}

		if err := g.Wait(); err != nil {
			t.Fatal(err)
		}
	})
}
