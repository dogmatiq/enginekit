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
			FormatCorrelation(c, nil),
		).To(Equal(
			"= <grandchild>  ∵ <child>  ⋲ <parent>",
		))
	})

	It("returns the expected output for an empty correlation", func() {
		Expect(
			FormatCorrelation(message.Correlation{}, nil),
		).To(Equal(
			"= -  ∵ -  ⋲ -",
		))
	})

	It("uses the format function", func() {
		Expect(
			FormatCorrelation(c, func(id string) string {
				if len(id) > 8 {
					return id[:8]
				}

				return id
			}),
		).To(Equal(
			"= <grandch  ∵ <child>  ⋲ <parent>",
		))
	})
})
