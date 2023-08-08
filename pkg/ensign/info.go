package ensign

import (
	"context"

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
	if !claims.HasAllPermissions(permissions.ReadTopics, permissions.ReadMetrics) {
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Get the project from the claims to get the info for
	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Msg("no project id specified in claims")
		return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Create a topic filter function that returns true if the topic should be included
	// in the returned project statistics.
	var includeTopic func(topicID ulid.ULID) bool
	if len(in.Topics) == 0 {
		// Do not filter any topics if an empty request was sent
		includeTopic = func(ulid.ULID) bool { return true }
	} else {
		// Create a set of all the topicIDs the user specified in the request
		included := make(map[ulid.ULID]struct{}, len(in.Topics))
		for _, topicID := range in.Topics {
			var tid ulid.ULID
			if tid, err = ulids.Parse(topicID); err != nil {
				return nil, status.Error(codes.InvalidArgument, "could not parse topic id in info request filter")
			}
			included[tid] = struct{}{}
		}

		includeTopic = func(topicID ulid.ULID) bool {
			_, ok := included[topicID]
			return ok
		}
	}

	// Prepare the response
	out = &api.ProjectInfo{
		ProjectId: projectID.Bytes(),
		Topics:    make([]*api.TopicInfo, 0),
	}

	// Loop through all topics in the project and get the info for them.
	iter := s.meta.ListTopics(projectID)
	defer iter.Release()

	for iter.Next() {
		var topic *api.Topic
		if topic, err = iter.Topic(); err != nil {
			sentry.Warn(ctx).Err(err).Bytes("topic_key", iter.Key()).Msg("could not deserialize topic from database")
			continue
		}

		// Parse the topicID
		var topicID ulid.ULID
		if topicID, err = ulids.Parse(topic.Id); err != nil {
			sentry.Warn(ctx).Err(err).Bytes("topic_id", topic.Id).Msg("could not parse topicID into ulid")
			continue
		}

		// filter topics based on the user's request
		if !includeTopic(topicID) {
			continue
		}

		// increment the project statistics
		out.NumTopics++
		if topic.Readonly {
			out.NumReadonlyTopics++
		}

		// get topic info for the specified topic and increment results
		var info *api.TopicInfo
		if info, err = s.meta.TopicInfo(topicID); err != nil {
			sentry.Warn(ctx).Err(err).Str("topic_id", topicID.String()).Msg("could not get topic info for topic")
			continue
		}

		out.Topics = append(out.Topics, info)
		out.Events += info.Events
		out.Duplicates += info.Duplicates
		out.DataSizeBytes += info.DataSizeBytes
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not retrieve topics from database")
		return nil, status.Error(codes.Internal, "unable to process project info request")
	}

	return out, nil
}
