package config_test

import (
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type HandlerType", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the type is valid", func() {
			AggregateHandlerType.MustValidate()
			ProcessHandlerType.MustValidate()
			IntegrationHandlerType.MustValidate()
			ProjectionHandlerType.MustValidate()
		})

		It("panics when the type is not valid", func() {
			Expect(func() {
				HandlerType("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is", func() {
		It("returns true when the type is in the given set", func() {
			Expect(AggregateHandlerType.Is(AggregateHandlerType, ProcessHandlerType)).To(BeTrue())
		})

		It("returns false when the type is not in the given set", func() {
			Expect(IntegrationHandlerType.Is(AggregateHandlerType, ProcessHandlerType)).To(BeFalse())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the type is in the given set", func() {
			AggregateHandlerType.MustBe(AggregateHandlerType, ProcessHandlerType)
		})

		It("panics when the type is not in the given set", func() {
			Expect(func() {
				IntegrationHandlerType.MustBe(AggregateHandlerType, ProcessHandlerType)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the type is not in the given set", func() {
			IntegrationHandlerType.MustNotBe(AggregateHandlerType, ProcessHandlerType)
		})

		It("panics when the type is in the given set", func() {
			Expect(func() {
				AggregateHandlerType.MustNotBe(AggregateHandlerType, ProcessHandlerType)
			}).To(Panic())
		})
	})

	Describe("func IsConsumerOf", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.IsConsumerOf(CommandMessageRole)).To(BeTrue())
			Expect(AggregateHandlerType.IsConsumerOf(EventMessageRole)).To(BeFalse())
			Expect(AggregateHandlerType.IsConsumerOf(TimeoutMessageRole)).To(BeFalse())

			Expect(ProcessHandlerType.IsConsumerOf(CommandMessageRole)).To(BeFalse())
			Expect(ProcessHandlerType.IsConsumerOf(EventMessageRole)).To(BeTrue())
			Expect(ProcessHandlerType.IsConsumerOf(TimeoutMessageRole)).To(BeTrue())

			Expect(IntegrationHandlerType.IsConsumerOf(CommandMessageRole)).To(BeTrue())
			Expect(IntegrationHandlerType.IsConsumerOf(EventMessageRole)).To(BeFalse())
			Expect(IntegrationHandlerType.IsConsumerOf(TimeoutMessageRole)).To(BeFalse())

			Expect(ProjectionHandlerType.IsConsumerOf(CommandMessageRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsConsumerOf(EventMessageRole)).To(BeTrue())
			Expect(ProjectionHandlerType.IsConsumerOf(TimeoutMessageRole)).To(BeFalse())
		})
	})

	Describe("func IsProducerOf", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.IsProducerOf(CommandMessageRole)).To(BeFalse())
			Expect(AggregateHandlerType.IsProducerOf(EventMessageRole)).To(BeTrue())
			Expect(AggregateHandlerType.IsProducerOf(TimeoutMessageRole)).To(BeFalse())

			Expect(ProcessHandlerType.IsProducerOf(CommandMessageRole)).To(BeTrue())
			Expect(ProcessHandlerType.IsProducerOf(EventMessageRole)).To(BeFalse())
			Expect(ProcessHandlerType.IsProducerOf(TimeoutMessageRole)).To(BeTrue())

			Expect(IntegrationHandlerType.IsProducerOf(CommandMessageRole)).To(BeFalse())
			Expect(IntegrationHandlerType.IsProducerOf(EventMessageRole)).To(BeTrue())
			Expect(IntegrationHandlerType.IsProducerOf(TimeoutMessageRole)).To(BeFalse())

			Expect(ProjectionHandlerType.IsProducerOf(CommandMessageRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(EventMessageRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(TimeoutMessageRole)).To(BeFalse())
		})
	})

	Describe("func Consumes", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Consumes()).To(ConsistOf(
				CommandMessageRole,
			))

			Expect(ProcessHandlerType.Consumes()).To(ConsistOf(
				EventMessageRole,
				TimeoutMessageRole,
			))

			Expect(IntegrationHandlerType.Consumes()).To(ConsistOf(
				CommandMessageRole,
			))

			Expect(ProjectionHandlerType.Consumes()).To(ConsistOf(
				EventMessageRole,
			))
		})
	})

	Describe("func Produces", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Produces()).To(ConsistOf(
				EventMessageRole,
			))

			Expect(ProcessHandlerType.Produces()).To(ConsistOf(
				CommandMessageRole,
				TimeoutMessageRole,
			))

			Expect(IntegrationHandlerType.Produces()).To(ConsistOf(
				EventMessageRole,
			))

			Expect(ProjectionHandlerType.Produces()).To(BeEmpty())
		})
	})

	Describe("func ShortString", func() {
		It("returns the type value as a short string", func() {
			Expect(AggregateHandlerType.ShortString()).To(Equal("agg"))
			Expect(ProcessHandlerType.ShortString()).To(Equal("prc"))
			Expect(IntegrationHandlerType.ShortString()).To(Equal("int"))
			Expect(ProjectionHandlerType.ShortString()).To(Equal("prj"))
		})
	})

	Describe("func String", func() {
		It("returns the type value as a string", func() {
			Expect(AggregateHandlerType.String()).To(Equal("aggregate"))
			Expect(ProcessHandlerType.String()).To(Equal("process"))
			Expect(IntegrationHandlerType.String()).To(Equal("integration"))
			Expect(ProjectionHandlerType.String()).To(Equal("projection"))
		})
	})

	Describe("func ConsumersOf", func() {
		It("returns the expected values", func() {
			Expect(ConsumersOf(CommandMessageRole)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ConsumersOf(EventMessageRole)).To(ConsistOf(
				ProcessHandlerType,
				ProjectionHandlerType,
			))

			Expect(ConsumersOf(TimeoutMessageRole)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})

	Describe("func ProducersOf", func() {
		It("returns the expected values", func() {
			Expect(ProducersOf(CommandMessageRole)).To(ConsistOf(
				ProcessHandlerType,
			))

			Expect(ProducersOf(EventMessageRole)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ProducersOf(TimeoutMessageRole)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})
})
