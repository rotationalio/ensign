package ensign_test

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
)

func (s *serverTestSuite) TestCreateTopic() {
	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GTW1R9MH8723JQDRMFE16CZ7"),
		Name:      "testing.testapp.test",
	}

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "",
		},
	}

	// Should not be able to create a topic when not authenticated

	// Should not be able to create a topic in the wrong project

	// Happy path: should be able to create a valid topic
	claims.ProjectID = "01GTW1R9MH8723JQDRMFE16CZ7"
	token, err := s.quarterdeck.CreateAccessToken(claims)
	s.NoError(err, "could not create valid claims for the user")

	out, err := s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.NoError(err, "could not execute create topic request")

	s.False(ulids.IsZero(ulids.MustParse(out.Id)))
	s.Equal(topic.ProjectId, out.ProjectId)
	s.Equal(topic.Name, out.Name)
	s.NotEmpty(out.Created)
	s.NotEmpty(out.Modified)

	// Should not be able to create a topic without a name
	topic.Name = ""
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing name field")
}
