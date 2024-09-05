package uuidpb_test

import (
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"pgregory.net/rapid"
)

func TestOrderedSet(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		set := OrderedSet{}
		var members Map[struct{}]

		keys := func() []*UUID {
			return slices.Collect(members.Keys())
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"": func(t *rapid.T) {
					if set.Len() != members.Len() {
						t.Fatalf("set cardinality is incorrect: got %d, want %d", set.Len(), members.Len())
					}

					sorted := keys()
					slices.SortFunc(sorted, (*UUID).Compare)

					i := 0
					for got := range set.All() {
						want := sorted[i]
						if !want.Equal(got) {
							t.Fatalf("set members are out of order at index %d: got %v, want %v", i, got, want)
						}
						i++
					}
				},
				"add a non-member": func(t *rapid.T) {
					m := &UUID{
						Upper: rapid.
							Uint64().
							Draw(t, "non-member (upper)"),
						Lower: rapid.
							Uint64().
							Draw(t, "non-member (lower)"),
					}

					if members.Has(m) {
						t.Skip("already a member")
					}

					set.Add(m)
					members.Set(m, struct{}{})
				},
				"re-add an existing member": func(t *rapid.T) {
					if members.Len() == 0 {
						t.Skip("set is empty")
					}

					m := rapid.
						SampledFrom(keys()).
						Draw(t, "member")

					set.Add(m)
				},
				"delete an existing member": func(t *rapid.T) {
					if members.Len() == 0 {
						t.Skip("set is empty")
					}

					m := rapid.
						SampledFrom(keys()).
						Draw(t, "member")

					set.Delete(m)
					members.Delete(m)
				},
				"delete a non-member": func(t *rapid.T) {
					m := &UUID{
						Upper: rapid.
							Uint64().
							Draw(t, "non-member (upper)"),
						Lower: rapid.
							Uint64().
							Draw(t, "non-member (lower)"),
					}

					if members.Has(m) {
						t.Skip("already a member")
					}

					set.Delete(m)
				},
				"check membership of member": func(t *rapid.T) {
					if members.Len() == 0 {
						t.Skip("set is empty")
					}

					m := rapid.
						SampledFrom(keys()).
						Draw(t, "member")

					if !set.Has(m) {
						t.Fatalf("membership check failed for %v", m)
					}
				},
				"check membership of non-member": func(t *rapid.T) {
					m := &UUID{
						Upper: rapid.
							Uint64().
							Draw(t, "non-member (upper)"),
						Lower: rapid.
							Uint64().
							Draw(t, "non-member (lower)"),
					}

					if members.Has(m) {
						t.Skip("already a member")
					}

					if set.Has(m) {
						t.Fatalf("membership check failed for %v", m)
					}
				},
			},
		)
	})
}
