package ensign_test

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverTestSuite) TestListTopics() {
	require := s.Require()
	s.store.UseError(store.ListTopics, errors.ErrIterReleased)
	defer s.store.Reset()

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
	}

	// Should not be able to create a topic when not authenticated
	_, err := s.client.ListTopics(context.Background(), &api.PageInfo{})
	s.GRPCErrorIs(err, codes.Unauthenticated, "missing credentials")

	// Should not be able to create a topic without the read topic permissions
	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// ProjectID is required in the claims
	claims.Permissions = []string{permissions.ReadTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// ProjectID must be valid and parseable
	claims.ProjectID = "foo"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Empty results should be returned on project not found
	claims.ProjectID = "01GV6G705RV812J20S6RKJHVGE"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	s.store.OnListTopics = func(ulid.ULID) iterator.TopicIterator {
		return store.NewTopicIterator(nil)
	}

	out, err := s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	require.NoError(err, "could not make a happy path request")
	require.Empty(out.NextPageToken, "expected no next page token on no results")
	require.Empty(out.Topics, "expected no topics on empty page request")

	// Results should be returned on project found
	s.store.UseFixture(store.ListTopics, "testdata/topics.json")

	out, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	require.NoError(err, "could not make a happy path request")
	require.Empty(out.NextPageToken, "expected no next page token on no results")
	require.Len(out.Topics, 4, "expected 3 topics returned")

	// TODO: test pagination
}

func (s *serverTestSuite) TestCreateTopic() {
	require := s.Require()
	s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)
	defer s.store.Reset()

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

	// Should not be able to create a topic without a project in the claims
	claims.Permissions = []string{permissions.CreateTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Should not be able to create a topic an invalid project in the claims
	claims.ProjectID = "invalidprojectid"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Should not be able to create a topic in the wrong project
	claims.ProjectID = "01GQFQCFC9P3S7QZTPYFVBJD7F"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Happy path: should be able to create a valid topic
	claims.ProjectID = "01GQ7P8DNR9MR64RJR9D64FFNT"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	s.store.OnCreateTopic = func(topic *api.Topic) error {
		if err := meta.ValidateTopic(topic, true); err != nil {
			return err
		}

		topic.Id = ulids.New().Bytes()
		topic.Created = timestamppb.Now()
		topic.Modified = topic.Created
		return nil
	}

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

	// Should not be able to create a topic without a project
	topic.Name = "testing.testapp.test"
	topic.ProjectId = nil
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing project id field")

	// Should not be able to create a topic without a valid projectID
	topic.ProjectId = []byte{118, 42}
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "invalid project id field")

	// Unhandled database error should create an internal error
	topic = &api.Topic{
		ProjectId: ulids.MustBytes("01GQ7P8DNR9MR64RJR9D64FFNT"),
		Name:      "testing.testapp.test",
	}
	s.store.UseError(store.CreateTopic, fmt.Errorf("something really bad happened"))
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Internal, "could not process create topic request")
}

func (s *serverTestSuite) TestDeleteTopic() {
	// Test common functionality for delete topic operations
	require := s.Require()
	defer s.store.Reset()

	for _, operation := range []api.TopicMod_Operation{api.TopicMod_ARCHIVE, api.TopicMod_DESTROY} {
		s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)
		claims := &tokens.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
			},
			OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
		}

		request := &api.TopicMod{
			Operation: operation,
		}

		token, err := s.quarterdeck.CreateAccessToken(claims)
		require.NoError(err, "could not create access token for request")

		// Should not be able to delete a topic when not authenticated
		_, err = s.client.DeleteTopic(context.Background(), request)
		s.GRPCErrorIs(err, codes.Unauthenticated, "missing credentials")

		// Should not be able to delete a topic without the correct permissions
		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

		// Should not be able to delete a topic without an ID
		claims.Permissions = []string{permissions.EditTopics, permissions.DestroyTopics}
		token, err = s.quarterdeck.CreateAccessToken(claims)
		require.NoError(err, "could not create access token for request")

		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.InvalidArgument, "missing id field")

		// Should not be able to delete a topic with an invalid ID
		request.Id = "foo"
		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.InvalidArgument, "invalid id field")

		// Should not be able to delete a topic that does not exist
		request.Id = ulids.New().String()
		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.NotFound, "topic not found")

		// Unhandled database exceptions should return an internal error
		s.store.UseError(store.RetrieveTopic, fmt.Errorf("somehing very bad happened"))
		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.Internal, "could not process delete topic request")

		// Able to retrieve the fixture to delete from this point on.
		err = s.store.UseFixture(store.RetrieveTopic, "testdata/topic.json")
		require.NoError(err, "could not load topic fixture")

		// Should receive a topic not found if the claims do not have a projectID
		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.NotFound, "topic not found")

		// Should receive an not found error if the claims projectID does not match the topic
		claims.ProjectID = ulids.New().String()
		token, err = s.quarterdeck.CreateAccessToken(claims)
		require.NoError(err, "could not create access token for request")

		_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		s.GRPCErrorIs(err, codes.NotFound, "topic not found")
	}
}

func (s *serverTestSuite) TestDeleteTopic_NOOP() {
	s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)
	defer s.store.Reset()

	request := &api.TopicMod{
		Operation: api.TopicMod_NOOP,
	}

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
	}

	require := s.Require()
	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	// Should not be able to delete a topic without a correct operation
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "invalid operation field")

	_, err = s.client.DeleteTopic(context.Background(), &api.TopicMod{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "invalid operation field")
}

func (s *serverTestSuite) TestDeleteTopic_Archive() {
	require := s.Require()
	err := s.store.UseFixture(store.RetrieveTopic, "testdata/topic.json")
	require.NoError(err, "could not load topic fixture")
	defer s.store.Reset()

	request := &api.TopicMod{
		Id:        "01GTSMQ3V8ASAPNCFEN378T8RD",
		Operation: api.TopicMod_ARCHIVE,
	}

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID:     "01GKHJRF01YXHZ51YMMKV3RCMK",
		ProjectID: "01GTSMMC152Q95RD4TNYDFJGHT",
	}

	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	// Should not be able to archive a topic without the topics:edit permission
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Should not be able to archive a topic with the topics:destroy permission
	claims.Permissions = []string{permissions.DestroyTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Happy path: should be able to mark topic as read-only
	s.store.OnUpdateTopic = func(topic *api.Topic) error {
		if !topic.Readonly {
			return fmt.Errorf("expected topic to be readonly")
		}
		return nil
	}

	claims.Permissions = []string{permissions.EditTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	out, err := s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	require.NoError(err, "could not execute happy path request")
	require.Equal(1, s.store.Calls(store.UpdateTopic))
	require.Zero(s.store.Calls(store.DeleteTopic))
	require.Equal(out.State, api.TopicTombstone_READONLY)
	require.Equal("01GTSMQ3V8ASAPNCFEN378T8RD", out.Id)
}

func (s *serverTestSuite) TestDeleteTopic_Destroy() {
	// Topic Destroy Tests
	require := s.Require()
	err := s.store.UseFixture(store.RetrieveTopic, "testdata/topic.json")
	require.NoError(err, "could not load topic fixture")
	defer s.store.Reset()

	request := &api.TopicMod{
		Id:        "01GTSMQ3V8ASAPNCFEN378T8RD",
		Operation: api.TopicMod_DESTROY,
	}

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID:     "01GKHJRF01YXHZ51YMMKV3RCMK",
		ProjectID: "01GTSMMC152Q95RD4TNYDFJGHT",
	}

	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	// Should not be able to destroy a topic without the topics:delete permission
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Should not be able to destroy a topic with the topics:edit permission
	claims.Permissions = []string{permissions.EditTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// Happy path: should be able to mark topic as read-only
	s.store.OnDeleteTopic = func(id ulid.ULID) error {
		if id.String() != "01GTSMQ3V8ASAPNCFEN378T8RD" {
			return errors.ErrNotFound
		}
		return nil
	}

	claims.Permissions = []string{permissions.DestroyTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	out, err := s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	require.NoError(err, "could not execute happy path request")
	require.Equal(1, s.store.Calls(store.DeleteTopic))
	require.Zero(s.store.Calls(store.UpdateTopic))
	require.Equal(out.State, api.TopicTombstone_DELETING)
	require.Equal("01GTSMQ3V8ASAPNCFEN378T8RD", out.Id)
}
