package tenant_test

import (
	"bytes"
	"context"
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
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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

	orgID := ulids.New()
	key, err := db.CreateKey(orgID, projectID)
	require.NoError(err, "could not create project key")

	keyData, err := key.MarshalValue()
	require.NoError(err, "could not marshal key data")

	// OnGet method should return the project key
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if !bytes.Equal(in.Key, projectID[:]) || in.Namespace != db.KeysNamespace {
			return nil, status.Error(codes.FailedPrecondition, "unexpected get request")
		}
		return &pb.GetReply{
			Value: keyData,
		}, nil
	}

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
	_, err = suite.client.ProjectTopicList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the user.
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// TODO: Add test for wrong orgID in claims

	// Should return an error if the project ID is not parseable.
	claims.OrgID = orgID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
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

	key, err := project.Key()
	require.NoError(err, "could not create project key")

	var data []byte
	data, err = project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	// Trtl Get should return project key or project data
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %s not found", in.Namespace)
		}
	}

	reply := &qd.LoginReply{
		AccessToken: "token",
	}

	// Connect to Quarterdeck mock.
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	enTopic := &en.Topic{
		ProjectId: project.ID[:],
		Id:        ulids.New().Bytes(),
		Name:      "topic01",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	// Connect to Ensign mock.
	suite.ensign.OnCreateTopic = func(ctx context.Context, t *en.Topic) (*en.Topic, error) {
		return enTopic, nil
	}

	// Call OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
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
	suite.requireError(err, http.StatusBadRequest, "could not parse project id from url", "expected error when project id is not a valid ULID")

	// Should return an error if topic id exists.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6", Name: "topic-example"})
	suite.requireError(err, http.StatusBadRequest, "topic id cannot be specified on create", "expected error when topic id exists")

	// Should return an error if topic name does not exist.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name does not exist")

	// Should return an error if org ID is not in the claims.
	claims.OrgID = ""
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: "topic-example"})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org ID is not in the claims")

	// TODO: Add test for wrong orgID in claims

	// Reset claims org ID for tests.
	claims.OrgID = project.OrgID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	req := &api.Topic{
		Name: enTopic.Name,
	}

	topic, err := suite.client.ProjectTopicCreate(ctx, projectID, req)
	require.NoError(err, "could not add topic")
	require.NotEmpty(topic.ID, "expected non-zero ulid to be populated")
	require.Equal(req.Name, topic.Name, "expected topic name to match")
	require.NotEmpty(topic.Created, "expected created to be populated")
	require.NotEmpty(topic.Modified, "expected modified to be populated")

	// Should return an error if Quarterdeck returns an error.
	suite.quarterdeck.OnProjects(mock.UseError(http.StatusBadRequest, "missing field project_id"), mock.RequireAuth())
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, req)
	suite.requireError(err, http.StatusBadRequest, "missing field project_id", "expected error when Quarterdeck returns an error")

	// Should return an error if Ensign returns an error.
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply), mock.RequireAuth())
	suite.ensign.OnCreateTopic = func(ctx context.Context, t *en.Topic) (*en.Topic, error) {
		return &en.Topic{}, status.Error(codes.Internal, "could not create topic")
	}
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, req)
	suite.requireError(err, http.StatusInternalServerError, "could not create topic", "expected error when Ensign returns an error")
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
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

}

func (suite *tenantTestSuite) TestTopicDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	id := "01GNA926JCTKDH3VZBTJM8MAF6"
	project := "01GNA91N6WMCWNG9MVSK47ZS88"
	org := "01GNA91N6WMCWNG9MVSK47ZS88"
	topic := &db.Topic{
		OrgID:     ulid.MustParse(org),
		ProjectID: ulid.MustParse(project),
		ID:        ulid.MustParse(id),
		Name:      "topic001",
		Created:   time.Now().Add(-time.Hour),
		Modified:  time.Now(),
	}

	key, err := topic.Key()
	require.NoError(err, "could not get topic key")

	// Marshal the topic data with msgpack.
	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")

	// Trtl Get should return the topic key or the topic data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.TopicNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %s not found", gr.Namespace)
		}
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the topic id is not parseable
	claims.OrgID = ulids.New().String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic does not exist")

	// TODO: Add test for wrong orgID in claims

	claims.OrgID = org
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	rep, err := suite.client.TopicDetail(ctx, id)
	require.NoError(err, "could not retrieve topic")
	require.Equal(topic.ID.String(), rep.ID, "expected topic ID to match")
	require.Equal(topic.Name, rep.Name, "expected topic name to match")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Errorf(codes.NotFound, "key not found")
	}

	// Should return an error if the topic ID is parsed but not found.
	_, err = suite.client.TopicDetail(ctx, id)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic ID is not found")
}

func (suite *tenantTestSuite) TestTopicUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	id := "01GNA926JCTKDH3VZBTJM8MAF6"
	orgID := "01GNA91N6WMCWNG9MVSK47ZS88"
	projectID := "01GNA91N6WMCWNG9MVSK47ZS88"
	topic := &db.Topic{
		OrgID:     ulid.MustParse(orgID),
		ProjectID: ulid.MustParse(projectID),
		ID:        ulid.MustParse(id),
		Name:      "topic001",
	}

	key, err := topic.Key()
	require.NoError(err, "could not create topic key")

	// Marshal the topic data with msgpack.
	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")

	// Trtl Get should return the topic key or the topic data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.TopicNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %s not found", gr.Namespace)
		}
	}

	// OnPut method should return a success response.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Configure Quarterdeck to return a success response on ProjectAccess requests.
	auth := &qd.LoginReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))

	// Configure Ensign to return a success response on DeleteTopic requests.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return &en.TopicTombstone{
			Id:    topic.ID.String(),
			State: en.TopicTombstone_READONLY,
		}, nil
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

	// Should return an error if the orgID is missing from the claims
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6"})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when orgID is missing")

	// Should return an error if the topic is not parseable.
	claims.OrgID = orgID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "invalid"})
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic is not parseable")

	// Should return an error if the topic name is missing.
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: id, ProjectID: projectID})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name is missing")

	// Should return an error if the topic name is invalid.
	req := &api.Topic{
		ID:        id,
		ProjectID: projectID,
		Name:      "New$Topic$Name",
		State:     en.TopicTombstone_UNKNOWN.String(),
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, db.ErrInvalidTopicName.Error(), "expected error when topic name is invalid")

	// Should return an error if the orgIDs do not match.
	req.Name = "NewTopicName"
	claims.OrgID = ulids.New().String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when orgIDs do not match")

	// Only update the name of a topic.
	claims.OrgID = orgID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	rep, err := suite.client.TopicUpdate(ctx, req)
	require.NoError(err, "could not update topic")
	require.Equal(topic.ID.String(), rep.ID, "expected topic ID to be unchanged")
	require.Equal(req.Name, rep.Name, "expected topic name to be updated")
	require.Equal(topic.State.String(), rep.State, "expected topic state to be unchanged")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	// Should return an error if the topic state is invalid
	req.State = en.TopicTombstone_DELETING.String()
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "topic state can only be set to READONLY", "expected error when topic state is invalid")

	// Should return an error if the topic is already being deleted.
	topic.State = en.TopicTombstone_DELETING
	data, err = topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")
	req.State = en.TopicTombstone_READONLY.String()
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "topic is already being deleted", "expected error when topic is already being deleted")

	// Sucessfully updating the topic state.
	topic.State = en.TopicTombstone_UNKNOWN
	data, err = topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")
	rep, err = suite.client.TopicUpdate(ctx, req)
	require.NoError(err, "could not update topic")
	require.Equal(topic.ID.String(), rep.ID, "expected topic ID to be unchanged")
	require.Equal(req.Name, rep.Name, "expected topic name to be updated")
	require.Equal(req.State, rep.State, "expected topic state to be updated")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}

	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic ID is not found")

	// Should return an error if Quarterdeck returns an error.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.TopicNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %q not found", gr.Namespace)
		}
	}
	suite.quarterdeck.OnProjects(mock.UseError(http.StatusInternalServerError, "could not get one time credentials"))
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not get one time credentials", "expected error when Quarterdeck returns an error")

	// Should return not found if Ensign returns not found.
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return nil, status.Error(codes.NotFound, "could not archive topic")
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when Ensign returns an error")

	// Should return an error if Ensign returns an error.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return nil, status.Error(codes.Internal, "could not archive topic")
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not update topic", "expected error when Ensign returns an error")
}

func (suite *tenantTestSuite) TestTopicDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	topicID := "01GNA926JCTKDH3VZBTJM8MAF6"
	orgID := "01GNA91N6WMCWNG9MVSK47ZS88"
	projectID := "02ABC91N6WMCWNG9MVSK47ZSYZ"
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	topic := &db.Topic{
		OrgID:     ulid.MustParse(orgID),
		ProjectID: ulid.MustParse(projectID),
		ID:        ulid.MustParse(topicID),
		Name:      "mytopic",
	}

	key, err := topic.Key()
	require.NoError(err, "could not create topic key")

	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal topic data")

	// Configure Trtl to return the topic key or the topic data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.TopicNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %q not found", gr.Namespace)
		}
	}

	// Configure Trtl to return a success response on Put requests.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Call OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	// Configure Quarterdeck to return a success response on ProjectAccess requests.
	auth := &qd.LoginReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))

	// Configure Ensign to return a success response on DeleteTopic requests.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return &en.TopicTombstone{
			Id:    topic.ID.String(),
			State: en.TopicTombstone_DELETING,
		}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set client csrf protection")
	req := &api.Confirmation{
		ID: topicID,
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.DestroyTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the orgID is not in the claims
	req.ID = "invalid"
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when orgID is not in claims")

	// Should return an error if the topic does not exist.
	claims.OrgID = orgID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic does not exist")

	// Should return an error if the orgIDs don't match
	req.ID = topicID
	claims.OrgID = ulids.New().String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when orgIDs don't match")

	// Retrieve a confirmation from the first successful request.
	claims.OrgID = orgID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	reply, err := suite.client.TopicDelete(ctx, req)
	require.NoError(err, "could not delete topic")
	require.Equal(reply.ID, topicID, "expected topic ID to match")
	require.Equal(reply.Name, topic.Name, "expected topic name to match")
	require.NotEmpty(reply.Token, "expected confirmation token to be set")

	// Simulate the backend saving the token in the database.
	topic.ConfirmDeleteToken = reply.Token
	data, err = topic.MarshalValue()
	require.NoError(err, "could not marshal topic data")

	// Should return an error if the token is invalid
	req.Token = "invalid"
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusPreconditionFailed, "invalid confirmation token", "expected error when token is invalid")

	// Should return an error if the token has expired.
	tokenData := &db.ResourceToken{
		ID:        ulids.New(),
		Secret:    reply.Token,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	token, err := tokenData.Create()
	require.NoError(err, "could not create string token from data")
	req.Token = token
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusPreconditionFailed, "invalid confirmation token", "expected error when token is expired")

	// Should return an error if the wrong token is provided.
	tokenData.ExpiresAt = time.Now().Add(1 * time.Minute)
	token, err = tokenData.Create()
	require.NoError(err, "could not create string token from data")
	req.Token = token
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusPreconditionFailed, "invalid confirmation token", "expected error when wrong token is provided")

	// Successfully requesting the topic delete
	req.Token = reply.Token
	expected := &api.Confirmation{
		ID:     topicID,
		Name:   topic.Name,
		Token:  reply.Token,
		Status: en.TopicTombstone_DELETING.String(),
	}
	reply, err = suite.client.TopicDelete(ctx, req)
	require.NoError(err, "could not delete topic")
	require.Equal(expected, reply, "expected confirmation reply to match")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Error(codes.NotFound, "key not found")
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic ID is not found")

	// Should return an error if Quarterdeck returns an error.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.TopicNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %q not found", gr.Namespace)
		}
	}
	suite.quarterdeck.OnProjects(mock.UseError(http.StatusInternalServerError, "could not get one time credentials"))
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not get one time credentials", "expected error when Quarterdeck returns an error")

	// Should return not found if Ensign returns not found.
	suite.quarterdeck.OnProjects(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return nil, status.Error(codes.NotFound, "could not delete topic")
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when Ensign returns an error")

	// Should return an error if Ensign returns an error.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *en.TopicMod) (*en.TopicTombstone, error) {
		return nil, status.Error(codes.Internal, "could not delete topic")
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not delete topic", "expected error when Ensign returns an error")
}
