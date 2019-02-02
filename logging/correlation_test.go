package logging_test

import (
	. "github.com/dogmatiq/enginekit/logging"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FormatCorrelation", func() {
	c := message.Correlation{
		MessageID:     "<grandchild>",
		CausationID:   "<child>",
		CorrelationID: "<parent>",
	}

	It("returns the expected output", func() {
		Expect(
			FormatCorrelation(c, 0),
		).To(Equal(
			"= <grandchild>  ∵ <child>  ⋲ <parent>",
		))
	})

	It("truncates message IDs (tail)", func() {
		Expect(
			FormatCorrelation(c, 8),
		).To(Equal(
			"= <grandch  ∵  <child>  ⋲ <parent>",
		))
	})

	It("truncates message IDs (head)", func() {
		Expect(
			FormatCorrelation(c, -8),
		).To(Equal(
			"= ndchild>  ∵  <child>  ⋲ <parent>",
		))
	})
})
