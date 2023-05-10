package tenant_test

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	tk "github.com/rotationalio/ensign/pkg/utils/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	sdk "github.com/rotationalio/go-ensign/api/v1beta1"
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
			State:     sdk.TopicTombstone_DELETING,
			Created:   time.Unix(1672161102, 0),
			Modified:  time.Unix(1672161102, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399KP7ZYFBHMD565EQBQQ4"),
			Name:      "topic002",
			State:     sdk.TopicTombstone_READONLY,
			Created:   time.Unix(1673659941, 0),
			Modified:  time.Unix(1673659941, 0),
		},
		{
			ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
			ID:        ulid.MustParse("01GQ399RREX32HRT1YA0YEW4JW"),
			Name:      "topic003",
			State:     sdk.TopicTombstone_UNKNOWN,
			Created:   time.Unix(1674073941, 0),
			Modified:  time.Unix(1674073941, 0),
		},
	}

	expected := []*api.Topic{
		{
			ID:       topics[0].ID.String(),
			Name:     topics[0].Name,
			Status:   db.TopicStatusDeleting,
			Created:  topics[0].Created.Format(time.RFC3339Nano),
			Modified: topics[0].Modified.Format(time.RFC3339Nano),
		},
		{
			ID:       topics[1].ID.String(),
			Name:     topics[1].Name,
			Status:   db.TopicStatusArchived,
			Created:  topics[1].Created.Format(time.RFC3339Nano),
			Modified: topics[1].Modified.Format(time.RFC3339Nano),
		},
		{
			ID:       topics[2].ID.String(),
			Name:     topics[2].Name,
			Status:   db.TopicStatusActive,
			Created:  topics[2].Created.Format(time.RFC3339Nano),
			Modified: topics[2].Modified.Format(time.RFC3339Nano),
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
		switch in.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: keyData,
			}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: projectID[:],
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Call the OnCursor method
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, topic := range topics {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, topic.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := topic.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       topic.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	req := &api.PageQuery{}

	// Set the initial claims fixture.
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated.
	_, err = suite.client.ProjectTopicList(ctx, "invalid", req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicList(ctx, "invalid", req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the user.
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicList(ctx, projectID.String(), &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the project ID is not parseable.
	claims.OrgID = projectID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicList(ctx, "invalid", req)
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project does not exist")

	rep, err := suite.client.ProjectTopicList(ctx, projectID.String(), req)
	require.NoError(err, "could not list project topics")
	require.Len(rep.Topics, 3, "expected 3 topics")
	require.Empty(rep.NextPageToken, "did not expect next page token since there is only 1 page")

	// Verify topic data has been populated.
	for i := range topics {
		require.Equal(expected[i].ID, rep.Topics[i].ID, "expected topic id to match")
		require.Equal(expected[i].Name, rep.Topics[i].Name, "expected topic name to match")
		require.Equal(expected[i].Status, rep.Topics[i].Status, "expected topic status to match")
		require.Equal(expected[i].Created, rep.Topics[i].Created, "expected topic created to match")
		require.Equal(expected[i].Modified, rep.Topics[i].Modified, "expected topic modified to match")
	}

	// Set page size and test pagination.
	req.PageSize = 2
	rep, err = suite.client.ProjectTopicList(ctx, projectID.String(), req)
	require.NoError(err, "could not list topics")
	require.Len(rep.Topics, 2, "expected 2 topics")
	require.NotEmpty(rep.NextPageToken, "next page token should be set")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.ProjectTopicList(ctx, projectID.String(), req)
	require.NoError(err, "could not list topics")
	require.Len(rep2.Topics, 1, "expected 1 topic")
	require.NotEqual(rep.Topics[0].ID, rep2.Topics[0].ID, "should not have same topic ID")
	require.Empty(rep2.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 3, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 3; i++ {
		page, err := suite.client.ProjectTopicList(ctx, projectID.String(), req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Topics)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 2, "expected 3 results in 2 pages")
	require.Equal(nResults, 3, "expected 3 results in 2 pages")
}

func (suite *tenantTestSuite) TestProjectTopicCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := "01GNA91N6WMCWNG9MVSK47ZS88"
	defer cancel()
	defer suite.ResetTasks()

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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: project.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %s not found", in.Namespace)
		}
	}

	reply := &qd.LoginReply{
		AccessToken: "token",
	}

	// Connect to Quarterdeck mock.
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	enTopic := &sdk.Topic{
		ProjectId: project.ID[:],
		Id:        ulids.New().Bytes(),
		Name:      "topic01",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	// Connect to Ensign mock.
	suite.ensign.OnCreateTopic = func(ctx context.Context, t *sdk.Topic) (*sdk.Topic, error) {
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
		OrgID:       "01GNA91N6WMCWNG9MVSK47ZS88",
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
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project id is not a valid ULID")

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

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{Name: "topic-example"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Reset claims org ID for tests.
	claims.OrgID = project.ID.String()
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
	suite.quarterdeck.OnProjects("access", mock.UseError(http.StatusBadRequest, "missing field project_id"), mock.RequireAuth())
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, req)
	suite.requireError(err, http.StatusBadRequest, "missing field project_id", "expected error when Quarterdeck returns an error")

	// Should return an error if Ensign returns an error.
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply), mock.RequireAuth())
	suite.ensign.OnCreateTopic = func(ctx context.Context, t *sdk.Topic) (*sdk.Topic, error) {
		return &sdk.Topic{}, status.Error(codes.Internal, "could not create topic")
	}
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, req)
	suite.requireError(err, http.StatusInternalServerError, "could not create topic", "expected error when Ensign returns an error")

	// Ensure project stats update task finishes.
	suite.StopTasks()
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

	// Call the OnCursor method
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, topic := range topics {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, topic.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := topic.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       topic.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	req := &api.PageQuery{}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GNA91N6WMCWNG9MVSK47ZS88",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := suite.client.TopicList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadTopics}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Retrieve all topics.
	rep, err := suite.client.TopicList(ctx, req)
	require.NoError(err, "could not list topics")
	require.Len(rep.Topics, 3, "expected 3 topics")
	require.Empty(rep.NextPageToken, "did not expect next page token since there is only 1 page")

	// Verify topic data has been populated.
	for i := range topics {
		require.Equal(topics[i].ID.String(), rep.Topics[i].ID, "expected topic id to match")
		require.Equal(topics[i].Name, rep.Topics[i].Name, "expected topic name to match")
		require.Equal(topics[i].Created.Format(time.RFC3339Nano), rep.Topics[i].Created, "expected topic created to match")
		require.Equal(topics[i].Modified.Format(time.RFC3339Nano), rep.Topics[i].Modified, "expected topic modified to match")
	}

	// Set page size and test pagination.
	req.PageSize = 2
	rep, err = suite.client.TopicList(ctx, req)
	require.NoError(err, "could not list topics")
	require.Len(rep.Topics, 2, "expected 2 topics")
	require.NotEmpty(rep.NextPageToken, "next page token should be set")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.TopicList(ctx, req)
	require.NoError(err, "could not list topics")
	require.Len(rep2.Topics, 1, "expected 1 topic")
	require.NotEqual(rep.Topics[0].ID, rep2.Topics[0].ID, "should not have same topic ID")
	require.Empty(rep2.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 3, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 3; i++ {
		page, err := suite.client.TopicList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Topics)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 2, "expected 3 results in 2 pages")
	require.Equal(nResults, 3, "expected 3 results in 2 pages")

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
		State:     sdk.TopicTombstone_READONLY,
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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: topic.ID[:]}, nil
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

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDetail(ctx, id)
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the topic id is not parseable
	claims.OrgID = id
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic does not exist")

	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	rep, err := suite.client.TopicDetail(ctx, id)
	require.NoError(err, "could not retrieve topic")
	require.Equal(topic.ID.String(), rep.ID, "expected topic ID to match")
	require.Equal(topic.Name, rep.Name, "expected topic name to match")
	require.Equal(db.TopicStatusArchived, rep.Status, "expected topic state to match")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")
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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: topic.ID[:]}, nil
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
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))

	// Configure Ensign to return a success response on DeleteTopic requests.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return &sdk.TopicTombstone{
			Id:    topic.ID.String(),
			State: sdk.TopicTombstone_READONLY,
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

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: id, ProjectID: projectID, Name: "project01"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the topic is not parseable.
	claims.OrgID = id
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "invalid"})
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic is not parseable")

	// Should return an error if the topic name is missing.
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: id, ProjectID: projectID})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name is missing")

	// Should return an error if the topic name is invalid.
	req := &api.Topic{
		ID:        id,
		ProjectID: projectID,
		Name:      "New$Topic$Name",
		Status:    db.TopicStatusActive,
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, db.ErrInvalidTopicName.Error(), "expected error when topic name is invalid")

	// Only update the name of a topic.
	req.Name = "NewTopicName"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	rep, err := suite.client.TopicUpdate(ctx, req)
	require.NoError(err, "could not update topic")
	require.Equal(topic.ID.String(), rep.ID, "expected topic ID to be unchanged")
	require.Equal(req.Name, rep.Name, "expected topic name to be updated")
	require.Equal(req.Status, rep.Status, "expected topic state to be unchanged")
	require.NotEmpty(rep.Created, "expected topic created to be set")
	require.NotEmpty(rep.Modified, "expected topic modified to be set")

	// Should return an error if the topic state is invalid
	req.Status = db.TopicStatusDeleting
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "topic state can only be set to Archived", "expected error when topic state is invalid")

	// Should return an error if the topic is already being deleted.
	topic.State = sdk.TopicTombstone_DELETING
	data, err = topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")
	req.Status = db.TopicStatusArchived
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "topic is already being deleted", "expected error when topic is already being deleted")

	// Valid request to update the topic state.
	// TODO: Update this test when topic archive is implemented in the SDK.
	topic.State = sdk.TopicTombstone_UNKNOWN
	data, err = topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "archiving a topic is not supported")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: topic.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: topic.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %q not found", gr.Namespace)
		}
	}
	suite.quarterdeck.OnProjects("access", mock.UseError(http.StatusInternalServerError, "could not get one time credentials"))
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not get one time credentials", "expected error when Quarterdeck returns an error")

	// Should return not found if Ensign returns not found.
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return nil, status.Error(codes.NotFound, "could not archive topic")
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "archiving a topic is not supported", "expected error when Ensign returns an error")

	// Should return an error if Ensign returns an error.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return nil, status.Error(codes.Internal, "could not archive topic")
	}
	_, err = suite.client.TopicUpdate(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "archiving a topic is not supported", "expected error when Ensign returns an error")
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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: topic.ID[:]}, nil
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
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))

	// Configure Ensign to return a success response on DeleteTopic requests.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return &sdk.TopicTombstone{
			Id:    topic.ID.String(),
			State: sdk.TopicTombstone_DELETING,
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
	claims.OrgID = topicID
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotFound, "topic not found", "expected error when topic does not exist")

	// Retrieve a confirmation from the first successful request.
	req.ID = topicID
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
	tokenData := &tk.Confirmation{
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

	// Valid delete request
	// TODO: Update when the DestroyTopic is implemented in the Go SDK.
	req.Token = reply.Token
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "deleting a topic is not supported")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: topic.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
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
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: topic.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %q not found", gr.Namespace)
		}
	}
	suite.quarterdeck.OnProjects("access", mock.UseError(http.StatusInternalServerError, "could not get one time credentials"))
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not get one time credentials", "expected error when Quarterdeck returns an error")

	// Should return not found if Ensign returns not found.
	suite.quarterdeck.OnProjects("access", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(auth))
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return nil, status.Error(codes.NotFound, "could not delete topic")
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "deleting a topic is not supported", "expected error when Ensign returns an error")

	// Should return an error if Ensign returns an error.
	suite.ensign.OnDeleteTopic = func(ctx context.Context, req *sdk.TopicMod) (*sdk.TopicTombstone, error) {
		return nil, status.Error(codes.Internal, "could not delete topic")
	}
	_, err = suite.client.TopicDelete(ctx, req)
	suite.requireError(err, http.StatusNotImplemented, "deleting a topic is not supported", "expected error when Ensign returns an error")
}
