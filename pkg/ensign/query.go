package ensign

import (
	"context"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EnSQL parses an incoming query and executes the query request, sending all results
// onto the query stream. The EnSQL query is guaranteed to terminate, e.g. it is not a
// long running query that waits for subscription events to come from publishers. Once
// the query has been completed the stream will close. Errors are returned for standard
// SQL operations errors - for example if the query cannot be parsed or no results would
// be returned from the query.
//
// Permissions: subscriber
func (s *Server) EnSQL(in *api.Query, stream api.Ensign_EnSQLServer) (err error) {
	ctx := stream.Context()
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// This should never happen but the check prevents nil panics.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return status.Error(codes.Unauthenticated, "missing credentials")
	}

	// The user must have the subscriber permission to execute a query
	// TODO: remove the read topics permission when we update Quarterdeck permissions.
	if !claims.HasAnyPermission(permissions.Subscriber, permissions.ReadTopics) {
		return status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Msg("no project id specified in claims")
		return status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	// Parse the incoming query
	if in.Query == "" {
		return status.Error(codes.InvalidArgument, "invalid query")
	}

	var query ensql.Query
	if query, err = ensql.Parse(in.Query); err != nil {
		log.Debug().Err(err).Str("query", in.Query).Msg("could not parse query")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// Identify the topic in the query
	var topicID ulid.ULID
	if topicID, err = s.meta.LookupTopicID(query.Topic.Topic, projectID); err != nil {
		log.Debug().Err(err).Str("topic", query.Topic.Topic).Msg("could not lookup topic in query")
		if errors.Is(err, errors.ErrNotFound) {
			return status.Error(codes.InvalidArgument, "unknown topic in query")
		}

		sentry.Error(ctx).Err(err).Msg("could not lookup topic name")
		return status.Error(codes.Internal, "could not execute query")
	}

	// Begin simple execution of query
	log.Debug().Str("query", query.Raw).Str("topic", topicID.String()).Msg("starting ensql query execution")
	events := s.data.List(topicID)
	defer events.Release()

	// Skip over events in the offset
	// NOTE: offset will include duplicates when skipping over ...
	// TODO: this is very slow, we need to do a binary search for the offset instead.
	if query.HasOffset {
		for i := uint64(0); i < query.Offset; i++ {
			if !events.Next() {
				break
			}
		}
	}

	nSent := uint64(0)
	for events.Next() {
		var event *api.EventWrapper
		if event, err = events.Event(); err != nil {
			sentry.Error(ctx).Bytes("key", events.Key()).Err(err).Msg("could not parse event")
			continue
		}

		// Skip over duplicates unless specified by the query
		if !in.IncludeDuplicates && event.IsDuplicate {
			continue
		}

		// If we're including duplicates, and the event is a duplicate, then dereference
		// the duplicate from the database so there is correct event information.
		if event.IsDuplicate {
			var target *api.EventWrapper
			if target, err = s.data.Retrieve(topicID, rlid.RLID(event.DuplicateId)); err != nil {
				sentry.Error(ctx).Bytes("duplicate_id", event.DuplicateId).Str("topic_id", topicID.String()).Msg("could not fetch duplicate reference target")
				continue
			}

			if err = event.DuplicateFrom(target); err != nil {
				sentry.Error(ctx).Bytes("duplicate_id", event.DuplicateId).Str("topic_id", topicID.String()).Msg("could not dereference duplicate event")
				continue
			}
		}

		// TODO: evaluate WHERE clause
		if err = stream.Send(event); err != nil {
			if streamClosed(err) {
				log.Debug().Msg("publish stream closed by client")
				return nil
			}
			sentry.Warn(ctx).Err(err).Msg("ensql query stream crashed")
			return status.Error(codes.Aborted, "query stream aborted")
		}

		// Check the limit to return a fixed number of events
		nSent++
		if query.HasLimit {
			if nSent >= query.Limit {
				break
			}
		}
	}

	if err := events.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not retrieve events from database")
		return status.Error(codes.Internal, "could not execute query")
	}

	return nil
}

// Explain parses the input query and returns an explanation consisting of the query
// plan and approximate number of results any any possible errors.
//
// Permissions: subscriber
// TODO: implement explanation
func (s *Server) Explain(ctx context.Context, in *api.Query) (out *api.QueryExplanation, err error) {
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// This should never happen but the check prevents nil panics.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// The user must have the subscriber permission to execute a query
	if !claims.HasPermission(permissions.Subscriber) {
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Msg("no project id specified in claims")
		return nil, status.Error(codes.PermissionDenied, "not authorized to perform this action")
	}

	return nil, status.Error(codes.Unimplemented, "explain query is not implemented yet")
}
