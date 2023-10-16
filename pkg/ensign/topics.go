package ensign

import (
	"context"
	"strings"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListTopics associated with the project ID in the claims of the request. This unary
// request is paginated to prevent a huge amount of data transfer.
//
// Permissions: topics:read
func (s *Server) ListTopics(ctx context.Context, in *api.PageInfo) (out *api.TopicsPage, err error) {
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to list the topics in the project
	if !claims.HasPermission(permissions.ReadTopics) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(claims.ProjectID); err != nil || ulids.IsZero(projectID) {
		sentry.Warn(ctx).Err(err).Msg("could not parse projectID from claims")
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// Fetch the results from the database
	iter := s.meta.ListTopics(projectID)
	defer iter.Release()

	if out, err = iter.NextPage(in); err != nil {
		// TODO: handle invalid argument errors
		sentry.Error(ctx).Err(err).Msg("could not process next page of results from the database")
		return nil, status.Error(codes.Internal, "unable to process list topics request")
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not retrieve topics from the database")
		return nil, status.Error(codes.Internal, "unable to process list topics request")
	}
	return out, nil
}

// CreateTopic is a user-facing request to create a Topic. Ensign first verifies that
// the topic is eligible to be created, then stores the topic in a pending state to disk
// and returns success to the user. Afterwards, Ensign sends a notification to the
// placement service in order to figure out where the topic should be assigned so that
// it can start receiving events.
//
// Permissions: topics:create
func (s *Server) CreateTopic(ctx context.Context, in *api.Topic) (_ *api.Topic, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to create the topic in the project
	if !claims.HasPermission(permissions.CreateTopics) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	if len(in.ProjectId) == 0 {
		// If the projectID hasn't been specified in the request, set it from the claims
		if projectID := claims.ParseProjectID(); !ulids.IsZero(projectID) {
			in.ProjectId = projectID.Bytes()
		} else {
			return nil, status.Error(codes.InvalidArgument, "missing project id field")
		}
	}

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(in.ProjectId); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not parse projectId from user create topic request")
		return nil, status.Error(codes.InvalidArgument, "invalid project id field")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// TODO: set the topic status as pending

	// Create the topic in the store: note that the store will validate the topic
	if err = s.meta.CreateTopic(in); err != nil {
		if errors.Is(err, errors.ErrUniqueTopicName) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, errors.ErrInvalidTopic) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		sentry.Error(ctx).Err(err).Msg("could not create topic")
		return nil, status.Error(codes.Internal, "could not process create topic request")
	}

	// TODO: send topic to placement service

	// The store method updates the in reference in place, preventing an allocation.
	return in, nil
}

// RetrieveTopic is a user-face request to fetch a single Topic and is typically used
// for existence checks; e.g. does this topic exist or not. The user only has to specify
// a TopicID in the request and then a complete topic is returned. If the topic is not
// found a status error with codes.NotFound is returned.
//
// Permissions: topics:read
func (s *Server) RetrieveTopic(ctx context.Context, in *api.Topic) (out *api.Topic, err error) {
	// Collect credentials from the context
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to retrieve the topic in the project
	if !claims.HasPermission(permissions.ReadTopics) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	if len(in.Id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing topic id field")
	}

	var topicID ulid.ULID
	if topicID, err = ulids.Parse(in.Id); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not parse topic id field from user retrieve topic request")
		return nil, status.Error(codes.InvalidArgument, "invalid topic id field")
	}

	// Retrieve the topic from the store
	if out, err = s.meta.RetrieveTopic(topicID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		sentry.Error(ctx).Err(err).Msg("could not retrieve topic")
		return nil, status.Error(codes.Internal, "could not complete retrieve topic request")
	}

	// Ensure the topic that is retrieved is in the same project as the claims.
	// Should not be able to delete a topic from another project
	// TODO: should this be part of the retrieve process, e.g. should retrieve take a projectID?
	// NOTE: because the object key is projectID:topicID, we could do a direct Get instead of a retrieve here
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(out.ProjectId); err != nil {
		sentry.Warn(ctx).Err(err).Msg("topic retrieved from database has unparsable project ID field")
		return nil, status.Error(codes.NotFound, "could not complete retrieve topic request")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}

	return out, nil
}

// DeleteTopic is a user-facing request to modify a topic and either archive it, which
// will make it read-only permanently or to destroy it, which will also have the effect
// of removing all of the data in the topic. This method is a stateful method, e.g. the
// topic will be updated to the current status then the Ensign placement server will
// take action from there.
//
// Permissions: topics:edit (archive) or topics:destroy (destroy)
func (s *Server) DeleteTopic(ctx context.Context, in *api.TopicMod) (out *api.TopicStatus, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
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
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// An ID is required to be able to delete a topic
	if strings.TrimSpace(in.Id) == "" {
		return nil, status.Error(codes.InvalidArgument, "missing id field")
	}

	var topicID ulid.ULID
	if topicID, err = ulids.Parse(in.Id); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not parse id from user delete topic request")
		return nil, status.Error(codes.InvalidArgument, "invalid id field")
	}

	// Update the local database with the record
	// HACK: this mechanism in a single node is not concurrency safe but will provide the functionality
	var topic *api.Topic
	if topic, err = s.meta.RetrieveTopic(topicID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		sentry.Error(ctx).Err(err).Msg("could not retrieve topic for deletion")
		return nil, status.Error(codes.Internal, "could not process delete topic request")
	}

	// Should not be able to delete a topic from another project
	// TODO: should this be part of the retrieve process, e.g. should retrieve take a projectID?
	// NOTE: because the object key is projectID:topicID, we could do a direct Get instead of a retrieve here
	var projectID ulid.ULID
	if projectID, err = topic.ParseProjectID(); err != nil {
		sentry.Error(ctx).Err(err).Str("topic_id", topicID.String()).Bytes("project_id", topic.ProjectId).Msg("unable to parse project id defined on stored topic")
		return nil, status.Error(codes.Internal, "could not process delete topic request")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}

	// TODO: send topic deletion to the placement service
	// TODO: if destroy, create a job to delete all the data for the specified topic
	// TODO: update broker to prevent any additional writes.

	out = &api.TopicStatus{Id: topicID.String()}
	switch in.Operation {
	case api.TopicMod_ARCHIVE:
		topic.Readonly = true
		out.State = api.TopicState_READONLY

		if err = s.meta.UpdateTopic(topic); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not update topic as readonly")
			return nil, status.Error(codes.Internal, "could not process delete topic request")
		}

	case api.TopicMod_DESTROY:
		out.State = api.TopicState_DELETING

		if err = s.meta.DeleteTopic(topicID); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not delete topic from meta store")
			return nil, status.Error(codes.Internal, "could not process delete topic request")
		}

		// TODO: queue a job to delete all events associated with the topic
	}

	return out, nil
}

// SetTopicPolicy allows the user to specify topic management policies such as
// deduplication or sharding. If the topic is already in the policies specified, then
// READY is returned. Otherwise a job is queued to modify the topic policy and PENDING
// is returned. This is a patch endpoint, so if any policy is set to UNKNOWN, then it is
// ignored; only named policies initiate changes on the topic.
//
// Permissions: topics:edit
func (s *Server) SetTopicPolicy(ctx context.Context, in *api.TopicPolicy) (out *api.TopicStatus, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to edit the topic policies
	if !claims.HasPermission(permissions.EditTopics) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// An ID is required to modify the topic policy
	if strings.TrimSpace(in.Id) == "" {
		return nil, status.Error(codes.InvalidArgument, "missing id field")
	}

	// If no policy change has been specified, return invalid argument
	if in.DeduplicationPolicy.Strategy == api.Deduplication_UNKNOWN && in.ShardingStrategy == api.ShardingStrategy_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "no policies defined to set on topic")
	}

	var topicID ulid.ULID
	if topicID, err = ulids.Parse(in.Id); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not parse topic id from user set topic policy request")
		return nil, status.Error(codes.InvalidArgument, "invalid id field")
	}

	// Retrieve the topic from the database
	var topic *api.Topic
	if topic, err = s.meta.RetrieveTopic(topicID); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		sentry.Error(ctx).Err(err).Msg("could not retrieve topic for policy change")
		return nil, status.Error(codes.Internal, "could not process set topic policy request")
	}

	// Should not be able to delete a topic from another project
	// TODO: should this be part of the retrieve process, e.g. should retrieve take a projectID?
	// NOTE: because the object key is projectID:topicID, we could do a direct Get instead of a retrieve here
	var projectID ulid.ULID
	if projectID, err = topic.ParseProjectID(); err != nil {
		sentry.Error(ctx).Err(err).Str("topic_id", topicID.String()).Bytes("project_id", topic.ProjectId).Msg("unable to parse project id defined on stored topic")
		return nil, status.Error(codes.Internal, "could not process set topic policy request")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}

	// NOTE: updating the sharding strategy is currently not supported
	// TODO: support updating the sharding strategy!
	if in.ShardingStrategy != api.ShardingStrategy_UNKNOWN {
		return nil, status.Error(codes.Unimplemented, "changing the sharding strategy of a topic is currently not supported")
	}

	// If there is no change to the deduplication strategy, then return READY
	if topic.Deduplication.Equals(in.DeduplicationPolicy) {
		return &api.TopicStatus{Id: topicID.String(), State: topic.Status}, nil
	}

	// Handle any changes to the deduplication strategy of the topic
	// TODO: normalize/validate incoming policy
	// TODO: update the policy of the topic
	// TODO: Update the broker with the new policy
	// TODO: create a job to update the topic's policy
	return nil, status.Error(codes.Unimplemented, "changing deduplication policies is coming soon!")
}

// TopicNames returns a paginated response that maps topic names to IDs for all topics
// in the project. The claims must have any of the topics:read, publisher, or subscriber
// permissions in order to access this endpoint.
//
// Permissions: topics:read OR publisher OR subscriber
func (s *Server) TopicNames(ctx context.Context, in *api.PageInfo) (out *api.TopicNamesPage, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to retrieve the topic in the project
	if !claims.HasAnyPermission(permissions.ReadTopics, permissions.Publisher, permissions.Subscriber) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		// If there is an invalid projectID in the claims return empty list
		return &api.TopicNamesPage{TopicNames: nil, NextPageToken: ""}, nil
	}

	iter := s.meta.ListTopicNames(projectID)
	defer iter.Release()

	if out, err = iter.NextPage(in); err != nil {
		// TODO: handle invalid argument errors
		sentry.Error(ctx).Err(err).Msg("could not process next page of topic name index results from the database")
		return nil, status.Error(codes.Internal, "unable to process list topic names request")
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not retrieve topic names from the database")
		return nil, status.Error(codes.Internal, "unable to process list topic names request")
	}
	return out, nil
}

// TopicExists does a quick check to see if the topic ID or name exists in the project
// and returns a simple yes or no bool with the original query. If both ID and name are
// specified then this method checks if a topic with the specified name has the specified
// ID (e.g. it is a more strict existence check). The claims must have any of the
// topics:read, publisher, or subscriber permissions in order to access this endpoint.
//
// Permissions: topics:read OR publisher OR subscriber
func (s *Server) TopicExists(ctx context.Context, in *api.TopicName) (out *api.TopicExistsInfo, err error) {
	// Collect credentials from the context
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against future development.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has the permissions to retrieve the topic in the project
	if !claims.HasAnyPermission(permissions.ReadTopics, permissions.Publisher, permissions.Subscriber) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	if in.ProjectId == "" {
		// If the projectID hasn't been specified in the request, set it from the claims
		if projectID := claims.ParseProjectID(); !ulids.IsZero(projectID) {
			in.ProjectId = projectID.String()
		} else {
			return nil, status.Error(codes.InvalidArgument, "missing project id field")
		}
	}

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(in.ProjectId); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not parse projectId from user topic exists request")
		return nil, status.Error(codes.InvalidArgument, "invalid project id field")
	}

	if !claims.ValidateProject(projectID) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// Check the user input isn't completely empty
	if in.Name == "" && in.TopicId == "" {
		return nil, status.Error(codes.InvalidArgument, "must specify either or both topic name and id")
	}

	if out, err = s.meta.TopicExists(in); err != nil {
		// TODO: handle invalid argument errors
		sentry.Error(ctx).Err(err).Msg("could not check topic existence in topic names index")
		return nil, status.Error(codes.Internal, "unable to process topic exists request")
	}
	return out, nil
}
