package ensign_test

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
)

func (s *serverTestSuite) TestCreateTopic() {
	require := s.Require()

	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GQ7P8DNR9MR64RJR9D64FFNT"),
		Name:      "testing.testapp.test",
	}

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
	}

	// Should not be able to create a topic when not authenticated
	_, err := s.client.CreateTopic(context.Background(), topic)
	s.GRPCErrorIs(err, codes.Unauthenticated, "missing credentials")

	// Should not be able to create a topic without the create topic permissions
	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Should not be able to create a topic without a project
	claims.Permissions = []string{permissions.CreateTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "missing project id")

	// Should not be able to create a topic an invalid project
	claims.ProjectID = "invalidprojectid"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "invalid project id")

	// Should not be able to create a topic in the wrong project
	claims.ProjectID = "01GQFQCFC9P3S7QZTPYFVBJD7F"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Happy path: should be able to create a valid topic
	claims.ProjectID = "01GTW1R9MH8723JQDRMFE16CZ7"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	out, err := s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	require.NoError(err, "could not execute create topic request")

	require.False(ulids.IsZero(ulids.MustParse(out.Id)))
	require.Equal(topic.ProjectId, out.ProjectId)
	require.Equal(topic.Name, out.Name)
	require.NotEmpty(out.Created)
	require.NotEmpty(out.Modified)

	// Should not be able to create a topic without a name
	topic.Name = ""
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing name field")
}
