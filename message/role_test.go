package message_test

import (
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Role", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the role is valid", func() {
			CommandRole.MustValidate()
			EventRole.MustValidate()
			TimeoutRole.MustValidate()
		})

		It("panics when the role is not valid", func() {
			Expect(func() {
				Role(-1).MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is", func() {
		It("returns true when the role is in the given set", func() {
			Expect(CommandRole.Is(CommandRole, EventRole)).To(BeTrue())
		})

		It("returns false when the role is not in the given set", func() {
			Expect(TimeoutRole.Is(CommandRole, EventRole)).To(BeFalse())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the role is in the given set", func() {
			CommandRole.MustBe(CommandRole, EventRole)
		})

		It("panics when the role is not in the given set", func() {
			Expect(func() {
				TimeoutRole.MustBe(CommandRole, EventRole)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the role is not in the given set", func() {
			TimeoutRole.MustNotBe(CommandRole, EventRole)
		})

		It("panics when the role is in the given set", func() {
			Expect(func() {
				CommandRole.MustNotBe(CommandRole, EventRole)
			}).To(Panic())
		})
	})

	Describe("func String", func() {
		It("returns the role value as a string", func() {
			Expect(CommandRole.String()).To(Equal("command"))
			Expect(EventRole.String()).To(Equal("event"))
			Expect(TimeoutRole.String()).To(Equal("timeout"))
		})
	})
})
