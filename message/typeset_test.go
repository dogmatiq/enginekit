package message_test

import (
	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ TypeContainer = TypeSet{}

var _ = Describe("type TypeSet", func() {
	Describe("func NewTypeSet", func() {
		It("returns a set containing the given types", func() {
			Expect(NewTypeSet(
				fixtures.MessageAType,
				fixtures.MessageBType,
			)).To(Equal(TypeSet{
				fixtures.MessageAType: struct{}{},
				fixtures.MessageBType: struct{}{},
			}))
		})
	})

	Describe("func TypesOf", func() {
		It("returns a set containing the types of the given messages", func() {
			Expect(TypesOf(
				fixtures.MessageA1,
				fixtures.MessageB1,
			)).To(Equal(TypeSet{
				fixtures.MessageAType: struct{}{},
				fixtures.MessageBType: struct{}{},
			}))
		})
	})

	Describe("func Has", func() {
		set := TypesOf(
			fixtures.MessageA1,
			fixtures.MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.Has(fixtures.MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM", func() {
		set := TypesOf(
			fixtures.MessageA1,
			fixtures.MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.HasM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.HasM(fixtures.MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add", func() {
		It("adds the type to the set", func() {
			s := TypesOf()
			s.Add(fixtures.MessageAType)

			Expect(
				s.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.Add(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := TypesOf()
			s.Add(fixtures.MessageAType)

			Expect(
				s.Add(fixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func AddM", func() {
		It("adds the type of the message to the set", func() {
			s := TypesOf()
			s.AddM(fixtures.MessageA1)

			Expect(
				s.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.AddM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := TypesOf()
			s.Add(fixtures.MessageAType)

			Expect(
				s.AddM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove", func() {
		It("removes the type from the set", func() {
			s := TypesOf(fixtures.MessageA1)
			s.Remove(fixtures.MessageAType)

			Expect(
				s.Has(fixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := TypesOf()
			s.Add(fixtures.MessageAType)

			Expect(
				s.Remove(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.Remove(fixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM", func() {
		It("removes the type of the message from the set", func() {
			s := TypesOf(fixtures.MessageA1)
			s.RemoveM(fixtures.MessageA1)

			Expect(
				s.Has(fixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := TypesOf()
			s.Add(fixtures.MessageAType)

			Expect(
				s.RemoveM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.RemoveM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		s := NewTypeSet(
			fixtures.MessageAType,
			fixtures.MessageBType,
		)

		It("calls fn for each type in the container", func() {
			var types []message.Type

			all := s.Each(func(t message.Type) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(fixtures.MessageAType, fixtures.MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Each(func(t message.Type) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
