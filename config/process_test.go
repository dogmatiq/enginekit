package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
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
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
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

			Describe("func Identity()", func() {
				It("returns the handler identity", func() {
					Expect(cfg.Identity()).To(Equal(
						MustNewIdentity("<name>", "<key>"),
					))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns ProcessHandlerType", func() {
					Expect(cfg.HandlerType()).To(Equal(ProcessHandlerType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						message.RoleMap{
							fixtures.MessageAType: message.EventRole,
							fixtures.MessageBType: message.EventRole,
							fixtures.MessageTType: message.TimeoutRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						message.RoleMap{
							fixtures.MessageCType: message.CommandRole,
							fixtures.MessageTType: message.TimeoutRole,
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
				"when the handler does not configure an identity",
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.Identity()`,
				func(c dogma.ProcessConfigurer) {
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures multiple identities",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.Identity("<name>", "<key>")`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("\t \n", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures an invalid key",
				`*fixtures.ProcessMessageHandler.Configure() called ProcessConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "\t \n")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler does not configure any consumed event types",
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ConsumesEventType()`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures the same consumed event type multiple times",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ConsumesEventType(fixtures.MessageA)`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures an event that was previously configured as a timeout",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.SchedulesTimeoutType(fixtures.MessageA)`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.SchedulesTimeoutType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler does not configure any produced commands",
				`*fixtures.ProcessMessageHandler.Configure() did not call ProcessConfigurer.ProducesCommandType()`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures the same produced command type multiple times",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.ProducesCommandType(fixtures.MessageC)`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
			Entry(
				"when the handler configures a command that was previously configured as a timeout",
				`*fixtures.ProcessMessageHandler.Configure() has already called ProcessConfigurer.SchedulesTimeoutType(fixtures.MessageC)`,
				func(c dogma.ProcessConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.SchedulesTimeoutType(fixtures.MessageC{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			),
		)
	})
})
