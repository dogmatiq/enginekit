package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ Config = &ApplicationConfig{}

var _ = Describe("type ApplicationConfig", func() {
	Describe("func NewApplicationConfig", func() {
		var (
			aggregate   *fixtures.AggregateMessageHandler
			process     *fixtures.ProcessMessageHandler
			integration *fixtures.IntegrationMessageHandler
			projection  *fixtures.ProjectionMessageHandler
			app         *fixtures.Application
		)

		BeforeEach(func() {
			aggregate = &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", "<aggregate-key>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}

			process = &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ConsumesEventType(fixtures.MessageE{}) // shared with <projection>
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			integration = &fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", "<integration-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				},
			}

			projection = &fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection>", "<projection-key>")
					c.ConsumesEventType(fixtures.MessageD{})
					c.ConsumesEventType(fixtures.MessageE{}) // shared with <process>
				},
			}

			app = &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterAggregate(aggregate)
					c.RegisterProcess(process)
					c.RegisterIntegration(integration)
					c.RegisterProjection(projection)
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *ApplicationConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewApplicationConfig(app)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("sets the roles map", func() {
				Expect(cfg.Roles).To(Equal(
					message.RoleMap{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
						fixtures.MessageCType: message.CommandRole,
						fixtures.MessageDType: message.EventRole,
						fixtures.MessageEType: message.EventRole,
						fixtures.MessageFType: message.EventRole,
						fixtures.MessageTType: message.TimeoutRole,
					},
				))
			})

			It("sets the consumers map", func() {
				Expect(cfg.Consumers).To(Equal(
					map[message.Type][]HandlerConfig{
						fixtures.MessageAType: {cfg.HandlersByName["<aggregate>"]},
						fixtures.MessageBType: {cfg.HandlersByName["<process>"]},
						fixtures.MessageCType: {cfg.HandlersByName["<integration>"]},
						fixtures.MessageDType: {cfg.HandlersByName["<projection>"]},
						fixtures.MessageEType: {cfg.HandlersByName["<process>"], cfg.HandlersByName["<projection>"]},
						fixtures.MessageTType: {cfg.HandlersByName["<process>"]},
					},
				))
			})

			It("sets the producers map", func() {
				Expect(cfg.Producers).To(Equal(
					map[message.Type][]HandlerConfig{
						fixtures.MessageCType: {cfg.HandlersByName["<process>"]},
						fixtures.MessageEType: {cfg.HandlersByName["<aggregate>"]},
						fixtures.MessageFType: {cfg.HandlersByName["<integration>"]},
						fixtures.MessageTType: {cfg.HandlersByName["<process>"]},
					},
				))
			})

			Describe("func Name()", func() {
				It("returns the app name", func() {
					Expect(cfg.Name()).To(Equal("<app>"))
				})
			})

			Describe("func Key()", func() {
				It("returns the app key", func() {
					Expect(cfg.Key()).To(Equal("<app-key>"))
				})
			})
		})

		When("the app does not configure an identity", func() {
			BeforeEach(func() {
				app.ConfigureFunc = nil
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						"*fixtures.Application.Configure() did not call ApplicationConfigurer.Identity()",
					),
				))
			})
		})

		When("the app configures multiple identities", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() has already called ApplicationConfigurer.Identity("<name>", "<key>")`,
					),
				))
			})
		})

		When("the app configures an invalid application name", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("\t \n", "<app-key>")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid name "\t \n"`,
					),
				))
			})
		})

		When("the app configures an invalid application key", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "\t \n")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid key "\t \n"`,
					),
				))
			})
		})

		When("the app contains an invalid handler configurations", func() {
			It("returns an error when an aggregate is misconfigured", func() {
				aggregate.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when a process is misconfigured", func() {
				process.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when an integration is misconfigured", func() {
				integration.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when a projection is misconfigured", func() {
				projection.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})
		})

		When("the app contains conflicting handler identities", func() {
			It("returns an error when an aggregate name is in conflict", func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<process>", "<aggregate-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when an aggregate key is in conflict", func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", "<process-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when a process name is in conflict", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<aggregate>", "<process-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler can not use the handler name "<aggregate>", because it is already used by *fixtures.AggregateMessageHandler`,
					),
				))
			})

			It("returns an error when a process key is in conflict", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process>", "<aggregate-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler can not use the handler key "<aggregate-key>", because it is already used by *fixtures.AggregateMessageHandler`,
					),
				))
			})

			It("returns an error when an integration name is in conflict", func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<process>", "<integration-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when an integration key is in conflict", func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", "<process-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when a projection name is in conflict", func() {
				projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					c.Identity("<integration>", "<projection-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageD{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProjectionMessageHandler can not use the handler name "<integration>", because it is already used by *fixtures.IntegrationMessageHandler`,
					),
				))
			})

			It("returns an error when a projection key is in conflict", func() {
				projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection>", "<integration-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageD{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProjectionMessageHandler can not use the handler key "<integration-key>", because it is already used by *fixtures.IntegrationMessageHandler`,
					),
				))
			})
		})

		It("returns an error when the app contains multiple consumers of the same command", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageA{}) // conflict with <aggregate>
				c.ProducesEventType(fixtures.MessageF{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				Error(
					`the "<integration>" handler can not consume fixtures.MessageA commands because they are already consumed by "<aggregate>"`,
				),
			))
		})

		It("returns an error when the app contains multiple producers of the same event", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageE{}) // conflict with <aggregate>
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				Error(
					`the "<integration>" handler can not produce fixtures.MessageE events because they are already produced by "<aggregate>"`,
				),
			))
		})

		It("does not return an error when the app contains multiple processes that schedule the same timeout", func() {
			process1 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-1>", "<process-1-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			process2 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-2>", "<process-2-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			app := &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProcess(process1)
					c.RegisterProcess(process2)
				},
			}

			_, err := NewApplicationConfig(app)

			Expect(err).ShouldNot(HaveOccurred())
		})

		When("multiple handlers use a single message type in differing roles", func() {
			It("returns an error when a conflict occurs with a consumed message", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageA{}) // conflict with <aggregate>
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`the "<process>" handler configures fixtures.MessageA as an event but "<aggregate>" configures it as a command`,
					),
				))
			})

			It("returns an error when a conflict occurs with a produced message", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageE{}) // conflict with <aggregate>
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`the "<process>" handler configures fixtures.MessageE as a command but "<aggregate>" configures it as an event`,
					),
				))
			})
		})
	})
})
