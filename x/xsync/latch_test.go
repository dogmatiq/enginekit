package xsync

import (
	"testing"
	"testing/synctest"
)

func TestLatch(t *testing.T) {
	t.Run("latch set after waiters are waiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var latch Latch

			if latch.IsSet() {
				t.Fatal("did not expect latch to be set")
			}

			waitUnblocked := false
			go func() {
				latch.Wait()
				waitUnblocked = true
			}()

			chanUnblocked := false
			go func() {
				<-latch.Chan()
				chanUnblocked = true
			}()

			synctest.Wait()

			latch.Set()

			if !latch.IsSet() {
				t.Fatal("expected latch to be set")
			}

			synctest.Wait()

			if !waitUnblocked {
				t.Fatal("expected Wait() to unblock")
			}

			if !chanUnblocked {
				t.Fatal("expected Chan() to unblock")
			}
		})
	})

	t.Run("latch set before waiters are waiting", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var latch Latch
			latch.Set()

			waitUnblocked := false
			go func() {
				latch.Wait()
				waitUnblocked = true
			}()

			chanUnblocked := false
			go func() {
				<-latch.Chan()
				chanUnblocked = true
			}()

			synctest.Wait()

			if !waitUnblocked {
				t.Fatal("expected Wait() to unblock")
			}

			if !chanUnblocked {
				t.Fatal("expected Chan() to unblock")
			}
		})
	})

	t.Run("setting an already-set latch is a no-op", func(t *testing.T) {
		var latch Latch
		latch.Set()
		latch.Set()

		if !latch.IsSet() {
			t.Fatal("expected latch to be set")
		}
	})
}
