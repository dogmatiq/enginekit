package uuidpb_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"github.com/dogmatiq/enginekit/x/xrapid"
	"pgregory.net/rapid"
)

func TestSet(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		var (
			subject  *Set
			expected = map[string]struct{}{}
		)

		t.Repeat(
			map[string]func(*rapid.T){
				"add a new member": func(t *rapid.T) {
					if subject == nil {
						subject = &Set{}
					}

					v := uuidpb.Generate()

					subject.Add(v)
					expected[v.AsString()] = struct{}{}
				},
				"add an existing member": func(t *rapid.T) {
					v := xrapid.SampledFromSeq(maps.Keys(expected)).Draw(t, "existing member")

					subject.Add(uuidpb.MustParse(v))
				},
				"delete a member": func(t *rapid.T) {
					v := xrapid.SampledFromSeq(maps.Keys(expected)).Draw(t, "existing member")
					subject.Delete(uuidpb.MustParse(v))
					delete(expected, v)
				},
				"delete a non-member": func(t *rapid.T) {
					v := uuidpb.Generate()
					subject.Delete(v)
				},
				"clear the set": func(t *rapid.T) {
					subject.Clear()
					clear(expected)
				},
				"": func(t *rapid.T) {
					if subject.Len() != len(expected) {
						t.Fatalf("unexpected length: got %d, want %d", subject.Len(), len(expected))
					}

					want := slices.Sorted(maps.Keys(expected))

					// check Has()
					{
						if subject.Has(uuidpb.Generate()) {
							t.Fatalf("did not expect random value to be in the set")
						}

						for k := range expected {
							if !subject.Has(uuidpb.MustParse(k)) {
								t.Fatalf("expected %q to be in the set", k)
							}
						}
					}

					// check All()
					{
						var got []string

						for v := range subject.All() {
							got = append(got, v.AsString())

							_, ok := expected[v.AsString()]
							if !ok {
								t.Fatalf("set contains an unexpected member: %q", v)
							}
						}

						slices.Sort(got)
						if !slices.Equal(got, want) {
							t.Fatalf("unexpected keys: got %v, want %v", got, want)
						}

						// partial iteration (coverage)
						for range subject.All() {
							break
						}
					}

					// check Clone()
					{
						clone := subject.Clone()

						if clone.Len() != subject.Len() {
							t.Fatalf("unexpected length of cloned set: got %d, want %d", clone.Len(), subject.Len())
						}

						for v := range expected {
							if !clone.Has(uuidpb.MustParse(v)) {
								t.Fatalf("expected %q to be in the cloned set", v)
							}
						}

						if clone != nil {
							v := uuidpb.Generate()
							clone.Add(v)

							if subject.Has(v) {
								t.Fatalf("adding to cloned set modified the original set")
							}
						}
					}
				},
			},
		)
	})
}
