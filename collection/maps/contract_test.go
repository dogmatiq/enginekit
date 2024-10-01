package maps_test

import (
	"iter"
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/collection/maps"
	"pgregory.net/rapid"
)

type contract[K, V, M any] interface {
	Len() int
	Has(keys ...K) bool
	Get(K) V
	TryGet(K) (V, bool)

	Clone() M
	Merge(M) M
	Select(func(K, V) bool) M
	Project(func(K, V) (K, V, bool)) M

	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
}

type orderedContract[K, V, M any] interface {
	contract[K, V, M]

	Reverse() iter.Seq2[K, V]
	ReverseKeys() iter.Seq[K]
	ReverseValues() iter.Seq[V]
}

type pointer[K, V, M any] interface {
	*M

	Set(K, V)
	Remove(...K)
	Clear()
}

func testMap[
	P pointer[K, int, M],
	M contract[K, int, M],
	K any,
](
	t *testing.T,
	newMap func(...Pair[K, int]) M,
	isEqual func(K, K) bool,
	gen *rapid.Generator[K],
) {
	t.Parallel()
	t.Helper()

	rapid.Check(t, func(t *rapid.T) {
		var (
			subject  M
			expected []Pair[K, int]
		)

		add := func(k K, v int) {
			expected = append(expected, Pair[K, int]{k, v})
		}

		remove := func(k K) {
			for i, p := range expected {
				if isEqual(p.Key, k) {
					expected[i] = expected[len(expected)-1]
					expected = expected[:len(expected)-1]
					return
				}
			}
		}

		replace := func(k K, v int) {
			for i, p := range expected {
				if isEqual(p.Key, k) {
					expected[i].Value = v
					return
				}
			}
		}

		get := func(k K) (int, bool) {
			for _, p := range expected {
				if isEqual(p.Key, k) {
					return p.Value, true
				}
			}
			return 0, false
		}

		drawExistingKey := func(t *rapid.T) K {
			if len(expected) == 0 {
				t.Skip("map is empty")
			}

			var keys []K
			for _, p := range expected {
				keys = append(keys, p.Key)
			}

			return rapid.
				SampledFrom(keys).
				Draw(t, "existing key")
		}

		drawNewKey := func(t *rapid.T) K {
			for {
				k := gen.Draw(t, "new key")
				if _, ok := get(k); !ok {
					return k
				}
			}
		}

		drawValue := func(t *rapid.T) int {
			return rapid.Int().Draw(t, "value")
		}

		for range 3 {
			k := drawNewKey(t)
			v := drawValue(t)
			add(k, v)
			subject = newMap(expected...)
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"add a new key": func(t *rapid.T) {
					k := drawNewKey(t)
					v := drawValue(t)
					P(&subject).Set(k, v)
					add(k, v)
				},
				"overwrite an existing key": func(t *rapid.T) {
					k := drawExistingKey(t)
					v := drawValue(t)
					P(&subject).Set(k, v)
					replace(k, v)
				},
				"remove an existing key": func(t *rapid.T) {
					k := drawExistingKey(t)
					P(&subject).Remove(k)
					remove(k)
				},
				"remove a key that is not in the map": func(t *rapid.T) {
					k := drawNewKey(t)
					P(&subject).Remove(k)
				},
				"check for presence of key": func(t *rapid.T) {
					k := drawExistingKey(t)
					if !subject.Has(k) {
						t.Fatalf("expected %#v to be in the map", k)
					}
				},
				"check for absence of key": func(t *rapid.T) {
					k := drawNewKey(t)
					if subject.Has(k) {
						t.Fatalf("did not expect %#v to be in the map", k)
					}
				},
				"get the value associated with a key": func(t *rapid.T) {
					k := drawExistingKey(t)
					v := subject.Get(k)

					if x, _ := get(k); v != x {
						t.Fatalf("unexpected value for key %#v: got %#v, want %#v", k, v, x)
					}
				},
				"get the value associated with a key that is not in the map": func(t *rapid.T) {
					k := drawNewKey(t)
					v := subject.Get(k)

					if v != 0 {
						t.Fatalf("unexpected value for key %#v: got %#v, want 0", k, v)
					}
				},
				"try to get the value associated with a key": func(t *rapid.T) {
					k := drawExistingKey(t)

					v, ok := subject.TryGet(k)
					if !ok {
						t.Fatalf("expected %#v to be in the map", k)
					}

					if x, _ := get(k); v != x {
						t.Fatalf("unexpected value for key %#v: got %#v, want %#v", k, v, x)
					}
				},
				"try to get the value associated with a key that is not in the map": func(t *rapid.T) {
					k := drawNewKey(t)

					v, ok := subject.TryGet(k)
					if v != 0 {
						t.Fatalf("unexpected value for key %#v: got %#v, want 0", k, v)
					}

					if ok {
						t.Fatalf("did not expect %#v to be in the map", k)
					}
				},
				"clear the map": func(t *rapid.T) {
					P(&subject).Clear()
					expected = nil
				},
				"clone the map": func(t *rapid.T) {
					subject = subject.Clone()
				},
				"merge with disjoint map": func(t *rapid.T) {
					n := rapid.
						IntRange(0, 3).
						Draw(t, "number of elements")

					other := newMap()

					for range n {
						k := drawNewKey(t)
						v := drawValue(t)
						P(&other).Set(k, v)
						add(k, v)
					}

					subject = subject.Merge(other)
				},
				"merge with itself": func(t *rapid.T) {
					subject = subject.Merge(subject)
				},
				"select a subset": func(t *rapid.T) {
					subject = subject.Select(
						func(k K, v int) bool {
							if v%2 == 0 {
								return true
							} else {
								remove(k)
								return false
							}
						},
					)
				},
				"project a subset": func(t *rapid.T) {
					subject = subject.Project(
						func(k K, v int) (K, int, bool) {
							if v%2 == 0 {
								replace(k, v*2)
								return k, v * 2, true
							} else {
								remove(k)
								return k, v, false
							}
						},
					)
				},
				"project with key modifications": func(t *rapid.T) {
					subject = subject.Project(
						func(k K, v int) (K, int, bool) {
							remove(k)
							k = drawNewKey(t)
							add(k, v)
							return k, v, true
						},
					)
				},
				"": func(t *rapid.T) {
					if subject.Len() != len(expected) {
						t.Fatalf("unexpected length: got %d, want %d", subject.Len(), len(expected))
					}

					for k, v := range subject.All() {
						x, ok := get(k)
						if !ok {
							t.Fatalf("map contains an unexpected key: %#v", k)
						}

						if x != v {
							t.Fatalf("unexpected value for key %#v: got %#v, want %#v", k, v, x)
						}
					}

					for k := range subject.Keys() {
						if _, ok := get(k); !ok {
							t.Fatalf("map contains an unexpected key: %#v", k)
						}
					}

					values := slices.Sorted(subject.Values())
					var expectedValues []int
					for _, p := range expected {
						expectedValues = append(expectedValues, p.Value)
					}
					slices.Sort(expectedValues)

					if len(values) != len(expectedValues) {
						t.Fatalf("unexpected number of values: got %d, want %d", len(values), len(expectedValues))
					}

					for i, want := range expectedValues {
						got := values[i]
						if got != want {
							t.Fatalf("unexpected value at index %d: got %d, want %d", i, got, want)
						}
					}
				},
			},
		)
	})
}

func testOrderedMap[
	P pointer[K, int, M],
	M orderedContract[K, int, M],
	K any,
](
	t *testing.T,
	newMap func(...Pair[K, int]) M,
	cmp func(K, K) int,
	gen *rapid.Generator[K],
) {
	t.Helper()

	testMap[P](
		t,
		newMap,
		func(a, b K) bool { return cmp(a, b) == 0 },
		gen,
	)

	t.Run("keys are iterated in order", func(t *testing.T) {
		t.Parallel()

		rapid.Check(t, func(t *rapid.T) {
			subject := newMap()

			n := rapid.
				IntRange(0, 10).
				Draw(t, "number of elements")

			for range n {
				k := gen.Draw(t, "key")
				v := rapid.Int().Draw(t, "value")
				P(&subject).Set(k, v)
			}

			var prev *K
			for k := range subject.All() {
				if prev != nil {
					if cmp(k, *prev) <= 0 {
						t.Fatalf("keys are not in order: %#v <= %#v", k, prev)
					}
				}

				prev = &k
			}

			prev = nil
			for k := range subject.Reverse() {
				if prev != nil {
					if cmp(k, *prev) >= 0 {
						t.Fatalf("keys are not in reverse order: %#v >= %#v", k, prev)
					}
				}
				prev = &k
			}

			prev = nil
			for k := range subject.Keys() {
				if prev != nil {
					if cmp(k, *prev) <= 0 {
						t.Fatalf("keys are not in order: %#v <= %#v", k, prev)
					}
				}
				prev = &k
			}

			prev = nil
			for k := range subject.ReverseKeys() {
				if prev != nil {
					if cmp(k, *prev) >= 0 {
						t.Fatalf("keys are not in reverse order: %#v >= %#v", k, prev)
					}
				}
				prev = &k
			}

			forwardValues := slices.Collect(subject.Values())
			reverseValues := slices.Collect(subject.ReverseValues())

			i := 0
			for k, v := range subject.All() {
				if forwardValues[i] != v {
					t.Fatalf("unexpected value for key %#v: got %#v, want %#v", k, v, forwardValues[i])
				}
				if reverseValues[len(reverseValues)-1-i] != v {
					t.Fatalf("unexpected value for key %#v (reverse): got %#v, want %#v", k, v, reverseValues[len(reverseValues)-1-i])
				}
				i++
			}
		})
	})
}
