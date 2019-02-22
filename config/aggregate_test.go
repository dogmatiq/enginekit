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

var _ HandlerConfig = &AggregateConfig{}

var _ = Describe("type AggregateConfig", func() {
	Describe("func NewAggregateConfig", func() {
		var handler *fixtures.AggregateMessageHandler

		BeforeEach(func() {
			handler = &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.AcceptsCommandType(fixtures.MessageB{})
					c.RecordsEventType(fixtures.MessageE{})
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

			It("the handler name is set", func() {
				Expect(cfg.HandlerName).To(Equal("<name>"))
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
				_, err := NewAggregateConfig(handler)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("the handler does not configure a name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.Name()",
					),
				))
			})
		})

		When("the handler configures multiple names", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.Name("<name>")`,
					),
				))
			})
		})

		When("the handler configures an invalid name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("\t \n")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler.Configure() called AggregateConfigurer.Name("\t \n") with an invalid name`,
					),
				))
			})
		})

		When("the handler does not configure any accepted command types", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.AcceptsCommandType()",
					),
				))
			})
		})

		When("the handler configures the same accepted command type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.AcceptsCommandType(fixtures.MessageA)",
					),
				))
			})
		})

		When("the handler does not configure any recorded events", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.AggregateMessageHandler.Configure() did not call AggregateConfigurer.RecordsEventType()",
					),
				))
			})
		})

		When("the handler configures the same recorded event type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Name("<name>")
					c.AcceptsCommandType(fixtures.MessageA{})
					c.RecordsEventType(fixtures.MessageE{})
					c.RecordsEventType(fixtures.MessageE{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewAggregateConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.AggregateMessageHandler.Configure() has already called AggregateConfigurer.RecordsEventType(fixtures.MessageE)",
					),
				))
			})
		})
	})
})
