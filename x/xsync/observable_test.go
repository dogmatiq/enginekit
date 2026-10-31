package xsync_test

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	. "github.com/dogmatiq/enginekit/x/xsync"
	"golang.org/x/sync/errgroup"
)

func TestObservable(t *testing.T) {
	t.Run("notifies multiple readers of changes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var (
				value   Observable[int]
				done    = make(chan struct{})
				want    = []int{0}
				readers errgroup.Group
			)

			go func() {
				for range 10 {
					synctest.Wait()

					n := rand.Int()
					want = append(want, n)
					value.Store(n)
				}

				close(done)
			}()

			for range 10 {
				readers.Go(func() error {
					var got []int

					for {
						v, changed := value.Load()
						got = append(got, v)

						select {
						case <-changed:
						case <-done:
							if !slices.Equal(got, want) {
								return fmt.Errorf("unexpected values: got %v, want %v", got, want)
							}
							return nil
						}
					}
				})
			}

			if err := readers.Wait(); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("does not produce unexpected values with multiple writers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var (
				value   Observable[int]
				stored  sync.Map
				writers sync.WaitGroup
			)

			stored.Store(0, true)

			for range 3 {
				writers.Go(func() {
					for range 100 {
						n := rand.Int()
						stored.Store(n, true)
						value.Store(n)
					}
				})
			}

			done := make(chan struct{})
			go func() {
				writers.Wait()
				close(done)
			}()

			for {
				v, changed := value.Load()

				if _, ok := stored.Load(v); !ok {
					t.Fatalf("unexpected value: got %d", v)
				}

				select {
				case <-changed:
				case <-done:
					return
				}
			}
		})
	})
}

func TestChangeObservable(t *testing.T) {
	t.Run("it changes the value and notifies readers", func(t *testing.T) {
		var value Observable[int]
		_, changed := value.Load()

		if !ChangeObservable(&value, 1) {
			t.Fatal("expected the value to be changed")
		}

		select {
		case <-changed:
		default:
			t.Fatal("expected readers to be notified")
		}

		if v, _ := value.Load(); v != 1 {
			t.Fatalf("unexpected value: got %d, want 1", v)
		}
	})

	t.Run("it does not notify readers when the value is unchanged", func(t *testing.T) {
		var value Observable[int]
		value.Store(1)

		_, changed := value.Load()

		if ChangeObservable(&value, 1) {
			t.Fatal("expected the value not to be changed")
		}

		select {
		case <-changed:
			t.Fatal("expected readers not to be notified")
		default:
		}
	})

	t.Run("it treats the zero value as the initial value", func(t *testing.T) {
		var value Observable[int]
		_, changed := value.Load()

		if ChangeObservable(&value, 0) {
			t.Fatal("expected the value not to be changed")
		}

		select {
		case <-changed:
			t.Fatal("expected readers not to be notified")
		default:
		}
	})

	t.Run("it changes the value only when it differs, with multiple writers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var (
				value   Observable[bool]
				changes atomic.Int64
				writers sync.WaitGroup
			)

			for range 3 {
				writers.Go(func() {
					for range 100 {
						if ChangeObservable(&value, rand.N(2) == 1) {
							changes.Add(1)
						}
					}
				})
			}

			writers.Wait()

			// Each call that reports a change flips the value, so the final
			// value is true if and only if there was an odd number of them.
			got, _ := value.Load()
			want := changes.Load()%2 == 1

			if got != want {
				t.Fatalf("unexpected final value after %d changes: got %t, want %t", changes.Load(), got, want)
			}
		})
	})
}
