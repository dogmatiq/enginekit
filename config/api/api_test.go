package api_test

import (
	"context"
	"net"
	"reflect"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/api"
	"github.com/dogmatiq/marshalkit"
	"github.com/dogmatiq/marshalkit/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("type Client", func() {
	var (
		ctx        context.Context
		cancel     func()
		app1, app2 dogma.Application
		cfg1, cfg2 *config.ApplicationConfig
		marshaler  *marshalkit.Marshaler
		listener   net.Listener
		gserver    *grpc.Server
		client     *Client
	)

	BeforeEach(func() {
		var fn func()
		ctx, fn = context.WithTimeout(context.Background(), 1*time.Second)
		cancel = fn // bypass linter warning about cancel being unused

		var err error

		app1 = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-1>", "<app-key-1>")

				c.RegisterAggregate(&AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", "<aggregate-key>")
						c.ConsumesCommandType(MessageC{})
						c.ProducesEventType(MessageE{})
					},
				})

				c.RegisterProcess(&ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", "<process-key>")
						c.ConsumesEventType(MessageE{})
						c.ProducesCommandType(MessageC{})
						c.SchedulesTimeoutType(MessageT{})
					},
				})
			},
		}

		app2 = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-2>", "<app-key-2>")

				c.RegisterIntegration(&IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", "<integration-key>")
						c.ConsumesCommandType(MessageI{})
						c.ProducesEventType(MessageJ{})
					},
				})

				c.RegisterProjection(&ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", "<projection-key>")
						c.ConsumesEventType(MessageE{})
						c.ConsumesEventType(MessageJ{})
					},
				})
			},
		}

		cfg1, err = config.NewApplicationConfig(app1)
		Expect(err).ShouldNot(HaveOccurred())
		cfg1.Accept(context.Background(), stripEntities{})

		cfg2, err = config.NewApplicationConfig(app2)
		Expect(err).ShouldNot(HaveOccurred())
		cfg2.Accept(context.Background(), stripEntities{})

		var types []reflect.Type
		for mt := range cfg1.Roles {
			types = append(types, mt.ReflectType())
		}
		for mt := range cfg2.Roles {
			types = append(types, mt.ReflectType())
		}
		marshaler, err = marshalkit.NewMarshaler(
			types,
			[]marshalkit.Codec{
				&json.Codec{},
			},
		)
		Expect(err).ShouldNot(HaveOccurred())

		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()
		RegisterServer(gserver, marshaler, cfg1, cfg2)

		go gserver.Serve(listener)

		conn, err := grpc.Dial(
			listener.Addr().String(),
			grpc.WithInsecure(),
		)
		Expect(err).ShouldNot(HaveOccurred())

		client = &Client{
			conn,
			marshaler,
		}
	})

	AfterEach(func() {
		if listener != nil {
			listener.Close()
		}

		if gserver != nil {
			gserver.Stop()
		}

		cancel()
	})

	Describe("func ListApplicationIdentities()", func() {
		It("returns the application identities", func() {
			idents, err := client.ListApplicationIdentities(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(idents).To(ConsistOf(
				config.MustNewIdentity("<app-1>", "<app-key-1>"),
				config.MustNewIdentity("<app-2>", "<app-key-2>"),
			))
		})

		It("returns an error if the gRPC call fails", func() {
			gserver.Stop()
			_, err := client.ListApplicationIdentities(ctx)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func ListApplications()", func() {
		It("returns the application configurations", func() {
			configs, err := client.ListApplications(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(configs).To(ConsistOf(
				cfg1,
				cfg2,
			))
		})

		It("returns an error if the gRPC call fails", func() {
			gserver.Stop()
			_, err := client.ListApplications(ctx)
			Expect(err).Should(HaveOccurred())
		})
	})
})

// stripEntities is a config visitor that sets all application/handler values to
// nil, to mimic how the config would appear when obtained via the API.
type stripEntities struct{}

func (v stripEntities) VisitApplicationConfig(
	ctx context.Context,
	cfg *config.ApplicationConfig,
) error {
	cfg.Application = nil

	for _, h := range cfg.HandlersByKey {
		h.Accept(ctx, v)
	}

	return nil
}

func (v stripEntities) VisitAggregateConfig(
	ctx context.Context,
	cfg *config.AggregateConfig,
) error {
	cfg.Handler = nil
	return nil
}

func (v stripEntities) VisitProcessConfig(
	ctx context.Context,
	cfg *config.ProcessConfig,
) error {
	cfg.Handler = nil
	return nil
}

func (v stripEntities) VisitIntegrationConfig(
	ctx context.Context,
	cfg *config.IntegrationConfig,
) error {
	cfg.Handler = nil
	return nil
}

func (v stripEntities) VisitProjectionConfig(
	ctx context.Context,
	cfg *config.ProjectionConfig,
) error {
	cfg.Handler = nil
	return nil
}
