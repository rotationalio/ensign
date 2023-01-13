package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTopicModel(t *testing.T) {
	topic := &db.Topic{
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

func (s *dbTestSuite) TestCreateTopic() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
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
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	_, err = db.RetrieveTopic(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListTopics() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
		ProjectID: ulid.MustParse("01GNA91N6WMCWNG9MVSK47ZS88"),
		ID:        ulid.MustParse("01GNA926JCTKDH3VZBTJM8MAF6"),
		Name:      "topic001",
		Created:   time.Unix(1672161102, 0),
		Modified:  time.Unix(1672161102, 0),
	}

	prefix := topic.ProjectID[:]
	namespace := "topics"

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate.
		for i := 0; i < 7; i++ {
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     []byte(fmt.Sprintf("value %d", i)),
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	values, err := db.List(ctx, prefix, namespace)
	require.NoError(err, "could not get topic values")
	require.Len(values, 7)

	topics := make([]*db.Topic, 0, len(values))
	topics = append(topics, topic)
	require.Len(topics, 1)

	_, err = db.ListTopics(ctx, topic.ProjectID)
	require.Error(err, "could not list topics")
}

func (s *dbTestSuite) TestUpdateTopic() {
	require := s.Require()
	ctx := context.Background()
	topic := &db.Topic{
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

	// Test NotFound path.
	// TODO: Use crypto rand and monotonic entropy with ulid.New.
	err = db.UpdateTopic(ctx, &db.Topic{ProjectID: ulid.Make(), ID: ulid.Make(), Name: "topic002"})
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
	// TODO: Use crypto rand and monotonic entropy with ulid.New.
	err = db.DeleteTopic(ctx, ulid.Make())
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
