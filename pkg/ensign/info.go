package ensign

import (
	"context"
	"encoding/base64"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The Info RPC returns a summary of the current state of the project retrieved from the
// claims of the request. This RPC requires the ReadTopics permission in order to return
// any information. The status request can be filtered by a list of topics to specify
// exactly what statistics are returned.
func (s *Server) Info(ctx context.Context, in *api.InfoRequest) (out *api.ProjectInfo, err error) {
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// This should never happen but the check prevents nil panics.
		sentry.Error(ctx).Msg("could not get user claims from authenticated request")
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// The user must have the read topics permission to view project info
	if !claims.HasPermission(permissions.ReadTopics) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Get the project from the claims to get the info for
	var projectID ulid.ULID
	if projectID, err = ulids.Parse(claims.ProjectID); err != nil || ulids.IsZero(projectID) {
		sentry.Warn(ctx).Err(err).Msg("could not parse projectID from claims")
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Create a topic filter function that returns true if the topic should be included
	// in the returned project statistics.
	var includeTopic func(topicID []byte) bool
	if len(in.Topics) == 0 {
		// Do not filter any topics if an empty request was sent
		includeTopic = func(topicID []byte) bool { return true }
	} else {
		// Create a set of all the topicIDs the user specified in the request
		included := make(map[ulid.ULID]struct{})
		for _, topicID := range in.Topics {
			var tid ulid.ULID
			if tid, err = ulids.Parse(topicID); err != nil {
				return nil, status.Error(codes.InvalidArgument, "could not parse topic id in info request filter")
			}
			included[tid] = struct{}{}
		}

		includeTopic = func(topicID []byte) bool {
			tid, _ := ulids.Parse(topicID)
			_, ok := included[tid]
			return ok
		}
	}

	// Prepare the response
	out = &api.ProjectInfo{
		ProjectId: projectID.String(),
	}

	// Loop through all topics in the project and get the info for them.
	iter := s.meta.ListTopics(projectID)
	defer iter.Release()

	for iter.Next() {
		var topic *api.Topic
		if topic, err = iter.Topic(); err != nil {
			sentry.Warn(ctx).Err(err).Str("topicKey", base64.StdEncoding.EncodeToString(iter.Key())).Msg("could not deserialize topic from database")
			continue
		}

		// filter topics based on the user's request
		if !includeTopic(topic.Id) {
			continue
		}

		// increment the project statistics
		out.Topics++
		if topic.Readonly {
			out.ReadonlyTopics++
		}

		out.Events += topic.Offset
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not retrieve topics from database")
		return nil, status.Error(codes.Internal, "unable to process project info request")
	}

	return out, nil
}
