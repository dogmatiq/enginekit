package config_test

import (
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ MessageTypeContainer = MessageRoleMap{}

var _ = Describe("type MessageRoleMap", func() {
	Describe("func Has()", func() {
		rm := MessageRoleMap{
			MessageAType: CommandMessageRole,
			MessageBType: EventMessageRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.Has(MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		rm := MessageRoleMap{
			MessageAType: CommandMessageRole,
			MessageBType: EventMessageRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the type to the map", func() {
			rm := MessageRoleMap{}
			rm.Add(MessageAType, CommandMessageRole)

			Expect(
				rm.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.Add(MessageAType, CommandMessageRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := MessageRoleMap{}
			rm.Add(MessageAType, CommandMessageRole)

			Expect(
				rm.Add(MessageAType, EventMessageRole),
			).To(BeFalse())

			Expect(
				rm[MessageAType],
			).To(Equal(CommandMessageRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the map", func() {
			rm := MessageRoleMap{}
			rm.AddM(MessageA1, CommandMessageRole)

			Expect(
				rm.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.AddM(MessageA1, CommandMessageRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := MessageRoleMap{}
			rm.AddM(MessageA1, CommandMessageRole)

			Expect(
				rm.AddM(MessageA1, EventMessageRole),
			).To(BeFalse())

			Expect(
				rm[MessageAType],
			).To(Equal(CommandMessageRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			rm := MessageRoleMap{MessageAType: CommandMessageRole}
			rm.Remove(MessageAType)

			Expect(
				rm.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := MessageRoleMap{MessageAType: CommandMessageRole}

			Expect(
				rm.Remove(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.Remove(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			rm := MessageRoleMap{MessageAType: CommandMessageRole}
			rm.RemoveM(MessageA1)

			Expect(
				rm.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := MessageRoleMap{MessageAType: CommandMessageRole}

			Expect(
				rm.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		rm := MessageRoleMap{
			MessageAType: CommandMessageRole,
			MessageBType: EventMessageRole,
		}

		It("calls fn for each type in the container", func() {
			var types []MessageType

			all := rm.Each(func(t MessageType) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(MessageAType, MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := rm.Each(func(t MessageType) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
