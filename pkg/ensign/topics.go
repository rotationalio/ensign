package ensign

import (
	"context"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
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

	// Verify that the user has the permissions to create the topic in the project

	// TODO: set the topic status as pending

	// Create the topic in the store: note that the store will validate the topic
	if err = s.meta.CreateTopic(in); err != nil {
		if errors.Is(err, errors.ErrInvalidTopic) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: send topic to placement service

	return in, nil
}
