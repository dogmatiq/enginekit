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

var _ HandlerConfig = &ProcessConfig{}

var _ = Describe("type ProcessConfig", func() {
	Describe("func NewProcessConfig", func() {
		var handler *fixtures.ProcessMessageHandler

		BeforeEach(func() {
			handler = &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.AcceptsEventType(fixtures.MessageB{})
					c.ExecutesCommandType(fixtures.MessageC{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *ProcessConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewProcessConfig(handler)
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
				It("returns handler.ProcessType", func() {
					Expect(cfg.HandlerType()).To(Equal(handlerkit.ProcessType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						map[message.Type]message.Role{
							fixtures.MessageAType: message.EventRole,
							fixtures.MessageBType: message.EventRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						map[message.Type]message.Role{
							fixtures.MessageCType: message.CommandRole,
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
				_, err := NewProcessConfig(handler)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("the handler does not configure a name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.AcceptsEventType(fixtures.MessageA{})
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.Name()",
					),
				))
			})
		})

		When("the handler configures multiple names", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.Name("<name>")`,
					),
				))
			})
		})

		When("the handler configures an invalid name", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("\t \n")
					c.AcceptsEventType(fixtures.MessageA{})
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Name("\t \n") with an invalid name`,
					),
				))
			})
		})

		When("the handler does not configure any accepted event types", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.AcceptsEventType()",
					),
				))
			})
		})

		When("the handler configures the same accepted event type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.AcceptsEventType(fixtures.MessageA{})
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.AcceptsEventType(fixtures.MessageA)",
					),
				))
			})
		})

		When("the handler does not configure any executed command types", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ExecutesCommandType()",
					),
				))
			})
		})

		When("the handler configures the same executed command type multiple times", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.ExecutesCommandType(fixtures.MessageC{})
					c.ExecutesCommandType(fixtures.MessageC{})
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewProcessConfig(handler)

				Expect(err).To(Equal(
					Error(
						"*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ExecutesCommandType(fixtures.MessageC)",
					),
				))
			})
		})
	})
})
