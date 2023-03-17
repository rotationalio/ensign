/*
Package mock implements an in-memory gRPC mock Ensign server that can be connected to
using a bufconn. The mock is useful for testing client side code for publishers and
subscribers without actually connecting to an Ensign server.
*/
package mock

import (
	"context"
	"errors"
	"fmt"
	"os"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	health "github.com/rotationalio/ensign/pkg/utils/probez/grpc/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// RPC Name constants based on the FullMethod that is returned from gRPC info. These
// constants can be used to reference RPCs in the mock code.
const (
	PublishRPC     = "/ensign.v1beta1.Ensign/Publish"
	SubscribeRPC   = "/ensign.v1beta1.Ensign/Subscribe"
	ListTopicsRPC  = "/ensign.v1beta1.Ensign/ListTopics"
	CreateTopicRPC = "/ensign.v1beta1.Ensign/CreateTopic"
	DeleteTopicRPC = "/ensign.v1beta1.Ensign/DeleteTopic"
	StatusRPC      = "/ensign.v1beta1.Ensign/Status"
)

// New creates a mock Ensign server for testing Ensign responses to RPC calls. If the
// bufnet is nil, the default bufconn is created for use in testing. Arbitrary server
// options (e.g. for authentication or to add interceptors) can be passed in as well.
func New(bufnet *bufconn.Listener, opts ...grpc.ServerOption) *Ensign {
	if bufnet == nil {
		bufnet = bufconn.New()
	}

	remote := &Ensign{
		bufnet: bufnet,
		srv:    grpc.NewServer(opts...),
		Calls:  make(map[string]int),
	}

	api.RegisterEnsignServer(remote.srv, remote)
	health.RegisterHealthServer(remote.srv, remote)
	go remote.srv.Serve(remote.bufnet.Sock())

	remote.Healthy()
	return remote
}

// Implements a mock gRPC server for testing Ensign client connections. The desired
// response of the Ensign server can be set by test users using the OnRPC functions or
// the WithFixture or WithError methods. The Calls map can be used to count the number
// of times a specific RPC was called.
type Ensign struct {
	health.ProbeServer
	api.UnimplementedEnsignServer

	bufnet        *bufconn.Listener
	srv           *grpc.Server
	client        api.EnsignClient
	Calls         map[string]int
	OnPublish     func(api.Ensign_PublishServer) error
	OnSubscribe   func(api.Ensign_SubscribeServer) error
	OnListTopics  func(context.Context, *api.PageInfo) (*api.TopicsPage, error)
	OnCreateTopic func(context.Context, *api.Topic) (*api.Topic, error)
	OnDeleteTopic func(context.Context, *api.TopicMod) (*api.TopicTombstone, error)
	OnStatus      func(context.Context, *api.HealthCheck) (*api.ServiceState, error)
}

// Create and connect an Ensign client to the mock server
func (s *Ensign) Client(ctx context.Context, opts ...grpc.DialOption) (client api.EnsignClient, err error) {
	if s.client == nil {
		if len(opts) == 0 {
			opts = make([]grpc.DialOption, 0, 1)
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		var cc *grpc.ClientConn
		if cc, err = s.bufnet.Connect(ctx, opts...); err != nil {
			return nil, err
		}
		s.client = api.NewEnsignClient(cc)
	}
	return s.client, nil
}

// Reset the client with the new dial options
func (s *Ensign) ResetClient(ctx context.Context, opts ...grpc.DialOption) (api.EnsignClient, error) {
	s.client = nil
	return s.Client(ctx, opts...)
}

func (s *Ensign) HealthClient(ctx context.Context, opts ...grpc.DialOption) (client health.HealthClient, err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0, 1)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	var cc *grpc.ClientConn
	if cc, err = s.bufnet.Connect(ctx, opts...); err != nil {
		return nil, err
	}
	return health.NewHealthClient(cc), nil
}

// Return the bufconn channel (helpful for dialing)
func (s *Ensign) Channel() *bufconn.Listener {
	return s.bufnet
}

// Shutdown the sever and cleanup (cannot be used after shutdown)
func (s *Ensign) Shutdown() {
	s.NotHealthy()
	s.srv.GracefulStop()
	s.bufnet.Close()
}

// Reset the calls map and all associated handlers in preparation for a new test.
func (s *Ensign) Reset() {
	for key := range s.Calls {
		s.Calls[key] = 0
	}

	s.OnPublish = nil
	s.OnSubscribe = nil
	s.OnListTopics = nil
	s.OnCreateTopic = nil
	s.OnDeleteTopic = nil
	s.OnStatus = nil
}

// UseFixture loads a JSON fixture from disk (usually in the testdata folder) to use as
// the protocol buffer response to the specified RPC, simplifying handler mocking.
func (s *Ensign) UseFixture(rpc, path string) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture: %v", err)
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch rpc {
	case PublishRPC, SubscribeRPC:
		return errors.New("cannot use fixture for a streaming RPC (yet)")
	case ListTopicsRPC:
		out := &api.TopicsPage{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnListTopics = func(context.Context, *api.PageInfo) (*api.TopicsPage, error) {
			return out, nil
		}
	case CreateTopicRPC:
		out := &api.Topic{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnCreateTopic = func(context.Context, *api.Topic) (*api.Topic, error) {
			return out, nil
		}
	case DeleteTopicRPC:
		out := &api.TopicTombstone{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnDeleteTopic = func(context.Context, *api.TopicMod) (*api.TopicTombstone, error) {
			return out, nil
		}
	case StatusRPC:
		out := &api.ServiceState{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return out, nil
		}
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}
	return nil
}

// UseError allows you to specify a gRPC status error to return from the specified RPC.
func (s *Ensign) UseError(rpc string, code codes.Code, msg string) error {
	switch rpc {
	case PublishRPC:
		s.OnPublish = func(api.Ensign_PublishServer) error {
			return status.Error(code, msg)
		}
	case SubscribeRPC:
		s.OnSubscribe = func(api.Ensign_SubscribeServer) error {
			return status.Error(code, msg)
		}
	case ListTopicsRPC:
		s.OnListTopics = func(context.Context, *api.PageInfo) (*api.TopicsPage, error) {
			return nil, status.Error(code, msg)
		}
	case CreateTopicRPC:
		s.OnCreateTopic = func(context.Context, *api.Topic) (*api.Topic, error) {
			return nil, status.Error(code, msg)
		}
	case DeleteTopicRPC:
		s.OnDeleteTopic = func(context.Context, *api.TopicMod) (*api.TopicTombstone, error) {
			return nil, status.Error(code, msg)
		}
	case StatusRPC:
		s.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return nil, status.Error(code, msg)
		}
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}
	return nil
}

func (s *Ensign) Publish(stream api.Ensign_PublishServer) error {
	s.Calls[PublishRPC]++
	return s.OnPublish(stream)
}

func (s *Ensign) Subscribe(stream api.Ensign_SubscribeServer) error {
	s.Calls[SubscribeRPC]++
	return s.OnSubscribe(stream)
}

func (s *Ensign) ListTopics(ctx context.Context, in *api.PageInfo) (*api.TopicsPage, error) {
	s.Calls[ListTopicsRPC]++
	return s.OnListTopics(ctx, in)
}

func (s *Ensign) CreateTopic(ctx context.Context, in *api.Topic) (*api.Topic, error) {
	s.Calls[CreateTopicRPC]++
	return s.OnCreateTopic(ctx, in)
}

func (s *Ensign) DeleteTopic(ctx context.Context, in *api.TopicMod) (*api.TopicTombstone, error) {
	s.Calls[DeleteTopicRPC]++
	return s.OnDeleteTopic(ctx, in)
}

func (s *Ensign) Status(ctx context.Context, in *api.HealthCheck) (*api.ServiceState, error) {
	s.Calls[StatusRPC]++
	return s.OnStatus(ctx, in)
}
