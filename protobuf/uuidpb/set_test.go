package uuidpb_test

import (
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"pgregory.net/rapid"
)

func TestSet(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		set := Set{}
		members := Map[struct{}]{}
		keys := func() []*UUID {
			var keys []*UUID
			for k := range members {
				keys = append(keys, k.AsUUID())
			}
			return keys
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"": func(t *rapid.T) {
					if len(set) != len(members) {
						t.Fatalf("set cardinality is incorrect: got %d, want %d", len(set), len(members))
					}

					sorted := keys()
					slices.SortFunc(sorted, (*UUID).Compare)

					for i, id := range sorted {
						if !set[i].Equal(id) {
							t.Fatalf("set members are out of order at index %d: got %v, want %v", i, set[i], id)
						}
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
					if len(members) == 0 {
						t.Skip("set is empty")
					}

					m := rapid.
						SampledFrom(keys()).
						Draw(t, "member")

					set.Add(m)
				},
				"delete an existing member": func(t *rapid.T) {
					if len(members) == 0 {
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
					if len(members) == 0 {
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
