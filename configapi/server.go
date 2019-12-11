package configapi

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/configapi/internal/pb"
	"google.golang.org/grpc"
)

// RegisterServer registers a config server for the config applications.
func RegisterServer(
	s *grpc.Server,
	apps ...configkit.Application,
) {
	svr := &server{}

	for _, in := range apps {
		out, err := marshalApplication(in)
		if err != nil {
			panic(err)
		}

		svr.ListApplicationIdentitiesResponse.Identities = append(
			svr.ListApplicationIdentitiesResponse.Identities,
			out.Identity,
		)

		svr.ListApplicationsResponse.Applications = append(
			svr.ListApplicationsResponse.Applications,
			out,
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
