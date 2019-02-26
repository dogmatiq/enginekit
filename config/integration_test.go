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

var _ HandlerConfig = &IntegrationConfig{}

var _ = Describe("type IntegrationConfig", func() {
	Describe("func NewIntegrationConfig", func() {
		var handler *fixtures.IntegrationMessageHandler

		BeforeEach(func() {
			handler = &fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *IntegrationConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewIntegrationConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Describe("func Name()", func() {
				It("returns the handler name", func() {
					Expect(cfg.Name()).To(Equal("<name>"))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns handler.IntegrationType", func() {
					Expect(cfg.HandlerType()).To(Equal(handlerkit.IntegrationType))
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

			When("the handler does not configure any produced events", func() {
				BeforeEach(func() {
					handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
						c.Name("<name>")
						c.ConsumesCommandType(fixtures.MessageA{})
					}
				})

				It("does not return an error", func() {
					_, err := NewIntegrationConfig(handler)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		DescribeTable(
			"when the configuration is invalid",
			func(
				msg string,
				fn func(dogma.IntegrationConfigurer),
			) {
				handler.ConfigureFunc = fn

				_, err := NewIntegrationConfig(handler)
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
				`*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.Name()`,
				func(c dogma.IntegrationConfigurer) {
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures multiple names",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.Name("<name>")`,
				func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.IntegrationMessageHandler.Configure() called IntegrationConfigurer.Name("\t \n") with an invalid name`,
				func(c dogma.IntegrationConfigurer) {
					c.Name("\t \n")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any consumed command types",
				`*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.ConsumesCommandType()`,
				func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures the same consumed command type multiple times",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.ConsumesCommandType(fixtures.MessageA)`,
				func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
			Entry(
				"when the handler configures the same produced event type multiple times",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.ProducesEventType(fixtures.MessageE)`,
				func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			),
		)
	})
})
