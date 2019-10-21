package message_test

import (
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type describer struct{}

func (describer) MessageDescription() string {
	return "<description>"
}

func (describer) String() string {
	panic("unexpected call")
}

type stringer struct{}

func (stringer) String() string {
	return "<string>"
}

type indescribable struct {
	Value int
}

var _ = Describe("func Description()", func() {
	It("returns the result of MessageString() if the message implements message.Stringer", func() {
		Expect(
			Description(describer{}),
		).To(Equal(
			"<description>",
		))
	})

	It("returns the result of String() if the message implements fmt.Stringer", func() {
		Expect(
			Description(stringer{}),
		).To(Equal(
			"<string>",
		))
	})

	It("returns the standard Go representation if the message implements does not implement a Stringer interface", func() {
		Expect(
			Description(indescribable{100}),
		).To(Equal(
			"{100}",
		))
	})
})
