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
		log.Warn().Err(err).Msg("could not parse projectId from user create topic request")
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

		log.Error().Err(err).Msg("could not create topic")
		return nil, status.Error(codes.Internal, "could not process create topic request")
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

	var topicID ulid.ULID
	if topicID, err = ulids.Parse(in.Id); err != nil {
		log.Warn().Err(err).Msg("could not parse id from user delete topic request")
		return nil, status.Error(codes.InvalidArgument, "invalid id field")
	}

	// TODO: send topic deletion to the placement service
	// TODO: if destroy, create a job to delete all the data for the specified topic

	// Update the local database with the record
	// HACK: this mechanism in a single node is not concurrency safe but will provide the functionality
	var topic *api.Topic
	if topic, err = s.meta.RetrieveTopic(topicID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		log.Error().Err(err).Msg("could not retrieve topic for deletion")
		return nil, status.Error(codes.Internal, "could not process delete topic request")
	}

	// Should not be able to delete a topic from another project
	// TODO: should this be part of the retrieve process, e.g. should retrieve take a projectID?
	// NOTE: because the object key is projectID:topicID, we could do a direct Get instead of a retrieve here
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(topic.ProjectId); err != nil {
		log.Error().Err(err).Str("topic_id", topicID.String()).Bytes("project_id", topic.ProjectId).Msg("unable to parse project id defined on stored topic")
		return nil, status.Error(codes.Internal, "could not process delete topic request")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}

	out = &api.TopicTombstone{Id: topicID.String()}
	switch in.Operation {
	case api.TopicMod_ARCHIVE:
		topic.Readonly = true
		out.State = api.TopicTombstone_READONLY

		if err = s.meta.UpdateTopic(topic); err != nil {
			log.Error().Err(err).Msg("could not update topic as readonly")
			return nil, status.Error(codes.Internal, "could not process delete topic request")
		}
	case api.TopicMod_DESTROY:
		out.State = api.TopicTombstone_DELETING

		if err = s.meta.DeleteTopic(topicID); err != nil {
			log.Error().Err(err).Msg("could not delete topic from meta store")
			return nil, status.Error(codes.Internal, "could not process delete topic request")
		}
	}

	return out, nil
}
