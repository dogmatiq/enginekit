package xsync

import (
	"testing"
	"testing/synctest"
)

func TestLatch(t *testing.T) {
	t.Run("Set()", func(t *testing.T) {
		t.Run("set unblocks existing waiters", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var latch Latch

				if latch.IsSet() {
					t.Fatal("did not expect latch to be set")
				}

				unblockedWait := false
				go func() {
					latch.Wait()
					unblockedWait = true
				}()

				unblockedChan := false
				go func() {
					<-latch.Chan()
					unblockedChan = true
				}()

				// Block until waiting goroutines are blocked on the latch.
				synctest.Wait()

				latch.Set()

				// Block until waiting goroutines exit (or remain blocked if
				// something is wrong).
				synctest.Wait()

				if !latch.IsSet() {
					t.Fatal("expected latch to be set")
				}

				if !unblockedWait {
					t.Fatal("expected Wait() to unblock")
				}

				if !unblockedChan {
					t.Fatal("expected Chan() to unblock")
				}
			})
		})

		t.Run("set unblocks future waiters", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var latch Latch

				// Set the latch before starting the waiting goroutines.
				latch.Set()

				unblockedWait := false
				go func() {
					latch.Wait()
					unblockedWait = true
				}()

				unblockedChan := false
				go func() {
					<-latch.Chan()
					unblockedChan = true
				}()

				// Block until waiting goroutines exit (or remain blocked if
				// something is wrong).
				synctest.Wait()

				if !unblockedWait {
					t.Fatal("expected Wait() to unblock")
				}

				if !unblockedChan {
					t.Fatal("expected Chan() to unblock")
				}
			})
		})

		t.Run("re-setting is a no-op", func(t *testing.T) {
			var latch Latch
			latch.Set()
			latch.Set()

			if !latch.IsSet() {
				t.Fatal("expected latch to be set")
			}
		})
	})
}
