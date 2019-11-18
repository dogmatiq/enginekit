package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ HandlerConfig = &AggregateConfig{}

var _ = Describe("type AggregateConfig", func() {
	Describe("func NewAggregateConfig", func() {
		var handler *AggregateMessageHandler

		BeforeEach(func() {
			handler = &AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ConsumesCommandType(MessageB{})
					c.ProducesEventType(MessageE{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *AggregateConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewAggregateConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Describe("func Identity()", func() {
				It("returns the handler identity", func() {
					Expect(cfg.Identity()).To(Equal(
						MustNewIdentity("<name>", "<key>"),
					))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns AggregateHandlerType", func() {
					Expect(cfg.HandlerType()).To(Equal(AggregateHandlerType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						MessageRoleMap{
							MessageAType: CommandMessageRole,
							MessageBType: CommandMessageRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						MessageRoleMap{
							MessageEType: EventMessageRole,
						},
					))
				})
			})
		})

		DescribeTable(
			"when the configuration is invalid",
			func(
				msg string,
				fn func(dogma.AggregateConfigurer),
			) {
				handler.ConfigureFunc = fn

				_, err := NewAggregateConfig(handler)
				Expect(err).Should(HaveOccurred())

				if msg != "" {
					Expect(err).To(MatchError(msg))
				}
			},
			Entry(
				"when the handler does not configure anything",
				"", // any error
				nil,
			),
			Entry(
				"when the handler does not configure an identity",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.Identity()`,
				func(c dogma.AggregateConfigurer) {
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures multiple identities",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.Identity("<name>", "<key>")`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("\t \n", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid key",
				`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "\t \n")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any consumed command types",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ConsumesCommandType()`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures the same consumed command type multiple times",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ConsumesCommandType(fixtures.MessageA)`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any produced events",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ProducesEventType()`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
				},
			),
			Entry(
				"when the handler configures the same produced event type multiple times",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ProducesEventType(fixtures.MessageE)`,
				func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
					c.ProducesEventType(MessageE{})
				},
			),
		)
	})
})
