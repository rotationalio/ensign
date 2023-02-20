package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	en "github.com/rotationalio/ensign/pkg/api/v1beta1"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *tenantTestSuite) TestProjectTopicList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88")

	topics := []*db.Topic{
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399DWFK3E94FV30WF7QMJ5"),
			Name:      "topic001",
			Created:   time.Unix(1672161102, 0),
			Modified:  time.Unix(1672161102, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399KP7ZYFBHMD565EQBQQ4"),
			Name:      "topic002",
			Created:   time.Unix(1673659941, 0),
			Modified:  time.Unix(1673659941, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399RREX32HRT1YA0YEW4JW"),
			Name:      "topic003",
			Created:   time.Unix(1674073941, 0),
			Modified:  time.Unix(1674073941, 0),
		},
	}

	prefix := projectID[:]
	namespace := "topics"

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i, topic := range topics {
			data, err := topic.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     data,
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	// Set the initial claims fixture.
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated.
	_, err := suite.client.ProjectTopicList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// User must have the correct permissions.
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the project does not exist.
	_, err = suite.client.ProjectTopicList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusBadRequest, "could not parse project ulid", "expected error when project does not exist")

	rep, err := suite.client.ProjectTopicList(ctx, projectID.String(), &api.PageQuery{})
	require.NoError(err, "could not list project topics")
	require.Len(rep.Topics, 3, "expected 3 topics")

	// Verify topic data has been populated.
	for i := range topics {
		require.Equal(topics[i].ID.String(), rep.Topics[i].ID, "expected topic id to match")
		require.Equal(topics[i].Name, rep.Topics[i].Name, "expected topic name to match")
		require.Equal(topics[i].Created.Format(time.RFC3339Nano), rep.Topics[i].Created, "expected topic created to match")
		require.Equal(topics[i].Modified.Format(time.RFC3339Nano), rep.Topics[i].Modified, "expected topic modified to match")
	}
}

func (suite *tenantTestSuite) TestProjectTopicCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := ulids.New().String()
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		Name:     "project001",
		Created:  time.Now().Add(-time.Hour),
		Modified: time.Now(),
	}

	var data []byte
	data, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	// Call trtl OnGet method
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	reply := &qd.LoginReply{
		AccessToken: "token",
	}

	// Connect to Quarterdeck mock.
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	enTopic := &en.Topic{
		ProjectId: project.ID[:],
		Name:      "topic01",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	// Connect to Ensign mock.
	suite.ensign.OnCreateTopic = func(ctx context.Context, t *en.Topic) (*en.Topic, error) {
		return enTopic, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"create:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set client csrf protection")
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: "topic-example"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: "topic-example"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.CreateTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if project id is not a valid ULID.
	_, err = suite.client.ProjectTopicCreate(ctx, "projectID", &api.Topic{ID: "", Name: "topic-example"})
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project id does not exist")

	// Should return an error if topic id exists.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6", Name: "topic-example"})
	suite.requireError(err, http.StatusBadRequest, "topic id cannot be specified on create", "expected error when topic id exists")

	// Should return an error if topic name does not exist.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name does not exist")

	req := &api.Topic{
		Name: "topic001",
	}

	topic, err := suite.client.ProjectTopicCreate(ctx, projectID, req)
	require.NoError(err, "could not add topic")
	require.NotEmpty(topic.ID, "expected non-zero ulid to be populated")
	require.Equal(req.Name, topic.Name, "expected topic name to match")
	require.NotEmpty(topic.Created, "expected created to be populated")
	require.NotEmpty(topic.Modified, "expected modified to be populated")

	// Create a test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "0000000000000000",
		Permissions: []string{perms.CreateTopics},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{})
	suite.requireError(err, http.StatusInternalServerError, "could not parse org id", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestTopicList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88")

	defer cancel()

	topics := []*db.Topic{
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399DWFK3E94FV30WF7QMJ5"),
			Name:      "topic001",
			Created:   time.Unix(1672161102, 0),
			Modified:  time.Unix(1672161102, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399KP7ZYFBHMD565EQBQQ4"),
			Name:      "topic002",
			Created:   time.Unix(1673659941, 0),
			Modified:  time.Unix(1673659941, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399RREX32HRT1YA0YEW4JW"),
			Name:      "topic003",
			Created:   time.Unix(1674073941, 0),
			Modified:  time.Unix(1674073941, 0),
		},
	}

	prefix := projectID[:]
	namespace := "topics"

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate.
		for i, topic := range topics {
			data, err := topic.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     data,
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GNA91N6WMCWNG9MVSK47ZS88",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := suite.client.TopicList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	rep, err := suite.client.TopicList(ctx, &api.PageQuery{})
	require.NoError(err, "could not list topics")
	require.Len(rep.Topics, 3, "expected 3 topics")

	// Verify topic data has been populated.
	for i := range topics {
		require.Equal(topics[i].ID.String(), rep.Topics[i].ID, "expected topic id to match")
		require.Equal(topics[i].Name, rep.Topics[i].Name, "expected topic name to match")
		require.Equal(topics[i].Created.Format(time.RFC3339Nano), rep.Topics[i].Created, "expected topic created to match")
		require.Equal(topics[i].Modified.Format(time.RFC3339Nano), rep.Topics[i].Modified, "expected topic modified to match")
	}

	// Set test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "0000000000000000",
		Permissions: []string{perms.ReadTopics},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.TopicList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusInternalServerError, "could not parse org id", "expected error when org id is missing or not a valid ulid")

}

func (suite *tenantTestSuite) TestTopicDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	topic := &db.Topic{
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
	}

	// Marshal the topic data with msgpack.
	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")

	// Unmarshal the topic data with msgpack.
	other := &db.Topic{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the topic data")

	// Call OnGet method and return a GetReply.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.TopicDetail(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDetail(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the topic does not exist.
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic does not exist")

	// Create a topic test fixture.
	req := &api.Topic{
		ID:   "01GNA926JCTKDH3VZBTJM8MAF6",
		Name: "topic001",
	}

	rep, err := suite.client.TopicDetail(ctx, req.ID)
	require.NoError(err, "could not retrieve topic")
	require.Equal(req.ID, rep.ID, "expected topic ID to match")
	require.Equal(req.Name, rep.Name, "expected topic name to match")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, errors.New("key not found")
	}

	// Should return an error if the topic ID is parsed but not found.
	_, err = suite.client.TopicDetail(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusNotFound, "could not retrieve topic", "expected error when topic ID is not found")
}

func (suite *tenantTestSuite) TestTopicUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	topic := &db.Topic{
		OrgID:     ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
	}

	// Marshal the topic data with msgpack.
	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")

	// Unmarshal the topic data with msgpack.
	other := &db.Topic{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the topic data")

	// OnGet method should return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// OnPut method should return a success response.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set client csrf protection")
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the topic is not parseable.
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "invalid"})
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic is not parseable")

	// Should return an error if the topic name is missing.
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6"})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name is missing")

	// Create a topic test fixture.
	req := &api.Topic{
		ID:   "01GNA926JCTKDH3VZBTJM8MAF6",
		Name: "topic001",
	}

	rep, err := suite.client.TopicUpdate(ctx, req)
	require.NoError(err, "could not update topic")
	require.NotEqual(req.ID, "", "topic id should not match")
	require.Equal(req.Name, rep.Name, "expected topic name to match")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, errors.New("key not found")
	}

	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6", Name: "topic001"})
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic ID is not found")
}

func (suite *tenantTestSuite) TestTopicDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	topicID := "01GNA926JCTKDH3VZBTJM8MAF6"
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set client csrf protection")
	err := suite.client.TopicDelete(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.TopicDelete(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.DestroyTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the topic does not exist.
	err = suite.client.TopicDelete(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic does not exist")

	err = suite.client.TopicDelete(ctx, topicID)
	require.NoError(err, "could not delete topic")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return nil, errors.New("key not found")
	}

	err = suite.client.TopicDelete(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusNotFound, "could not delete topic", "expected error when topic ID is not found")
}
