package tenant_test

import (
	"context"
	"sync"
	"time"

	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/metatopic"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/trisacrypto/directory/pkg/trtl/mock"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *tenantTestSuite) TestTopicSubscriber() {
	require := s.Require()
	orgID := ulids.New()
	projectID := ulids.New()
	topicID := ulids.New()

	createTopic := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateCreated,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "test-client",
		Topic: &metatopic.Topic{
			ID:        topicID[:],
			ProjectID: projectID[:],
			Name:      "test-topic",
			Events:    1000,
			Storage:   42,
			Publishers: &metatopic.Activity{
				Active:   1,
				Inactive: 2,
			},
			Subscribers: &metatopic.Activity{
				Active: 3,
			},
			Created:  time.Now(),
			Modified: time.Now(),
		},
	}
	data, err := createTopic.Marshal()
	require.NoError(err, "failed to marshal topic create event")
	event := &api.Event{
		Data:     data,
		Mimetype: mimetype.MIME_APPLICATION_MSGPACK,
		Type: &api.Type{
			Name:         metatopic.SchemaName,
			MajorVersion: 1,
		},
		Created: timestamppb.Now(),
	}
	topicCreateEvent := &api.EventWrapper{
		Id:        rlid.Make(1).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.Now(),
	}
	require.NoError(topicCreateEvent.Wrap(event), "failed to wrap topic create event")

	modifyTopic := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateModified,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "test-client",
		Topic: &metatopic.Topic{
			ID:         createTopic.Topic.ID,
			ProjectID:  createTopic.Topic.ProjectID,
			Name:       createTopic.Topic.Name,
			Events:     2000,
			Storage:    84,
			Publishers: createTopic.Topic.Publishers,
			Subscribers: &metatopic.Activity{
				Active: 4,
			},
			Created:  createTopic.Topic.Created,
			Modified: createTopic.Topic.Modified.Add(time.Hour),
		},
	}
	data, err = modifyTopic.Marshal()
	require.NoError(err, "failed to marshal topic modified event")
	event.Data = data
	topicModifiedEvent := &api.EventWrapper{
		Id:        rlid.Make(2).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(topicModifiedEvent.Wrap(event), "failed to wrap topic modified event")

	stateChange := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateStateChange,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "test-client",
		Topic:      nil,
	}
	data, err = stateChange.Marshal()
	require.NoError(err, "failed to marshal topic state change event")
	event.Data = data
	topicStateChangeEvent := &api.EventWrapper{
		Id:        rlid.Make(3).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(topicStateChangeEvent.Wrap(event), "failed to wrap topic state change event")

	deleteTopic := &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateDeleted,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "test-client",
		Topic:      nil,
	}
	data, err = deleteTopic.Marshal()
	require.NoError(err, "failed to marshal topic deleted event")
	event.Data = data
	topicDeletedEvent := &api.EventWrapper{
		Id:        rlid.Make(4).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(topicDeletedEvent.Wrap(event), "failed to wrap topic deleted event")

	trtl := db.GetMock()
	defer trtl.Reset()

	topic := &db.Topic{
		OrgID:     orgID,
		ProjectID: projectID,
		ID:        topicID,
		Name:      "test-topic",
	}
	key, err := topic.Key()
	require.NoError(err, "failed to get topic key")

	topicData, err := topic.MarshalValue()
	require.NoError(err, "failed to marshal topic value")

	// Configure trtl Put to verify that the correct topic is created.
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (reply *pb.PutReply, err error) {
		switch in.Namespace {
		case db.TopicNamespace:
			require.Equal(key, in.Key, "wrong key for topic put")

			topic := &db.Topic{}
			err = topic.UnmarshalValue(in.Value)
			require.NoError(err, "failed to unmarshal topic value")

			// TODO: Would be nice to be able to specify different behavior for each
			// call from the mock, rather than having to specify it here.
			calls := trtl.Calls[mock.PutRPC]
			switch {
			case calls <= 3:
				// First call should be the create.
				if topic.Name != createTopic.Topic.Name {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic name provided to put on topic create")
				}

				if topic.Events != createTopic.Topic.Events {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic events provided to put on topic create")
				}

				if topic.Storage != createTopic.Topic.Storage {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic storage provided to put on topic create")
				}

				if topic.Subscribers.Active != createTopic.Topic.Subscribers.Active {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic subscribers provided to put on topic create")
				}

			case calls == 4:
				// Second call should be the update.
				if topic.Name != modifyTopic.Topic.Name {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic name provided to put on topic update")
				}

				if topic.Events != modifyTopic.Topic.Events {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic events provided to put on topic update")
				}

				if topic.Storage != modifyTopic.Topic.Storage {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic storage provided to put on topic update")
				}

				if topic.Subscribers.Active != modifyTopic.Topic.Subscribers.Active {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic subscribers provided to put on topic update")
				}
			case calls == 5:
				// Third call should be the state change.
				if topic.Name != modifyTopic.Topic.Name {
					return nil, status.Errorf(codes.InvalidArgument, "wrong topic name provided to put on topic state change")
				}

				if topic.State != api.TopicState_READONLY {
					return nil, status.Errorf(codes.InvalidArgument, "expected topic provided to put to have readonly status on topic state change")
				}
			default:
				return nil, status.Errorf(codes.InvalidArgument, "unexpected call to trtl put")
			}

			return &pb.PutReply{}, nil
		case db.KeysNamespace, db.OrganizationNamespace:
			return &pb.PutReply{}, nil
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unexpected namespace")
		}
	}

	// Configure trtl Get to return a topic.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (reply *pb.GetReply, err error) {
		switch in.Namespace {
		case db.TopicNamespace:
			require.Equal(key, in.Key, "wrong key for topic get")
			return &pb.GetReply{
				Value: topicData,
			}, nil
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unexpected namespace")
		}
	}

	// Configure trtl Delete to verify that the correct key is deleted.
	trtl.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (reply *pb.DeleteReply, err error) {
		switch in.Namespace {
		case db.TopicNamespace:
			require.Equal(key, in.Key, "wrong key for topic delete")
			return &pb.DeleteReply{}, nil
		case db.KeysNamespace:
			return &pb.DeleteReply{}, nil
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unexpected namespace")
		}
	}

	// Ensure the server is done before shutting down the subscriber.
	server := sync.WaitGroup{}
	server.Add(1)

	// Configure the Ensign mock to emit the events.
	s.metatopic.OnSubscribe = func(stream api.Ensign_SubscribeServer) (err error) {
		defer server.Done()

		// Wait for the open subscribe request and send the stream ready response.
		if _, err = stream.Recv(); err != nil {
			return err
		}

		rep := &api.SubscribeReply{
			Embed: &api.SubscribeReply_Ready{
				Ready: &api.StreamReady{
					ClientId: "test-client",
					ServerId: "test-server",
					Topics: map[string][]byte{
						"topics": topicID[:],
					},
				},
			},
		}
		if err = stream.Send(rep); err != nil {
			return err
		}

		events := []*api.EventWrapper{
			topicCreateEvent,
			topicModifiedEvent,
			topicStateChangeEvent,
			topicDeletedEvent,
		}
		for _, event := range events {
			rep = &api.SubscribeReply{
				Embed: &api.SubscribeReply_Event{
					Event: event,
				},
			}
			if err = stream.Send(rep); err != nil {
				return err
			}
		}

		// Should receive all the acks from the subscriber.
		for i := 0; i < len(events); i++ {
			var req *api.SubscribeRequest
			if req, err = stream.Recv(); err != nil {
				return err
			}

			if req.GetAck() == nil {
				return status.Errorf(codes.InvalidArgument, "expected ack")
			}
		}

		return nil
	}

	// Run the topic subscriber with another waitgroup.
	sub := &sync.WaitGroup{}
	err = s.subscriber.Run(sub)
	require.NoError(err, "failed to run topic subscriber")

	// Wait for the server to process all the acks before stopping the subscriber.
	server.Wait()
	s.subscriber.Stop()
	sub.Wait()

	// Ensure that the subscriber actually made the database updates.
	require.Equal(5, trtl.Calls[mock.PutRPC], "expected 5 calls to Put, 3 for the topic create, 1 for the topic update, and 1 for the topic state change")
	require.Equal(2, trtl.Calls[mock.DeleteRPC], "expected 2 calls to Delete, 1 in the topic namespace and 1 in the keys namespace")
}

func (s *tenantTestSuite) TestTopicSubscriberBadEvents() {
	require := s.Require()
	topicID := ulids.New()
	projectID := ulids.New()
	orgID := ulids.New()

	event := &api.Event{
		Mimetype: mimetype.MIME_APPLICATION_JSON,
		Type: &api.Type{
			Name:         metatopic.SchemaName,
			MajorVersion: 1,
		},
		Created: timestamppb.Now(),
	}
	badMimetype := &api.EventWrapper{
		Id:        rlid.Make(1).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.Now(),
	}
	require.NoError(badMimetype.Wrap(event), "failed to wrap bad mimetype event")

	event.Mimetype = mimetype.MIME_APPLICATION_MSGPACK
	event.Type = nil
	badType := &api.EventWrapper{
		Id:        rlid.Make(2).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(badType.Wrap(event), "failed to wrap bad type event")

	event.Type = &api.Type{
		Name:         "wrong-schema",
		MajorVersion: 1,
	}
	badSchema := &api.EventWrapper{
		Id:        rlid.Make(3).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(badSchema.Wrap(event), "failed to wrap bad schema event")

	event.Data = []byte("bad-payload")
	event.Type = &api.Type{
		Name:         metatopic.SchemaName,
		MajorVersion: 1,
	}
	badPayload := &api.EventWrapper{
		Id:        rlid.Make(4).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(badPayload.Wrap(event), "failed to wrap bad payload event")

	var err error
	update := &metatopic.TopicUpdate{}
	event.Data, err = update.Marshal()
	require.NoError(err, "failed to marshal topic unknown event")
	badUpdateType := &api.EventWrapper{
		Id:        rlid.Make(5).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(badUpdateType.Wrap(event), "failed to wrap topic unknown event")

	update = &metatopic.TopicUpdate{
		UpdateType: metatopic.TopicUpdateStateChange,
		OrgID:      orgID,
		ProjectID:  projectID,
		TopicID:    topicID,
		ClientID:   "test-client",
	}
	event.Data, err = update.Marshal()
	require.NoError(err, "failed to marshal topic state change event")
	badStateChange := &api.EventWrapper{
		Id:        rlid.Make(6).Bytes(),
		TopicId:   topicID[:],
		Committed: timestamppb.New(time.Now()),
	}
	require.NoError(badStateChange.Wrap(event), "failed to wrap topic state change event")

	trtl := db.GetMock()
	defer trtl.Reset()

	// Configure trtl Get to return a not found error.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (out *pb.GetReply, err error) {
		return nil, status.Errorf(codes.NotFound, "topic not found")
	}

	// Ensure the server is done before shutting down the subscriber.
	server := sync.WaitGroup{}
	server.Add(1)

	// Configure the Ensign mock to emit the events.
	s.metatopic.OnSubscribe = func(stream api.Ensign_SubscribeServer) (err error) {
		defer server.Done()

		// Wait for the open subscribe request and send the stream ready response.
		if _, err = stream.Recv(); err != nil {
			return err
		}

		rep := &api.SubscribeReply{
			Embed: &api.SubscribeReply_Ready{
				Ready: &api.StreamReady{
					ClientId: "test-client",
					ServerId: "test-server",
					Topics: map[string][]byte{
						"topics": topicID[:],
					},
				},
			},
		}
		if err = stream.Send(rep); err != nil {
			return err
		}

		events := []*api.EventWrapper{
			badMimetype,
			badType,
			badSchema,
			badPayload,
			badUpdateType,
			badStateChange,
		}
		for _, ew := range events {
			rep = &api.SubscribeReply{
				Embed: &api.SubscribeReply_Event{
					Event: ew,
				},
			}
			if err = stream.Send(rep); err != nil {
				return err
			}
		}

		// Should receive all the nacks from the subscriber.
		for i := 0; i < len(events); i++ {
			var req *api.SubscribeRequest
			if req, err = stream.Recv(); err != nil {
				return err
			}

			if req.GetNack() == nil {
				return status.Errorf(codes.InvalidArgument, "expected nack")
			}
		}

		return nil
	}

	// Run the topic subscriber with another waitgroup.
	sub := &sync.WaitGroup{}
	err = s.subscriber.Run(sub)
	require.NoError(err, "failed to run topic subscriber")

	// Wait for the server to process all the nacks before stopping the subscriber.
	server.Wait()
	s.subscriber.Stop()
	sub.Wait()
}
