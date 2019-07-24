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

var _ HandlerConfig = &ProjectionConfig{}

var _ = Describe("type ProjectionConfig", func() {
	Describe("func NewProjectionConfig", func() {
		var handler *fixtures.ProjectionMessageHandler

		BeforeEach(func() {
			handler = &fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageB{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *ProjectionConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewProjectionConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Describe("func Name()", func() {
				It("returns the handler name", func() {
					Expect(cfg.Name()).To(Equal("<name>"))
				})
			})

			Describe("func Key()", func() {
				It("returns the handler key", func() {
					Expect(cfg.Key()).To(Equal("<key>"))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns handler.ProjectionType", func() {
					Expect(cfg.HandlerType()).To(Equal(handlerkit.ProjectionType))
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
				It("returns an empty map", func() {
					Expect(cfg.ProducedMessageTypes()).To(BeEmpty())
				})
			})
		})

		DescribeTable(
			"when the configuration is invalid",
			func(
				msg string,
				fn func(dogma.ProjectionConfigurer),
			) {
				handler.ConfigureFunc = fn

				_, err := NewProjectionConfig(handler)
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
				`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.Identity()`,
				func(c dogma.ProjectionConfigurer) {
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures multiple identities",
				`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.Identity("<name>", "<key>")`,
				func(c dogma.ProjectionConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.ProjectionMessageHandler.Configure() called ProjectionConfigurer.Identity() with an invalid name "\t \n"`,
				func(c dogma.ProjectionConfigurer) {
					c.Identity("\t \n", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures an invalid key",
				`*fixtures.ProjectionMessageHandler.Configure() called ProjectionConfigurer.Identity() with an invalid key "\t \n"`,
				func(c dogma.ProjectionConfigurer) {
					c.Identity("<name>", "\t \n")
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler does not configure any consumed event types",
				`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.ConsumesEventType()`,
				func(c dogma.ProjectionConfigurer) {
					c.Identity("<name>", "<key>")
				},
			),
			Entry(
				"when the handler configures the same consumed event type multiple times",
				`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.ConsumesEventType(fixtures.MessageA)`,
				func(c dogma.ProjectionConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesEventType(fixtures.MessageA{})
					c.ConsumesEventType(fixtures.MessageA{})
				},
			),
		)
	})
})
