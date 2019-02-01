package message_test

import (
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Durection", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the direction is valid", func() {
			InboundDirection.MustValidate()
			OutboundDirection.MustValidate()
		})

		It("panics when the direction is not valid", func() {
			Expect(func() {
				Direction("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the direction is equal to the given direction", func() {
			InboundDirection.MustBe(InboundDirection)
		})

		It("panics when the direction is not equal to the given direction", func() {
			Expect(func() {
				InboundDirection.MustBe(OutboundDirection)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the direction is not equal to the given direction", func() {
			InboundDirection.MustNotBe(OutboundDirection)
		})

		It("panics when the direction is equal to the given direction", func() {
			Expect(func() {
				InboundDirection.MustNotBe(InboundDirection)
			}).To(Panic())
		})
	})

	Describe("func String", func() {
		It("returns the direction value as a string", func() {
			Expect(InboundDirection.String()).To(Equal("inbound"))
			Expect(OutboundDirection.String()).To(Equal("outbound"))
		})
	})
})
