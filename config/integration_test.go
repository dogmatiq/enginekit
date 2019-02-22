package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	handlerkit "github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
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
					c.AcceptsCommandType(fixtures.MessageA{})
					c.AcceptsCommandType(fixtures.MessageB{})
					c.RecordsEventType(fixtures.MessageE{})
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

			It("the handler name is set", func() {
				Expect(cfg.HandlerName).To(Equal("<name>"))
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
						map[message.Type]message.Role{
							fixtures.MessageAType: message.CommandRole,
							fixtures.MessageBType: message.CommandRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						map[message.Type]message.Role{
							fixtures.MessageEType: message.EventRole,
						},
					))
				})
			})
		})

		When("the handler does not configure anything", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = nil
			})

			It("returns an error", func() {
				_, err := NewIntegrationConfig(handler)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("the handler does not configure a name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.Name()",
					),
				))
			})
		})

		When("the handler configures multiple names", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.Name("<name>")`,
					),
				))
			})
		})

		When("the handler configures an invalid name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("\t \n")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler.Configure() called IntegrationConfigurer.Name("\t \n") with an invalid name`,
					),
				))
			})
		})

		When("the handler does not configure any accepted command types", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.AcceptsCommandType()",
					),
				))
			})
		})

		When("the handler configures the same accepted command type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.AcceptsCommandType(fixtures.MessageA)",
					),
				))
			})
		})

		When("the handler does not configure any recorded events", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
				}
			})

			It("does not return an error", func() {
				_, err := NewIntegrationConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("the handler configures the same recorded event type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewIntegrationConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.RecordsEventType(fixtures.MessageE)",
					),
				))
			})
		})
	})
})
