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

var _ HandlerConfig = &ProcessConfig{}

var _ = Describe("type ProcessConfig", func() {
	Describe("func NewProcessConfig", func() {
		var handler *fixtures.ProcessMessageHandler

		BeforeEach(func() {
			handler = &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
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
						message.RoleMap{
							fixtures.MessageAType: message.EventRole,
							fixtures.MessageBType: message.EventRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						message.RoleMap{
							fixtures.MessageCType: message.CommandRole,
						},
					))
				})
			})
		})

		DescribeTable(
			"when the configuration is invalid",
			func(
				msg string,
				fn func(dogma.ProcessConfigurer),
			) {
				handler.ConfigureFunc = fn

				_, err := NewProcessConfig(handler)
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
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.Name()`,
				func(c dogma.ProcessConfigurer) {
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures multiple names",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.Name("<name>")`,
				func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Name("\t \n") with an invalid name`,
				func(c dogma.ProcessConfigurer) {
					c.Name("\t \n")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler does not configure any consumed event types",
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ConsumesEventType()`,
				func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures the same consumed event type multiple times",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ConsumesEventType(fixtures.MessageA)`,
				func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler does not configure any produced commands",
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ProducesCommandType()`,
				func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures the same produced command type multiple times",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ProducesCommandType(fixtures.MessageC)`,
				func(c dogma.ProcessConfigurer) {
					c.Name("<name>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
		)
	})
})
