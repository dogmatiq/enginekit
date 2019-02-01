package message_test

import (
	. "github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Correlation", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the correlation is valid", func() {
			NewCorrelation("<id>").MustValidate()
		})

		It("panics when the message ID is empty", func() {
			Expect(func() {
				Correlation{
					CausationID:   "<cause>",
					CorrelationID: "<corr>",
				}.MustValidate()
			}).To(Panic())
		})

		It("panics when the causation ID is empty", func() {
			Expect(func() {
				Correlation{
					MessageID:     "<id>",
					CorrelationID: "<corr>",
				}.MustValidate()
			}).To(Panic())
		})

		It("panics when the correlation ID is empty", func() {
			Expect(func() {
				Correlation{
					MessageID:   "<id>",
					CausationID: "<cause>",
				}.MustValidate()
			}).To(Panic())
		})
	})

	Describe("func New", func() {
		It("returns a correlation with a causal relationship", func() {
			Expect(
				NewCorrelation("<parent>").
					New("<child>"),
			).To(Equal(Correlation{
				MessageID:     "<child>",
				CausationID:   "<parent>",
				CorrelationID: "<parent>",
			}))
		})

		It("maintains the correlation ID across multiple generations", func() {
			Expect(
				NewCorrelation("<parent>").
					New("<child>").
					New("<grandchild>"),
			).To(Equal(Correlation{
				MessageID:     "<grandchild>",
				CausationID:   "<child>",
				CorrelationID: "<parent>",
			}))
		})

		It("panics if the given ID is the same as the parent's message ID", func() {
			Expect(func() {
				NewCorrelation("<id>").
					New("<id>")
			}).To(Panic())
		})

		It("panics if the given ID is the same as the parent's causation ID", func() {
			Expect(func() {
				NewCorrelation("<parent>").
					New("<child>").
					New("<parent>")
			}).To(Panic())
		})

		It("panics if the given ID is the same as the parent's correlation ID", func() {
			Expect(func() {
				NewCorrelation("<parent>").
					New("<child>").
					New("<grandchild>").
					New("<parent>")
			}).To(Panic())
		})
	})

	Describe("func IsRoot", func() {
		It("returns true for root messages", func() {
			Expect(
				NewCorrelation("<id>").IsRoot(),
			).To(BeTrue())
		})

		It("returns false for non-root messages", func() {
			Expect(
				NewCorrelation("<parent>").New("<child>").IsRoot(),
			).To(BeFalse())
		})
	})

	Describe("func IsCausedBy", func() {
		It("returns true for direct causal messages", func() {
			p := NewCorrelation("<parent>")
			c := p.New("<child>")

			Expect(
				c.IsCausedBy(p),
			).To(BeTrue())
		})

		It("returns false for correlated messages", func() {
			p := NewCorrelation("<parent>")
			c := p.New("<child>")
			g := c.New("<grandchild>")

			Expect(
				g.IsCausedBy(p),
			).To(BeFalse())
		})

		It("returns false for unrelated messages", func() {
			p := NewCorrelation("<parent>")
			u := NewCorrelation("<unrelated>")

			Expect(
				u.IsCausedBy(p),
			).To(BeFalse())
		})
	})

	Describe("func IsCorrelatedWith", func() {
		It("returns true for direct causal messages", func() {
			p := NewCorrelation("<parent>")
			c := p.New("<child>")

			Expect(
				c.IsCorrelatedWith(p),
			).To(BeTrue())
		})

		It("returns true for correlated messages", func() {
			p := NewCorrelation("<parent>")
			c := p.New("<child>")
			g := c.New("<grandchild>")

			Expect(
				g.IsCorrelatedWith(p),
			).To(BeTrue())
		})

		It("returns false for unrelated messages", func() {
			p := NewCorrelation("<parent>")
			u := NewCorrelation("<unrelated>")

			Expect(
				u.IsCorrelatedWith(p),
			).To(BeFalse())
		})
	})
})
