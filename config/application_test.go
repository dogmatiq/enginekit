package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ Config = &ApplicationConfig{}

var _ = Describe("type ApplicationConfig", func() {
	Describe("func NewApplicationConfig", func() {
		var (
			aggregate   *AggregateMessageHandler
			process     *ProcessMessageHandler
			integration *IntegrationMessageHandler
			projection  *ProjectionMessageHandler
			app         *Application
		)

		BeforeEach(func() {
			aggregate = &AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", "<aggregate-key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			}

			process = &ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process>", "<process-key>")
					c.ConsumesEventType(MessageB{})
					c.ConsumesEventType(MessageE{}) // shared with <projection>
					c.ProducesCommandType(MessageC{})
					c.SchedulesTimeoutType(MessageT{})
				},
			}

			integration = &IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", "<integration-key>")
					c.ConsumesCommandType(MessageC{})
					c.ProducesEventType(MessageF{})
				},
			}

			projection = &ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection>", "<projection-key>")
					c.ConsumesEventType(MessageD{})
					c.ConsumesEventType(MessageE{}) // shared with <process>
				},
			}

			app = &Application{
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
					MessageRoleMap{
						MessageAType: CommandMessageRole,
						MessageBType: EventMessageRole,
						MessageCType: CommandMessageRole,
						MessageDType: EventMessageRole,
						MessageEType: EventMessageRole,
						MessageFType: EventMessageRole,
						MessageTType: TimeoutMessageRole,
					},
				))
			})

			It("sets the consumers map", func() {
				Expect(cfg.Consumers).To(Equal(
					map[MessageType][]HandlerConfig{
						MessageAType: {cfg.HandlersByName["<aggregate>"]},
						MessageBType: {cfg.HandlersByName["<process>"]},
						MessageCType: {cfg.HandlersByName["<integration>"]},
						MessageDType: {cfg.HandlersByName["<projection>"]},
						MessageEType: {cfg.HandlersByName["<process>"], cfg.HandlersByName["<projection>"]},
						MessageTType: {cfg.HandlersByName["<process>"]},
					},
				))
			})

			It("sets the producers map", func() {
				Expect(cfg.Producers).To(Equal(
					map[MessageType][]HandlerConfig{
						MessageCType: {cfg.HandlersByName["<process>"]},
						MessageEType: {cfg.HandlersByName["<aggregate>"]},
						MessageFType: {cfg.HandlersByName["<integration>"]},
						MessageTType: {cfg.HandlersByName["<process>"]},
					},
				))
			})

			Describe("func Identity()", func() {
				It("returns the app identity", func() {
					Expect(cfg.Identity()).To(Equal(
						MustNewIdentity("<app>", "<app-key>"),
					))
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
						`*fixtures.Application.Configure() did not call ApplicationConfigurer.Identity()`,
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
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
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
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
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
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
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
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
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
					c.ConsumesEventType(MessageB{})
					c.ProducesCommandType(MessageC{})
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
					c.ConsumesEventType(MessageB{})
					c.ProducesCommandType(MessageC{})
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
					c.ConsumesCommandType(MessageC{})
					c.ProducesEventType(MessageF{})
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
					c.ConsumesCommandType(MessageC{})
					c.ProducesEventType(MessageF{})
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
					c.ConsumesEventType(MessageD{})
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
					c.ConsumesEventType(MessageD{})
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
				c.ConsumesCommandType(MessageA{}) // conflict with <aggregate>
				c.ProducesEventType(MessageF{})
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
				c.ConsumesCommandType(MessageC{})
				c.ProducesEventType(MessageE{}) // conflict with <aggregate>
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				Error(
					`the "<integration>" handler can not produce fixtures.MessageE events because they are already produced by "<aggregate>"`,
				),
			))
		})

		It("does not return an error when the app contains multiple processes that schedule the same timeout", func() {
			process1 := &ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-1>", "<process-1-key>")
					c.ConsumesEventType(MessageB{})
					c.ProducesCommandType(MessageC{})
					c.SchedulesTimeoutType(MessageT{})
				},
			}

			process2 := &ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-2>", "<process-2-key>")
					c.ConsumesEventType(MessageB{})
					c.ProducesCommandType(MessageC{})
					c.SchedulesTimeoutType(MessageT{})
				},
			}

			app := &Application{
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
					c.ConsumesEventType(MessageA{}) // conflict with <aggregate>
					c.ProducesCommandType(MessageC{})
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
					c.ConsumesEventType(MessageB{})
					c.ProducesCommandType(MessageE{}) // conflict with <aggregate>
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
