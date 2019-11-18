package config_test

import (
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ MessageTypeContainer = MessageTypeSet{}

var _ = Describe("type MessageTypeSet", func() {
	Describe("func NewMessageTypeSet", func() {
		It("returns a set containing the given types", func() {
			Expect(NewMessageTypeSet(
				MessageAType,
				MessageBType,
			)).To(Equal(MessageTypeSet{
				MessageAType: struct{}{},
				MessageBType: struct{}{},
			}))
		})
	})

	Describe("func MessageTypesOf", func() {
		It("returns a set containing the types of the given messages", func() {
			Expect(MessageTypesOf(
				MessageA1,
				MessageB1,
			)).To(Equal(MessageTypeSet{
				MessageAType: struct{}{},
				MessageBType: struct{}{},
			}))
		})
	})

	Describe("func Has", func() {
		set := MessageTypesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.Has(MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM", func() {
		set := MessageTypesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add", func() {
		It("adds the type to the set", func() {
			s := MessageTypesOf()
			s.Add(MessageAType)

			Expect(
				s.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.Add(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(MessageAType)

			Expect(
				s.Add(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func AddM", func() {
		It("adds the type of the message to the set", func() {
			s := MessageTypesOf()
			s.AddM(MessageA1)

			Expect(
				s.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.AddM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(MessageAType)

			Expect(
				s.AddM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove", func() {
		It("removes the type from the set", func() {
			s := MessageTypesOf(MessageA1)
			s.Remove(MessageAType)

			Expect(
				s.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(MessageAType)

			Expect(
				s.Remove(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.Remove(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM", func() {
		It("removes the type of the message from the set", func() {
			s := MessageTypesOf(MessageA1)
			s.RemoveM(MessageA1)

			Expect(
				s.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(MessageAType)

			Expect(
				s.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		s := NewMessageTypeSet(
			MessageAType,
			MessageBType,
		)

		It("calls fn for each type in the container", func() {
			var types []MessageType

			all := s.Each(func(t MessageType) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(MessageAType, MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Each(func(t MessageType) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
