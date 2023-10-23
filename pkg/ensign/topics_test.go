package ensign_test

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverTestSuite) TestListTopics() {
	require := s.Require()
	s.store.UseError(store.ListTopics, errors.ErrIterReleased)

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
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// ProjectID is required in the claims
	claims.Permissions = []string{permissions.ReadTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// ProjectID must be valid and parseable
	claims.ProjectID = "foo"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

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
	err = s.store.UseFixture(store.ListTopics, "testdata/topics.json")
	require.NoError(err, "could not load testdata/topics.json")

	out, err = s.client.ListTopics(context.Background(), &api.PageInfo{}, mock.PerRPCToken(token))
	require.NoError(err, "could not make a happy path request")
	require.Empty(out.NextPageToken, "expected no next page token on no results")
	require.Len(out.Topics, 5, "expected 5 topics returned")

	// TODO: test pagination
}

func (s *serverTestSuite) TestCreateTopic() {
	require := s.Require()
	s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)

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
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Should not be able to create a topic without a project in the claims
	claims.Permissions = []string{permissions.CreateTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Should not be able to create a topic without a project
	topic.ProjectId = nil
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing project id field")

	// Should not be able to create a topic an invalid project in the claims
	topic.ProjectId = ulids.New().Bytes()
	claims.ProjectID = "invalidprojectid"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Should not be able to create a topic in the wrong project
	claims.ProjectID = "01GQFQCFC9P3S7QZTPYFVBJD7F"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	claims.ProjectID = "01GQ7P8DNR9MR64RJR9D64FFNT"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	// Happy path: should be able to create a valid topic without projectID on topic
	// Because the projectID is set by the claims and should also be able to create a
	// topic with a projectID that matches the ones in the claims.
	for _, projectID := range [][]byte{nil, ulids.MustParse(claims.ProjectID).Bytes()} {
		topic.ProjectId = projectID
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
		require.Equal(ulids.MustParse(claims.ProjectID).Bytes(), out.ProjectId)
		require.Equal(topic.Name, out.Name)
		require.Equal(api.TopicState_READY, out.Status)
		require.NotEmpty(out.Created)
		require.NotEmpty(out.Modified)
	}

	// Should not be able to create a topic without a name
	topic.Name = ""
	_, err = s.client.CreateTopic(context.Background(), topic, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing name field")

	// Should not be able to create a topic without a valid projectID
	topic.ProjectId = []byte{118, 42}
	topic.Name = "testing.testapp.test"
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

func (s *serverTestSuite) TestRetrieveTopic() {
	require := s.Require()

	// Should not be able to retrieve a topic when not authenticated
	request := &api.Topic{Id: ulids.New().Bytes()}
	_, err := s.client.RetrieveTopic(context.Background(), request)
	s.GRPCErrorIs(err, codes.Unauthenticated, "missing credentials")

	// Should not be able to retrieve a topic with no project ID
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
	}

	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	// Should not be able to delete a topic without the correct permissions
	_, err = s.client.RetrieveTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	claims.Permissions = []string{permissions.ReadTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	// Should not be able to delete a topic without a topic ID
	_, err = s.client.RetrieveTopic(context.Background(), &api.Topic{}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "missing topic id field")

	// Should not be able to delete a topic with an invalid ID
	_, err = s.client.RetrieveTopic(context.Background(), &api.Topic{Id: []byte("foo")}, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "invalid topic id field")

	// Should receive a not found error if the topic doesn't exist
	s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)
	_, err = s.client.RetrieveTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.NotFound, "topic not found")

	// Should return an internal error if database is down
	s.store.UseError(store.RetrieveTopic, fmt.Errorf("something big exploded"))
	_, err = s.client.RetrieveTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Internal, "could not complete retrieve topic request")

	// At this point the topic should be returned from the database
	err = s.store.UseFixture(store.RetrieveTopic, "testdata/topic.json")
	require.NoError(err, "could not load test fixture")

	// Should not be able to retrieve topic if no project is on the claims
	_, err = s.client.RetrieveTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.NotFound, "topic not found")

	// Should be able to successfully retrieve a topic
	claims.ProjectID = "01GTSMMC152Q95RD4TNYDFJGHT"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	topic, err := s.client.RetrieveTopic(context.Background(), request, mock.PerRPCToken(token))
	require.NoError(err, "could not retrieve topic")

	require.Equal("testing.testapp.test", topic.Name)
}

func (s *serverTestSuite) TestDeleteTopic() {
	// Test common functionality for delete topic operations
	require := s.Require()

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
		s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

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
		s.store.UseError(store.RetrieveTopic, fmt.Errorf("something very bad happened"))
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

		// Should be able to successfully delete a topic
		claims.ProjectID = "01GTSMMC152Q95RD4TNYDFJGHT"
		token, err = s.quarterdeck.CreateAccessToken(claims)
		require.NoError(err, "could not create access token for request")

		s.store.UseError(store.DeleteTopic, nil)
		s.store.UseError(store.UpdateTopic, nil)

		// See operation-specific delete tests for more thorough assertions.
		rep, err := s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
		require.NoError(err, "could not delete topic")
		require.Equal(request.Id, rep.Id)
		require.True(rep.State == api.TopicState_DELETING || rep.State == api.TopicState_READONLY)
	}
}

func (s *serverTestSuite) TestDeleteTopicState() {
	require := s.Require()

	testCases := []struct {
		state     api.TopicState
		operation api.TopicMod_Operation
		err       error
	}{
		{api.TopicState_READY, api.TopicMod_ARCHIVE, nil},
		{api.TopicState_READY, api.TopicMod_DESTROY, nil},
		{api.TopicState_READONLY, api.TopicMod_ARCHIVE, nil},
		{api.TopicState_READONLY, api.TopicMod_DESTROY, nil},
		{api.TopicState_UNDEFINED, api.TopicMod_ARCHIVE, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_UNDEFINED, api.TopicMod_DESTROY, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_DELETING, api.TopicMod_ARCHIVE, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_DELETING, api.TopicMod_DESTROY, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_PENDING, api.TopicMod_ARCHIVE, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_PENDING, api.TopicMod_DESTROY, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_ALLOCATING, api.TopicMod_ARCHIVE, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_ALLOCATING, api.TopicMod_DESTROY, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_REPAIRING, api.TopicMod_ARCHIVE, status.Error(codes.FailedPrecondition, "--")},
		{api.TopicState_REPAIRING, api.TopicMod_DESTROY, status.Error(codes.FailedPrecondition, "--")},
	}

	ctx := context.Background()
	topicID := ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD")

	// Authorize access
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID:       "01GKHJRF01YXHZ51YMMKV3RCMK",
		ProjectID:   "01GTSMMC152Q95RD4TNYDFJGHT",
		Permissions: []string{"topics:edit", "topics:destroy"},
	}

	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")
	s.store.UseError(store.UpdateTopic, nil)

	for i, tc := range testCases {
		s.store.OnRetrieveTopic = func(topicID ulid.ULID) (*api.Topic, error) {
			return &api.Topic{
				Id:        topicID[:],
				ProjectId: ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
				Status:    tc.state,
			}, nil
		}

		_, err := s.client.DeleteTopic(ctx, &api.TopicMod{Id: topicID.String(), Operation: tc.operation}, mock.PerRPCToken(token))
		if tc.err != nil {
			require.Error(err, "expected an error on test case %d", i)
			require.Equal(status.Code(tc.err), status.Code(err), "expected failed precondition on test case %d")
		} else {
			require.NoError(err, "expected no error on test case %d", i)
		}
	}
}

func (s *serverTestSuite) TestDeleteTopic_NOOP() {
	s.store.UseError(store.RetrieveTopic, errors.ErrNotFound)

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
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Should not be able to archive a topic with the topics:destroy permission
	claims.Permissions = []string{permissions.DestroyTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

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
	require.Equal(out.State, api.TopicState_READONLY)
	require.Equal("01GTSMQ3V8ASAPNCFEN378T8RD", out.Id)
}

func (s *serverTestSuite) TestDeleteTopic_Destroy() {
	// Topic Destroy Tests
	require := s.Require()
	err := s.store.UseFixture(store.RetrieveTopic, "testdata/topic.json")
	require.NoError(err, "could not load topic fixture")

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
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Should not be able to destroy a topic with the topics:edit permission
	claims.Permissions = []string{permissions.EditTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")
	_, err = s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Happy path: should be able to mark topic as read-only
	s.store.OnDeleteTopic = func(id ulid.ULID) error {
		if id.String() != "01GTSMQ3V8ASAPNCFEN378T8RD" {
			return errors.ErrNotFound
		}
		return nil
	}
	s.store.OnDestroy = func(ulid.ULID) error { return nil }
	s.store.OnUpdateTopic = func(*api.Topic) error { return nil }

	claims.Permissions = []string{permissions.DestroyTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create access token for request")

	out, err := s.client.DeleteTopic(context.Background(), request, mock.PerRPCToken(token))
	require.NoError(err, "could not execute happy path request")
	require.Equal(1, s.store.Calls(store.UpdateTopic))
	// TODO: test that destroy and delete topic are queued in the tasks.
	// require.Equal(1, s.store.Calls(store.Destroy))
	// require.Equal(1, s.store.Calls(store.DeleteTopic))
	require.Equal(out.State, api.TopicState_DELETING)
	require.Equal("01GTSMQ3V8ASAPNCFEN378T8RD", out.Id)
}
