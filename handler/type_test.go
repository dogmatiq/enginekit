package handler_test

import (
	. "github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Type", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the type is valid", func() {
			AggregateType.MustValidate()
			ProcessType.MustValidate()
			IntegrationType.MustValidate()
			ProjectionType.MustValidate()
		})

		It("panics when the type is not valid", func() {
			Expect(func() {
				Type("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is", func() {
		It("returns true when the type is in the given set", func() {
			Expect(AggregateType.Is(AggregateType, ProcessType)).To(BeTrue())
		})

		It("returns false when the type is not in the given set", func() {
			Expect(IntegrationType.Is(AggregateType, ProcessType)).To(BeFalse())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the type is in the given set", func() {
			AggregateType.MustBe(AggregateType, ProcessType)
		})

		It("panics when the type is not in the given set", func() {
			Expect(func() {
				IntegrationType.MustBe(AggregateType, ProcessType)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the type is not in the given set", func() {
			IntegrationType.MustNotBe(AggregateType, ProcessType)
		})

		It("panics when the type is in the given set", func() {
			Expect(func() {
				AggregateType.MustNotBe(AggregateType, ProcessType)
			}).To(Panic())
		})
	})

	Describe("func IsConsumerOf", func() {
		It("returns the expected values", func() {
			Expect(AggregateType.IsConsumerOf(message.CommandRole)).To(BeTrue())
			Expect(AggregateType.IsConsumerOf(message.EventRole)).To(BeFalse())
			Expect(AggregateType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProcessType.IsConsumerOf(message.CommandRole)).To(BeFalse())
			Expect(ProcessType.IsConsumerOf(message.EventRole)).To(BeTrue())
			Expect(ProcessType.IsConsumerOf(message.TimeoutRole)).To(BeTrue())

			Expect(IntegrationType.IsConsumerOf(message.CommandRole)).To(BeTrue())
			Expect(IntegrationType.IsConsumerOf(message.EventRole)).To(BeFalse())
			Expect(IntegrationType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProjectionType.IsConsumerOf(message.CommandRole)).To(BeFalse())
			Expect(ProjectionType.IsConsumerOf(message.EventRole)).To(BeTrue())
			Expect(ProjectionType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())
		})
	})

	Describe("func IsProducerOf", func() {
		It("returns the expected values", func() {
			Expect(AggregateType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(AggregateType.IsProducerOf(message.EventRole)).To(BeTrue())
			Expect(AggregateType.IsProducerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProcessType.IsProducerOf(message.CommandRole)).To(BeTrue())
			Expect(ProcessType.IsProducerOf(message.EventRole)).To(BeFalse())
			Expect(ProcessType.IsProducerOf(message.TimeoutRole)).To(BeTrue())

			Expect(IntegrationType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(IntegrationType.IsProducerOf(message.EventRole)).To(BeTrue())
			Expect(IntegrationType.IsProducerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProjectionType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(ProjectionType.IsProducerOf(message.EventRole)).To(BeFalse())
			Expect(ProjectionType.IsProducerOf(message.TimeoutRole)).To(BeFalse())
		})
	})

	Describe("func Consumes", func() {
		It("returns the expected values", func() {
			Expect(AggregateType.Consumes()).To(ConsistOf(
				message.CommandRole,
			))

			Expect(ProcessType.Consumes()).To(ConsistOf(
				message.EventRole,
				message.TimeoutRole,
			))

			Expect(IntegrationType.Consumes()).To(ConsistOf(
				message.CommandRole,
			))

			Expect(ProjectionType.Consumes()).To(ConsistOf(
				message.EventRole,
			))
		})
	})

	Describe("func Produces", func() {
		It("returns the expected values", func() {
			Expect(AggregateType.Produces()).To(ConsistOf(
				message.EventRole,
			))

			Expect(ProcessType.Produces()).To(ConsistOf(
				message.CommandRole,
				message.TimeoutRole,
			))

			Expect(IntegrationType.Produces()).To(ConsistOf(
				message.EventRole,
			))

			Expect(ProjectionType.Produces()).To(BeEmpty())
		})
	})

	Describe("func ShortString", func() {
		It("returns the type value as a short string", func() {
			Expect(AggregateType.ShortString()).To(Equal("agg"))
			Expect(ProcessType.ShortString()).To(Equal("prc"))
			Expect(IntegrationType.ShortString()).To(Equal("int"))
			Expect(ProjectionType.ShortString()).To(Equal("prj"))
		})
	})

	Describe("func String", func() {
		It("returns the type value as a string", func() {
			Expect(AggregateType.String()).To(Equal("aggregate"))
			Expect(ProcessType.String()).To(Equal("process"))
			Expect(IntegrationType.String()).To(Equal("integration"))
			Expect(ProjectionType.String()).To(Equal("projection"))
		})
	})

	Describe("func ConsumersOf", func() {
		It("returns the expected values", func() {
			Expect(ConsumersOf(message.CommandRole)).To(ConsistOf(
				AggregateType,
				IntegrationType,
			))

			Expect(ConsumersOf(message.EventRole)).To(ConsistOf(
				ProcessType,
				ProjectionType,
			))

			Expect(ConsumersOf(message.TimeoutRole)).To(ConsistOf(
				ProcessType,
			))
		})
	})

	Describe("func ProducersOf", func() {
		It("returns the expected values", func() {
			Expect(ProducersOf(message.CommandRole)).To(ConsistOf(
				ProcessType,
			))

			Expect(ProducersOf(message.EventRole)).To(ConsistOf(
				AggregateType,
				IntegrationType,
			))

			Expect(ProducersOf(message.TimeoutRole)).To(ConsistOf(
				ProcessType,
			))
		})
	})
})
