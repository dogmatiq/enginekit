package api

import (
	"context"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/api/internal/pb"
	"github.com/dogmatiq/enginekit/marshaling"
	"google.golang.org/grpc"
)

// RegisterServer registers a config server for the config applications.
func RegisterServer(
	s *grpc.Server,
	m *marshaling.Marshaler,
	apps []*config.ApplicationConfig,
) {
	svr := &server{}

	for _, cfg := range apps {
		app := marshalApplication(m, cfg)

		svr.ListApplicationIdentitiesResponse.Identities = append(
			svr.ListApplicationIdentitiesResponse.Identities,
			app.Identity,
		)

		svr.ListApplicationsResponse.Applications = append(
			svr.ListApplicationsResponse.Applications,
			app,
		)
	}

	pb.RegisterConfigServer(s, svr)
}

var _ pb.ConfigServer = (*server)(nil)

type server struct {
	pb.ListApplicationIdentitiesResponse
	pb.ListApplicationsResponse
}

// ListApplicationIdentities returns the identity of all applications.
func (s *server) ListApplicationIdentities(
	ctx context.Context,
	req *pb.ListApplicationIdentitiesRequest,
) (*pb.ListApplicationIdentitiesResponse, error) {
	return &s.ListApplicationIdentitiesResponse, nil
}

// ListApplications returns the full configuration of all applications.
func (s *server) ListApplications(
	ctx context.Context,
	req *pb.ListApplicationsRequest,
) (*pb.ListApplicationsResponse, error) {
	return &s.ListApplicationsResponse, nil
}
