package logging_test

import (
	"fmt"

	. "github.com/dogmatiq/enginekit/logging"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Logger", func() {
	c := message.Correlation{
		MessageID:     "<grandchild>",
		CausationID:   "<child>",
		CorrelationID: "<parent>",
	}

	var (
		output string
		logger *Logger
	)

	BeforeEach(func() {
		output = ""

		logger = &Logger{
			Log: func(s string) {
				output = s
			},
			FormatMessageID: func(id string) string {
				if len(id) > 8 {
					return id[:8]
				}

				return fmt.Sprintf("%8s", id)
			},
		}
	})

	Describe("func LogGeneric", func() {
		It("logs a generic message", func() {
			logger.LogGeneric(
				c,
				[]string{InboundIcon, RetryIcon},
				"<foo>",
				"<bar>",
			)

			Expect(output).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
			))
		})

		It("includes padding for empty icons", func() {
			logger.LogGeneric(
				c,
				[]string{InboundIcon, ""},
				"<foo>",
				"<bar>",
			)

			Expect(output).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼    <foo> ● <bar>",
			))
		})

		It("skips empty text", func() {
			logger.LogGeneric(
				c,
				[]string{InboundIcon, RetryIcon},
				"<foo>",
				"",
				"<bar>",
			)

			Expect(output).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
			))
		})

		It("skips empty text in the first position", func() {
			logger.LogGeneric(
				c,
				[]string{InboundIcon, RetryIcon},
				"",
				"<foo>",
				"<bar>",
			)

			Expect(output).To(Equal(
				"= <grandch  ∵  <child>  ⋲ <parent>  ▼ ↻  <foo> ● <bar>",
			))
		})
	})
})
