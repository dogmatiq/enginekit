package message_test

import (
	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ TypeContainer = RoleMap{}

var _ = Describe("type RoleMap", func() {
	Describe("func Has()", func() {
		rm := RoleMap{
			fixtures.MessageAType: CommandRole,
			fixtures.MessageBType: EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.Has(fixtures.MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		rm := RoleMap{
			fixtures.MessageAType: CommandRole,
			fixtures.MessageBType: EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.HasM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.HasM(fixtures.MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the type to the map", func() {
			rm := RoleMap{}
			rm.Add(fixtures.MessageAType, CommandRole)

			Expect(
				rm.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := RoleMap{}

			Expect(
				rm.Add(fixtures.MessageAType, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := RoleMap{}
			rm.Add(fixtures.MessageAType, CommandRole)

			Expect(
				rm.Add(fixtures.MessageAType, EventRole),
			).To(BeFalse())

			Expect(
				rm[fixtures.MessageAType],
			).To(Equal(CommandRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the map", func() {
			rm := RoleMap{}
			rm.AddM(fixtures.MessageA1, CommandRole)

			Expect(
				rm.Has(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := RoleMap{}

			Expect(
				rm.AddM(fixtures.MessageA1, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := RoleMap{}
			rm.AddM(fixtures.MessageA1, CommandRole)

			Expect(
				rm.AddM(fixtures.MessageA1, EventRole),
			).To(BeFalse())

			Expect(
				rm[fixtures.MessageAType],
			).To(Equal(CommandRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			rm := RoleMap{fixtures.MessageAType: CommandRole}
			rm.Remove(fixtures.MessageAType)

			Expect(
				rm.Has(fixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := RoleMap{fixtures.MessageAType: CommandRole}

			Expect(
				rm.Remove(fixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := RoleMap{}

			Expect(
				rm.Remove(fixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			rm := RoleMap{fixtures.MessageAType: CommandRole}
			rm.RemoveM(fixtures.MessageA1)

			Expect(
				rm.Has(fixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := RoleMap{fixtures.MessageAType: CommandRole}

			Expect(
				rm.RemoveM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := RoleMap{}

			Expect(
				rm.RemoveM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})
})
