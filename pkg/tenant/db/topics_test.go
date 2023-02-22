package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTopicModel(t *testing.T) {
	topic := &db.Topic{
		OrgID:     ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
		Created:   time.Unix(1672161102, 0).In(time.UTC),
		Modified:  time.Unix(1672161102, 0).In(time.UTC),
	}

	err := topic.Validate()
	require.NoError(t, err, "could not validate topic data")

	key, err := topic.Key()
	require.NoError(t, err, "could not marshal the topic")
	require.Equal(t, topic.ProjectID[:], key[0:16], "unexpected marshaling of the project id half of the key")
	require.Equal(t, topic.ID[:], key[16:], "unexpected marshaling of the topic id half of the key")

	// Test marshal and unmarshal.
	data, err := topic.MarshalValue()
	require.NoError(t, err, "could not marshal the topic")

	other := &db.Topic{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the topic")

	TopicsEqual(t, topic, other, "unmarshal topic does not match the marshaled topic")
}

func TestTopicValidate(t *testing.T) {
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	projectID := ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88")
	topic := &db.Topic{
		OrgID:     orgID,
		ProjectID: projectID,
		Name:      "otters",
	}

	// Test missing orgID
	topic.OrgID = ulids.Null
	require.ErrorIs(t, topic.Validate(), db.ErrMissingOrgID, "expected missing org id to be an error")

	// Test missing projectID
	topic.OrgID = orgID
	topic.ProjectID = ulids.Null
	require.ErrorIs(t, topic.Validate(), db.ErrMissingProjectID, "expected missing project id to be an error")

	// Test missing name
	topic.ProjectID = projectID
	topic.Name = ""
	require.ErrorIs(t, topic.Validate(), db.ErrMissingTopicName, "expected missing name to be an error")

	// Test invalid name
	topic.Name = "otters;"
	require.ErrorIs(t, topic.Validate(), db.ErrInvalidTopicName, "expected invalid name to be an error")

	// Valid topic
	topic.Name = "otters"
	require.NoError(t, topic.Validate(), "expected valid topic to not be an error")
}

func (s *dbTestSuite) TestCreateTopic() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
		OrgID:     ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		Name:      "topic001",
	}

	err := topic.Validate()
	require.NoError(err, "could not validate topic data")

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TopicNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.CreateTopic(ctx, topic)
	require.NoError(err, "could not create topic")

	// Verify that below fields have been populated.
	require.NotEmpty(topic.ID, "expected non-zero ulid to be populated for topic id")
	require.NotEmpty(topic.Name, "topic name is required")
	require.NotZero(topic.Created, "expected topic to have a created timestamp")
	require.Equal(topic.Created, topic.Modified, "expected the same created and modified timestamp")

	// Should error if the topic is not valid.
	topic.Name = ""
	require.ErrorIs(db.CreateTopic(ctx, topic), db.ErrMissingTopicName, "expected missing topic id to be an error")

	// Test when trtl returns not found
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	topic.Name = "topic001"
	require.ErrorIs(db.CreateTopic(ctx, topic), db.ErrNotFound, "expected not found to be an error")
}

func (s *dbTestSuite) TestRetrieveTopic() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
		Created:   time.Unix(1672161102, 0),
		Modified:  time.Unix(1672161102, 0),
	}

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TopicNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}
		if !bytes.Equal(in.Key[16:], topic.ID[:]) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		// TODO: Add msgpack fixture helpers.

		// Marshal the data with msgpack.
		data, err := topic.MarshalValue()
		require.NoError(err, "could not marshal the data")

		other := &db.Topic{}
		err = other.UnmarshalValue(data)
		require.NoError(err, "could not unmarshal the data")

		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	topic, err := db.RetrieveTopic(ctx, topic.ID)
	require.NoError(err, "could not retrieve topic")

	// Verify the fields below have been populated.
	require.Equal(ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"), topic.ProjectID, "expected project id to match")
	require.Equal(ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"), topic.ID, "expected topic id to match")
	require.Equal("topic001", topic.Name, "expected topic name to match")
	require.Equal(time.Unix(1672161102, 0), topic.Created, "expected created timestamp to have not changed")

	// Test NotFound path.
	_, err = db.RetrieveTopic(ctx, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListTopics() {
	require := s.Require()
	ctx := context.Background()
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

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
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

	values, err := db.List(ctx, prefix, namespace)
	require.NoError(err, "could not get topic values")
	require.Len(values, 3, "expected 3 values")

	rep, err := db.ListTopics(ctx, projectID)
	require.NoError(err, "could not list topics")
	require.Len(rep, 3, "expected 3 topics")

	// Test first topic data has been populated.
	require.Equal(topics[0].ID, rep[0].ID, "expected topic id to match")
	require.Equal(topics[0].Name, rep[0].Name, "expected topic name to match")

	// Test second topic data has been populated.
	require.Equal(topics[1].ID, rep[1].ID, "expected topic id to match")
	require.Equal(topics[1].Name, rep[1].Name, "expected topic name to match")

	// Test third topic data has been populated.
	require.Equal(topics[2].ID, rep[2].ID, "expected topic id to match")
	require.Equal(topics[2].Name, rep[2].Name, "expected topic name to match")
}

func (s *dbTestSuite) TestUpdateTopic() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
		OrgID:     ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
		Created:   time.Unix(1672161102, 0),
		Modified:  time.Unix(1672161102, 0),
	}

	err := topic.Validate()
	require.NoError(err, "could not validate topic data")

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TopicNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}
		if !bytes.Equal(in.Key[16:], topic.ID[:]) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.UpdateTopic(ctx, topic)
	require.NoError(err, "could not update topic")

	// Fields below should have been populated.
	require.Equal(ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"), topic.ProjectID, "project ID should not have changed")
	require.Equal(ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"), topic.ID, "topic ID should not have changed")
	require.Equal(time.Unix(1672161102, 0), topic.Created, "expected created timestamp to have not changed")
	require.True(time.Unix(1672161102, 0).Before(topic.Modified), "expected modified timestamp to be updated")

	// If created timestamp is missing then it should be populated.
	topic.Created = time.Time{}
	require.NoError(db.UpdateTopic(ctx, topic), "could not update topic")
	require.Equal(topic.Modified, topic.Created, "expected created timestamp to be populated")

	// Should fail if topic ID is missing
	topic.ID = ulid.ULID{}
	require.ErrorIs(db.UpdateTopic(ctx, topic), db.ErrMissingID, "expected invalid topic ID error")

	// Should fail if topic is not valid
	topic.ID = ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6")
	topic.Name = ""
	require.ErrorIs(db.UpdateTopic(ctx, topic), db.ErrMissingTopicName, "expected invalid topic error")

	// Test NotFound path.
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "topic not found")
	}
	err = db.UpdateTopic(ctx, &db.Topic{OrgID: ulids.New(), ProjectID: ulids.New(), ID: ulids.New(), Name: "topic002"})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteTopic() {
	require := s.Require()
	ctx := context.Background()
	topicID := ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6")

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TopicNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}
		if !bytes.Equal(in.Key[16:], topicID[:]) {
			return nil, status.Errorf(codes.NotFound, "topic not found")
		}
		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteTopic(ctx, topicID)
	require.NoError(err, "could not delete topic")

	// Test NotFound path.
	err = db.DeleteTopic(ctx, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)

}

// TopicsEqual tests assertions in the TopicModel.
// Note: require.True compares the actual.Created and actual.Modified timestamps
// because MsgPack does not preserve time zone information.
func TopicsEqual(t *testing.T, expected, actual *db.Topic, msgAndArgs ...interface{}) {
	require.Equal(t, expected.ProjectID, actual.ProjectID, msgAndArgs...)
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}
