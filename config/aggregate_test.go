package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	handlerkit "github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ HandlerConfig = &AggregateConfig{}

var _ = Describe("type AggregateConfig", func() {
	Describe("func NewAggregateConfig", func() {
		var handler *fixtures.AggregateMessageHandler

		BeforeEach(func() {
			handler = &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
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

			Describe("func Name()", func() {
				It("returns the handler name", func() {
					Expect(cfg.Name()).To(Equal("<name>"))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns handler.AggregateType", func() {
					Expect(cfg.HandlerType()).To(Equal(handlerkit.AggregateType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						message.RoleMap{
							fixtures.MessageAType: message.CommandRole,
							fixtures.MessageBType: message.CommandRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						message.RoleMap{
							fixtures.MessageEType: message.EventRole,
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
				"when the handler does not configure a name",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.Name()`,
				func(c dogma.AggregateConfigurer) {
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures multiple names",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.Name("<name>")`,
				func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Name("\t \n") with an invalid name`,
				func(c dogma.AggregateConfigurer) {
					c.Name("\t \n")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any consumed command types",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ConsumesCommandType()`,
				func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures the same consumed command type multiple times",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ConsumesCommandType(fixtures.MessageA)`,
				func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any produced events",
				`*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.ProducesEventType()`,
				func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures the same produced event type multiple times",
				`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.ProducesEventType(fixtures.MessageE)`,
				func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
		)
	})
})
