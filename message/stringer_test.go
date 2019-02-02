package message_test

import (
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type messageStringerMessage struct{}

func (messageStringerMessage) MessageString() string {
	return "<string>"
}

func (messageStringerMessage) String() string {
	panic("unexpected call")
}

type fmtStringerMessage struct{}

func (fmtStringerMessage) String() string {
	return "<string>"
}

type nonStringerMessage struct {
	Value int
}

var _ = Describe("func ToString", func() {
	It("returns the result of MessageString() if the message implements message.Stringer", func() {
		Expect(
			ToString(messageStringerMessage{}),
		).To(Equal(
			"<string>",
		))
	})

	It("returns the result of String() if the message implements fmt.Stringer", func() {
		Expect(
			ToString(fmtStringerMessage{}),
		).To(Equal(
			"<string>",
		))
	})

	It("returns the standard Go representation if the message implements does not implement a Stringer interface", func() {
		Expect(
			ToString(nonStringerMessage{100}),
		).To(Equal(
			"{100}",
		))
	})
})
