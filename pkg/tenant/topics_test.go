package tenant_test

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (suite *tenantTestSuite) TestTopicDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	topic := &db.Topic{
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic-example",
	}

	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

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

	// Should return an error if the topic does not exist.
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic does not exist")

	// Create a topic test fixture.
	req := &api.Topic{
		ID:   "01GNA926JCTKDH3VZBTJM8MAF6",
		Name: "topic-example",
	}

	rep, err := suite.client.TopicDetail(ctx, req.ID)
	require.NoError(err, "could not retrieve topic")
	require.Equal(req.ID, rep.ID, "expected topic ID to match")
	require.Equal(req.Name, rep.Name, "expected topic name to match")

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
	topic := &db.Topic{
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic-example",
	}

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the topic data with msgpack.
	data, err := topic.MarshalValue()
	require.NoError(err, "could not marshal the topic data")

	// Unmarshal the topic data with msgpack.
	other := &db.Topic{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the topic data")

	// Call the OnGet method and return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return an error if the topic does not exist.
	_, err = suite.client.TopicDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse topic ulid", "expected error when topic does not exist")

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if the topic name does not exist.
	_, err = suite.client.TopicUpdate(ctx, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6"})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name does not exist")

	// Create a topic test fixture.
	req := &api.Topic{
		ID:   "01GNA926JCTKDH3VZBTJM8MAF6",
		Name: "topic-example",
	}

	rep, err := suite.client.TopicUpdate(ctx, req)
	require.NoError(err, "could not update topic")
	require.NotEqual(req.ID, "", "topic id should not match")
	require.Equal(req.Name, rep.Name, "expected topic name to match")

	// Should return an error if the topic ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, errors.New("key not found")
	}

	_, err = suite.client.TopicDetail(ctx, "01GNA926JCTKDH3VZBTJM8MAF6")
	suite.requireError(err, http.StatusNotFound, "could not retrieve topic", "expected error when topic ID is not found")
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

	// Should return an error if the topic does not exist.
	err := suite.client.TopicDelete(ctx, "invalid")
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

func (suite *tenantTestSuite) TestProjectTopicCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projectID := ulid.Make().String()

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if project id is not a valid ULID.
	_, err := suite.client.ProjectTopicCreate(ctx, "projectID", &api.Topic{ID: "", Name: "topic-example"})
	suite.requireError(err, http.StatusBadRequest, "could not parse project id", "expected error when project id does not exist")

	// Should return an error if topic id exists.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "01GNA926JCTKDH3VZBTJM8MAF6", Name: "topic-example"})
	suite.requireError(err, http.StatusBadRequest, "topic id cannot be specified on create", "expected error when topic id exists")

	// Should return an error if topic name does not exist.
	_, err = suite.client.ProjectTopicCreate(ctx, projectID, &api.Topic{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "topic name is required", "expected error when topic name does not exist")

	req := &api.Topic{
		Name: "topic-example",
	}

	topic, err := suite.client.ProjectTopicCreate(ctx, projectID, req)
	require.NoError(err, "could not add topic")
	require.Equal(req.Name, topic.Name, "expected topic name to match")
}
