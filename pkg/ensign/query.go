package ensign

import (
	"context"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/ensql"
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
	if !claims.HasPermission(permissions.Subscriber) {
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Msg("no project id specified in claims")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
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
	if topicID, err = s.meta.LookupTopicName(query.Topic.Topic, projectID); err != nil {
		log.Debug().Err(err).Str("topic", query.Topic.Topic).Msg("could not lookup topic in query")
		if errors.Is(err, errors.ErrNotFound) {
			return status.Error(codes.InvalidArgument, "unknown topic in query")
		}

		sentry.Error(ctx).Err(err).Msg("could not lookup topic name")
		return status.Error(codes.Internal, "could not execute query")
	}

	// TODO: Begin simple execution of query
	log.Debug().Str("query", query.Raw).Str("topic", topicID.String()).Msg("starting ensql query execution")
	return status.Error(codes.Unimplemented, "query execution not implemented yet")
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
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Msg("no project id specified in claims")
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	return nil, status.Error(codes.Unimplemented, "explain query is not implemented yet")
}
