package ensign_test

import (
	"context"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
)

func (s *serverTestSuite) TestCreateTopic() {
	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GTW1R9MH8723JQDRMFE16CZ7"),
		Name:      "testing.testapp.test",
	}

	// Should not be able to create a topic when not authenticated

	// Should not be able to create a topic in the wrong project

	// Happy path: should be able to create a valid topic
	out, err := s.client.CreateTopic(context.Background(), topic)
	s.NoError(err, "could not execute create topic request")

	s.False(ulids.IsZero(ulids.MustParse(out.Id)))
	s.Equal(topic.ProjectId, out.ProjectId)
	s.Equal(topic.Name, out.Name)
	s.NotEmpty(out.Created)
	s.NotEmpty(out.Modified)

	// Should not be able to create a topic without a name
	topic.Name = ""
	_, err = s.client.CreateTopic(context.Background(), topic)
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing name field")
}
