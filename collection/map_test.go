package collection_test

import (
	"maps"
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/collection"
	"pgregory.net/rapid"
)

func TestOrderedMap(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		omap := &OrderedMap[elem, string]{}

		expected := map[elem]string{}
		keysInOrder := func() []elem {
			return slices.Sorted(maps.Keys(expected))
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"": func(t *rapid.T) {
					if omap.Len() != len(expected) {
						t.Fatalf("map length is incorrect: got %d, want %d", omap.Len(), len(expected))
					}

					expectedKeys := keysInOrder()

					i := 0
					for gotKey, gotValue := range omap.Elements() {
						wantKey := expectedKeys[i]
						wantValue := expected[wantKey]

						if gotKey != wantKey {
							t.Fatalf("map elements are out of order at index %d: got %d, want %d", i, gotKey, wantKey)
						}

						if gotValue != wantValue {
							t.Fatalf("map value is incorrect for key %d: got %s, want %s", gotKey, gotValue, wantValue)
						}

						i++
					}

					i = 0
					for gotKey := range omap.Keys() {
						wantKey := expectedKeys[i]

						if gotKey != wantKey {
							t.Fatalf("map keys are out of order at index %d: got %d, want %d", i, gotKey, wantKey)
						}

						i++
					}

					i = 0
					for gotValue := range omap.Values() {
						key := expectedKeys[i]
						wantValue := expected[key]

						if gotValue != wantValue {
							t.Fatalf("map value is incorrect for key %d: got %s, want %s", key, gotValue, wantValue)
						}

						i++
					}
				},
				"set a new key": func(t *rapid.T) {
					for {
						k := elem(
							rapid.
								Int().
								Draw(t, "new key"),
						)

						v := rapid.
							String().
							Draw(t, "value")

						if _, ok := expected[k]; !ok {
							omap.Set(k, v)
							expected[k] = v
							return
						}
					}
				},
				"overwite an existing key": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("map is empty")
					}

					k := rapid.
						SampledFrom(keysInOrder()).
						Draw(t, "existing key")

					v := rapid.
						String().
						Draw(t, "value")

					omap.Set(k, v)
					expected[k] = v
				},
				"remove a key": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("map is empty")
					}

					k := rapid.
						SampledFrom(keysInOrder()).
						Draw(t, "existing key")

					omap.Remove(k)
					delete(expected, k)
				},
				"remove a key that is not in the map": func(t *rapid.T) {
					for {
						k := elem(
							rapid.
								Int().
								Draw(t, "new key"),
						)

						if _, ok := expected[k]; !ok {
							omap.Remove(k)
							return
						}
					}
				},
				"get the value associated with a key": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("map is empty")
					}

					k := rapid.
						SampledFrom(keysInOrder()).
						Draw(t, "existing key")

					got := omap.Get(k)
					want := expected[k]

					if got != want {
						t.Fatalf("unexpected value associated with %d: got %s, want %s", k, got, want)
					}
				},
				"get the value associated with a key that is not in the map": func(t *rapid.T) {
					for {
						k := elem(
							rapid.
								Int().
								Draw(t, "new key"),
						)

						if _, ok := expected[k]; !ok {
							got := omap.Get(k)
							want := ""

							if got != want {
								t.Fatalf("unexpected value associated with %d: got %s, want %s", k, got, want)
							}

							return
						}
					}
				},
				"check if the map has an element that it should have": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("map is empty")
					}

					k := rapid.
						SampledFrom(keysInOrder()).
						Draw(t, "existing key")

					if !omap.Has(k) {
						t.Fatalf("expected %d to be in the map", k)
					}
				},
				"check if the map has an element that it should not have": func(t *rapid.T) {
					for {
						k := elem(
							rapid.
								Int().
								Draw(t, "new key"),
						)

						if _, ok := expected[k]; !ok {
							if omap.Has(k) {
								t.Fatalf("did not expect %d to be in the map", k)
							}
							return
						}
					}
				},
				"clear the map": func(t *rapid.T) {
					omap.Clear()
					clear(expected)
				},
				"clone the map": func(t *rapid.T) {
					omap = omap.Clone()
				},
			},
		)
	})
}
