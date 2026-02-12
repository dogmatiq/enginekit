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

				// Block until waiting goroutines exit.
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

				// Block until waiting goroutines exit.
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

	t.Run("Link()", func(t *testing.T) {
		t.Run("unblocks waiters when upstream latch is set", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var upstream, downstreamA, downstreamB Latch

				downstreamA.Link(&upstream)
				downstreamB.Link(&upstream)

				{
					var downstreamC Latch
					downstreamC.Link(&upstream)
				}

				unblockedA := false
				go func() {
					downstreamA.Wait()
					unblockedA = true
				}()

				unblockedB := false
				go func() {
					downstreamB.Wait()
					unblockedB = true
				}()

				// Block until waiting goroutine is blocked on the latch.
				synctest.Wait()

				upstream.Set()

				// Block until waiting goroutine exits.
				synctest.Wait()

				if !unblockedA {
					t.Fatal("expected latch A to unblock")
				}

				if !unblockedB {
					t.Fatal("expected latch B to unblock")
				}
			})
		})
	})
}
