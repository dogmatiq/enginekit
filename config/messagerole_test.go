package config_test

import (
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type MessageRole", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the role is valid", func() {
			CommandMessageRole.MustValidate()
			EventMessageRole.MustValidate()
			TimeoutMessageRole.MustValidate()
		})

		It("panics when the role is not valid", func() {
			Expect(func() {
				MessageRole("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is", func() {
		It("returns true when the role is in the given set", func() {
			Expect(CommandMessageRole.Is(CommandMessageRole, EventMessageRole)).To(BeTrue())
		})

		It("returns false when the role is not in the given set", func() {
			Expect(TimeoutMessageRole.Is(CommandMessageRole, EventMessageRole)).To(BeFalse())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the role is in the given set", func() {
			CommandMessageRole.MustBe(CommandMessageRole, EventMessageRole)
		})

		It("panics when the role is not in the given set", func() {
			Expect(func() {
				TimeoutMessageRole.MustBe(CommandMessageRole, EventMessageRole)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the role is not in the given set", func() {
			TimeoutMessageRole.MustNotBe(CommandMessageRole, EventMessageRole)
		})

		It("panics when the role is in the given set", func() {
			Expect(func() {
				CommandMessageRole.MustNotBe(CommandMessageRole, EventMessageRole)
			}).To(Panic())
		})
	})

	Describe("func Marker", func() {
		It("returns the correct marker character", func() {
			Expect(CommandMessageRole.Marker()).To(Equal("?"))
			Expect(EventMessageRole.Marker()).To(Equal("!"))
			Expect(TimeoutMessageRole.Marker()).To(Equal("@"))
		})
	})

	Describe("func String", func() {
		It("returns the role value as a string", func() {
			Expect(CommandMessageRole.String()).To(Equal("command"))
			Expect(EventMessageRole.String()).To(Equal("event"))
			Expect(TimeoutMessageRole.String()).To(Equal("timeout"))
		})
	})
})
