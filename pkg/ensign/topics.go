package ensign

import (
	"context"
	"strings"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateTopic is a user-facing request to create a Topic. Ensign first verifies that
// the topic is eligible to be created, then stores the topic in a pending state to disk
// and returns success to the user. Afterwards, Ensign sends a notification to the
// placement service in order to figure out where the topic should be assigned so that
// it can start receiving events.
func (s *Server) CreateTopic(ctx context.Context, in *api.Topic) (_ *api.Topic, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to create the topic in the project
	if !claims.HasPermission(permissions.CreateTopics) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	if len(in.ProjectId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing project id field")
	}

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(in.ProjectId); err != nil {
		log.Warn().Err(err).Msg("could not parse projectId from user request")
		return nil, status.Error(codes.InvalidArgument, "invalid project id field")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// TODO: set the topic status as pending

	// Create the topic in the store: note that the store will validate the topic
	if err = s.meta.CreateTopic(in); err != nil {
		if errors.Is(err, errors.ErrInvalidTopic) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: send topic to placement service

	// The store method updates the in reference in place, preventing an allocation.
	return in, nil
}

// DeleteTopic is a user-facing request to modify a topic and either archive it, which
// will make it read-only permanently or to destroy it, which will also have the effect
// of removing all of the data in the topic. This method is a stateful method, e.g. the
// topic will be updated to the current status then the Ensign placement server will
// take action from there.
func (s *Server) DeleteTopic(ctx context.Context, in *api.TopicMod) (out *api.TopicTombstone, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// If the modification operation is archive then the user needs the EditTopics
	// permission, otherwise they need the DestroyTopics permission
	var permission string
	switch in.Operation {
	case api.TopicMod_ARCHIVE:
		permission = permissions.EditTopics
	case api.TopicMod_DESTROY:
		permission = permissions.DestroyTopics
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid operation field")
	}

	// Verify the user has the permissions to modify the topic in the project.
	if !claims.HasPermission(permission) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// An ID is required to be able to delete a topic
	if strings.TrimSpace(in.Id) == "" {
		return nil, status.Error(codes.InvalidArgument, "missing id field")
	}

	return out, nil
}
