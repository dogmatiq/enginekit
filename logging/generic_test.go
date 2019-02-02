package logging_test

import (
	. "github.com/dogmatiq/enginekit/logging"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Format", func() {
	c := message.Correlation{
		MessageID:     "<grandchild>",
		CausationID:   "<child>",
		CorrelationID: "<parent>",
	}

	It("returns the expected output", func() {
		Expect(
			Format(
				c,
				8,
				[]string{InboundIcon, RetryIcon},
				[]string{"<foo>", "<bar>"},
			),
		).To(Equal(
			"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
		))
	})

	It("includes padding for empty icons", func() {
		Expect(
			Format(
				c,
				8,
				[]string{InboundIcon, ""},
				[]string{"<foo>", "<bar>"},
			),
		).To(Equal(
			"= <grandch  ∵  <child>  ⋲ <parent>  ▼    <foo> ● <bar>",
		))
	})

	It("skips empty text", func() {
		Expect(
			Format(
				c,
				8,
				[]string{InboundIcon, RetryIcon},
				[]string{"<foo>", "", "<bar>"},
			),
		).To(Equal(
			"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
		))
	})

	It("skips empty text in the first position", func() {
		Expect(
			Format(
				c,
				8,
				[]string{InboundIcon, RetryIcon},
				[]string{"", "<foo>", "<bar>"},
			),
		).To(Equal(
			"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
		))
	})
})
