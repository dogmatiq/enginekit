package collection_test

import (
	"maps"
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/collection"
	"pgregory.net/rapid"
)

func TestEquivalentSet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name string
		A, B Set[elem]
		Want bool
	}{
		{
			"empty",
			NewUnorderedSet[elem](),
			NewUnorderedSet[elem](),
			true,
		},
		{
			"disjoint",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](4, 5, 6),
			false,
		},
		{
			"intersecting",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](3, 4, 5),
			false,
		},
		{
			"superset/subset",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](1, 2),
			false,
		},
		{
			"different types",
			NewUnorderedSet[elem](3, 1, 2),
			NewOrderedSet[elem](2, 3, 1),
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if got := IsEquivalentSet(c.A, c.B); got != c.Want {
				t.Fatalf("unexpected result: got %T, want %T", got, c.Want)
			}

			if got := IsEquivalentSet(c.B, c.A); got != c.Want {
				t.Fatalf("unexpected result: got %T, want %T", got, c.Want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name string
		A, B *UnorderedSet[elem]
		Want *UnorderedSet[elem]
	}{
		{
			"empty",
			NewUnorderedSet[elem](),
			NewUnorderedSet[elem](),
			NewUnorderedSet[elem](),
		},
		{
			"disjoint",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](4, 5, 6),
			NewUnorderedSet[elem](1, 2, 3, 4, 5, 6),
		},
		{
			"intersecting",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](3, 4, 5),
			NewUnorderedSet[elem](1, 2, 3, 4, 5),
		},
		{
			"superset/subset",
			NewUnorderedSet[elem](1, 2, 3),
			NewUnorderedSet[elem](1, 2),
			NewUnorderedSet[elem](1, 2, 3),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if got := Union(c.A, c.B); !IsEquivalentSet(got, c.Want) {
				t.Fatalf("unexpected result: got %v, want %v", got, c.Want)
			}
		})
	}

	t.Run("different types", func(t *testing.T) {
		a := NewUnorderedSet[elem](1, 2, 3)
		b := NewOrderedSet[elem](3, 4, 5)

		want := NewUnorderedSet[elem](1, 2, 3, 4, 5)

		if got := Union(a, b); !IsEquivalentSet(got, want) {
			t.Fatalf("unexpected result: got %v, want %v", got, a)
		}
	})
}

func TestSubset(t *testing.T) {
	t.Parallel()

	want := NewUnorderedSet(2, 4)
	got := Subset(
		NewUnorderedSet(1, 2, 3, 4, 5),
		func(e int) bool {
			return e%2 == 0
		},
	)

	if !IsEquivalentSet(got, want) {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

func TestUnorderedSet(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		set := NewUnorderedSet(1, 50, 110)

		expected := map[int]struct{}{
			1:   {},
			50:  {},
			110: {},
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"": func(t *rapid.T) {
					if set.Len() != len(expected) {
						t.Fatalf("set cardinality is incorrect: got %d, want %d", set.Len(), len(expected))
					}

					for e := range set.Elements() {
						if _, ok := expected[e]; !ok {
							t.Fatalf("unexpected element in set: %d", e)
						}
					}
				},
				"add a new element": func(t *rapid.T) {
					for {
						e := rapid.
							Int().
							Draw(t, "new element")

						if _, ok := expected[e]; !ok {
							set.Add(e)
							expected[e] = struct{}{}
							return
						}
					}
				},
				"duplicate an element": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(slices.Collect(maps.Keys(expected))).
						Draw(t, "existing element")

					set.Add(e)
				},
				"remove an element": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(slices.Collect(maps.Keys(expected))).
						Draw(t, "existing element")

					set.Remove(e)
					delete(expected, e)
				},
				"remove a value that is not in the set": func(t *rapid.T) {
					for {
						e := rapid.
							Int().
							Draw(t, "new element")

						if _, ok := expected[e]; !ok {
							set.Remove(e)
							return
						}
					}
				},
				"check if the set has an element that it should have": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(slices.Collect(maps.Keys(expected))).
						Draw(t, "existing element")

					if !set.Has(e) {
						t.Fatalf("expected %d to be in the set", e)
					}
				},
				"check if the set has an element that it should not have": func(t *rapid.T) {
					for {
						e := rapid.
							Int().
							Draw(t, "new element")

						if _, ok := expected[e]; !ok {
							if set.Has(e) {
								t.Fatalf("did not expect %d to be in the set", e)
							}
							return
						}
					}
				},
				"clear the set": func(t *rapid.T) {
					set.Clear()
					clear(expected)
				},
			},
		)
	})
}

func TestOrderedSet(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		set := NewOrderedSet[elem](1, 50, 110)

		expected := map[elem]struct{}{
			1:   {},
			50:  {},
			110: {},
		}

		expectedInOrder := func() []elem {
			return slices.Sorted(maps.Keys(expected))
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"": func(t *rapid.T) {
					if set.Len() != len(expected) {
						t.Fatalf("set cardinality is incorrect: got %d, want %d", set.Len(), len(expected))
					}

					elements := expectedInOrder()

					i := 0
					for got := range set.Elements() {
						want := elements[i]
						if got != want {
							t.Fatalf("set elements are out of order at index %d: got %d, want %d", i, got, want)
						}
						i++
					}
				},
				"add a new element": func(t *rapid.T) {
					for {
						e := elem(
							rapid.
								Int().
								Draw(t, "new element"),
						)

						if _, ok := expected[e]; !ok {
							set.Add(e)
							expected[e] = struct{}{}
							return
						}
					}
				},
				"duplicate an element": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(expectedInOrder()).
						Draw(t, "existing element")

					set.Add(e)
				},
				"remove an element": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(expectedInOrder()).
						Draw(t, "existing element")

					set.Remove(e)
					delete(expected, e)
				},
				"remove a value that is not in the set": func(t *rapid.T) {
					for {
						e := elem(
							rapid.
								Int().
								Draw(t, "new element"),
						)

						if _, ok := expected[e]; !ok {
							set.Remove(e)
							return
						}
					}
				},
				"check if the set has an element that it should have": func(t *rapid.T) {
					if len(expected) == 0 {
						t.Skip("set is empty")
					}

					e := rapid.
						SampledFrom(expectedInOrder()).
						Draw(t, "existing element")

					if !set.Has(e) {
						t.Fatalf("expected %d to be in the set", e)
					}
				},
				"check if the set has an element that it should not have": func(t *rapid.T) {
					for {
						e := elem(
							rapid.
								Int().
								Draw(t, "new element"),
						)

						if _, ok := expected[e]; !ok {
							if set.Has(e) {
								t.Fatalf("did not expect %d to be in the set", e)
							}
							return
						}
					}
				},
				"clear the set": func(t *rapid.T) {
					set.Clear()
					clear(expected)
				},
			},
		)
	})
}
