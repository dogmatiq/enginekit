package logging_test

import (
	"errors"

	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/dogmatiq/enginekit/logging"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FormatMessage", func() {
	md := message.MetaData{
		Correlation: message.Correlation{
			MessageID:     "<grandchild>",
			CausationID:   "<child>",
			CorrelationID: "<parent>",
		},
		Type:      fixtures.MessageAType,
		Role:      message.CommandRole,
		Direction: message.InboundDirection,
	}

	When("there is no error", func() {
		It("returns the expected output", func() {
			Expect(
				FormatMessage(
					md,
					8,
					false,
					nil,
					[]string{"<foo>", "<bar>"},
				),
			).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼    fixtures.MessageA? ● <foo> ● <bar>",
			))
		})
	})

	When("the message is being retried", func() {
		It("returns the expected output", func() {
			Expect(
				FormatMessage(
					md,
					8,
					true,
					nil,
					[]string{"<foo>", "<bar>"},
				),
			).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  fixtures.MessageA? ● <foo> ● <bar>",
			))
		})
	})

	When("there is an error", func() {
		It("returns the expected output", func() {
			Expect(
				FormatMessage(
					md,
					8,
					true, // setting isRetry to true should still render the error icon, not the retry icon
					errors.New("<error>"),
					[]string{"<foo>", "<bar>"},
				),
			).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▽ ✖  fixtures.MessageA? ● <error> ● <foo> ● <bar>",
			))
		})
	})
})
