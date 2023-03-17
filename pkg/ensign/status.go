package ensign

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Defines the expected heartbeat interval for remote clients.
const (
	minHBInterval = 1 * time.Minute
	maxHBInterval = 1 * time.Hour
)

// Status implements a simple heartbeat mechanism for checking on the state of the
// Ensign server and making sure that the node is up and responding.
func (s *Server) Status(ctx context.Context, in *api.HealthCheck) (out *api.ServiceState, err error) {
	out = &api.ServiceState{
		Status:    api.ServiceState_HEALTHY,
		Version:   pkg.Version(),
		Uptime:    durationpb.New(time.Since(s.started)),
		NotBefore: timestamppb.New(time.Now().Add(minHBInterval)),
		NotAfter:  timestamppb.New(time.Now().Add(maxHBInterval)),
	}

	if s.conf.Maintenance {
		out.Status = api.ServiceState_MAINTENANCE
	}
	return out, nil
}
