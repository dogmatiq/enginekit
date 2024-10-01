package sets_test

import (
	"iter"
	"testing"

	"pgregory.net/rapid"
)

type contract[T, S any] interface {
	Len() int
	Has(members ...T) bool
	IsEqual(S) bool
	IsSuperset(S) bool
	IsSubset(S) bool
	IsStrictSuperset(S) bool
	IsStrictSubset(S) bool

	Clone() S
	Union(S) S
	Select(func(T) bool) S

	All() iter.Seq[T]
}

type orderedContract[T, S any] interface {
	contract[T, S]
	Reverse() iter.Seq[T]
}

type pointer[T, S any] interface {
	*S

	Add(members ...T)
	Remove(members ...T)
	Clear()
}

func testSet[
	P pointer[T, S],
	S contract[T, S],
	T any,
](
	t *testing.T,
	newSet func(...T) S,
	isEqual func(T, T) bool,
	pred func(T) bool,
	gen *rapid.Generator[T],
) {
	t.Parallel()
	t.Helper()

	rapid.Check(t, func(t *rapid.T) {
		var (
			subject  S
			expected []T
		)

		add := func(m T) {
			expected = append(expected, m)
		}

		remove := func(m T) {
			for i, v := range expected {
				if isEqual(m, v) {
					expected[i] = expected[len(expected)-1]
					expected = expected[:len(expected)-1]
					return
				}
			}
		}

		has := func(m T) bool {
			for _, v := range expected {
				if isEqual(m, v) {
					return true
				}
			}
			return false
		}

		drawMember := func(t *rapid.T) T {
			if len(expected) == 0 {
				t.Skip("set is empty")
			}

			return rapid.
				SampledFrom(expected).
				Draw(t, "existing member")
		}

		drawNonMember := func(t *rapid.T) T {
			for {
				m := gen.Draw(t, "non-member")
				if !has(m) {
					return m
				}
			}
		}

		for range 3 {
			add(drawNonMember(t))
			subject = newSet(expected...)
		}

		t.Repeat(
			map[string]func(*rapid.T){
				"add a non-member": func(t *rapid.T) {
					m := drawNonMember(t)
					P(&subject).Add(m)
					add(m)
				},
				"re-add an existing member": func(t *rapid.T) {
					m := drawMember(t)
					P(&subject).Add(m)
				},
				"remove a member": func(t *rapid.T) {
					m := drawMember(t)
					P(&subject).Remove(m)
					remove(m)
				},
				"remove a non-member": func(t *rapid.T) {
					m := drawNonMember(t)
					P(&subject).Remove(m)
				},
				"check for membership": func(t *rapid.T) {
					m := drawMember(t)
					if !subject.Has(m) {
						t.Fatalf("expected %#v to be a member", m)
					}
				},
				"check for non-membership": func(t *rapid.T) {
					m := drawNonMember(t)
					if subject.Has(m) {
						t.Fatalf("did not expect %#v to be a member", m)
					}
				},
				"clear the set": func(t *rapid.T) {
					P(&subject).Clear()
					expected = nil
				},
				"clone the set": func(t *rapid.T) {
					subject = subject.Clone()
				},
				"union with disjoint set": func(t *rapid.T) {
					n := rapid.
						IntRange(0, 3).
						Draw(t, "number of members")

					other := newSet()

					for range n {
						m := drawNonMember(t)
						P(&other).Add(m)
						add(m)
					}

					subject = subject.Union(other)
				},
				"union with itself": func(t *rapid.T) {
					subject = subject.Union(subject)
				},
				"select a subset": func(t *rapid.T) {
					subject = subject.Select(
						func(m T) bool {
							if pred(m) {
								return true
							} else {
								remove(m)
								return false
							}
						},
					)
				},
				"": func(t *rapid.T) {
					if subject.Len() != len(expected) {
						t.Fatalf("unexpected length: got %d, want %d", subject.Len(), len(expected))
					}

					for m := range subject.All() {
						if !has(m) {
							t.Fatalf("set contains an unexpected member: %#v", m)
						}
					}

					equivalent := newSet()

					for _, m := range expected {
						if subject.IsEqual(equivalent) {
							t.Fatal("set should not be equal to its subset", subject)
						}

						if !subject.IsSuperset(equivalent) {
							t.Fatal("set should be a superset of its subset", subject)
						}

						if subject.IsSubset(equivalent) {
							t.Fatal("set should not be a subset of its subset", subject)
						}

						if !subject.IsStrictSuperset(equivalent) {
							t.Fatal("set should be a strict superset of its subset", subject)
						}

						if subject.IsStrictSubset(equivalent) {
							t.Fatal("set should not be a strict subset of its subset", subject)
						}

						P(&equivalent).Add(m)
					}

					if !subject.IsEqual(equivalent) {
						t.Fatal("expected the set to be equal to itself", subject)
					}

					if !subject.IsSuperset(equivalent) {
						t.Fatal("expected the set to be a (non-strict) superset of itself", subject)
					}

					if !subject.IsSubset(equivalent) {
						t.Fatal("expected the set to be a (non-strict) subset of itself", subject)
					}

					if subject.IsStrictSuperset(equivalent) {
						t.Fatal("did not expect the set to be a strict superset of itself", subject)
					}

					if subject.IsStrictSubset(equivalent) {
						t.Fatal("did not expect the set to be a strict subset of itself", subject)
					}

					if n := subject.Len(); n != 0 {
						disjoint := newSet()

						for range n {
							m := drawNonMember(t)
							P(&disjoint).Add(m)
						}

						if subject.IsEqual(disjoint) {
							t.Fatal("did not expect the set to be equal to a disjoint set", subject)
						}

						if subject.IsSuperset(disjoint) {
							t.Fatal("did not expect the set to be a superset of a disjoint set", subject)
						}

						if subject.IsSubset(disjoint) {
							t.Fatal("did not expect the set to be a subset of a disjoint set", subject)
						}

						if subject.IsStrictSuperset(disjoint) {
							t.Fatal("did not expect the set to be a strict superset of a disjoint set", subject)
						}

						if subject.IsStrictSubset(disjoint) {
							t.Fatal("did not expect the set to be a strict subset of a disjoint set", subject)
						}
					}
				},
			},
		)
	})
}

func testOrderedSet[
	P pointer[T, S],
	S orderedContract[T, S],
	T any,
](
	t *testing.T,
	newSet func(...T) S,
	cmp func(T, T) int,
	pred func(T) bool,
	gen *rapid.Generator[T],
) {
	t.Helper()

	testSet[P](
		t,
		newSet,
		func(a, b T) bool { return cmp(a, b) == 0 },
		pred,
		gen,
	)

	t.Run("members are iterated in order", func(t *testing.T) {
		t.Parallel()

		rapid.Check(t, func(t *rapid.T) {
			subject := newSet()

			n := rapid.
				IntRange(0, 10).
				Draw(t, "number of members")

			for range n {
				m := gen.Draw(t, "member")
				P(&subject).Add(m)
			}

			var prev *T
			for m := range subject.All() {
				if prev != nil {
					if cmp(m, *prev) <= 0 {
						t.Fatalf("members are not in order: %#v <= %#v", m, prev)
					}
				}
				prev = &m
			}

			prev = nil
			for m := range subject.Reverse() {
				if prev != nil {
					if cmp(m, *prev) >= 0 {
						t.Fatalf("members are not in reverse order: %#v >= %#v", m, prev)
					}
				}
				prev = &m
			}
		})
	})
}
