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
					fixtures.MessageA1,
					8,
					SystemIcon,
					nil,
					"<foo>",
					"<bar>",
				),
			).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ⚙  fixtures.MessageA? ● {A1} ● <foo> ● <bar>",
			))
		})
	})

	When("there is an error", func() {
		It("returns the expected output", func() {
			Expect(
				FormatMessage(
					md,
					fixtures.MessageA1,
					8,
					SystemIcon,
					errors.New("<error>"),
					"<foo>",
					"<bar>",
				),
			).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▽ ✖  fixtures.MessageA? ● <error> ● {A1} ● <foo> ● <bar>",
			))
		})
	})
})
