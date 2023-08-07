package mock

import (
	"context"
	"fmt"
	"io"
	"sync"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	StreamSend       = "Send"
	StreamRecv       = "Recv"
	StreamSetHeader  = "SetHeader"
	StreamSendHeader = "SendHeader"
	StreamSetTrailer = "SetTrailer"
	StreamContext    = "Context"
	StreamSendMsg    = "SendMsg"
	StreamRecvMsg    = "RecvMsg"
)

// Implements api.Ensign_PublishServer for testing the Publish streaming RPC.
type PublisherServer struct {
	ServerStream

	OnSend func(*api.PublisherReply) error
	OnRecv func() (*api.PublisherRequest, error)
}

func (s *PublisherServer) Send(m *api.PublisherReply) error {
	s.incrCalls(StreamSend)
	if s.OnSend != nil {
		return s.OnSend(m)
	}
	return nil
}

func (s *PublisherServer) Recv() (*api.PublisherRequest, error) {
	s.incrCalls(StreamRecv)
	if s.OnRecv != nil {
		return s.OnRecv()
	}
	return nil, nil
}

// WithError ensures that the next call to the specified method returns an error.
func (s *PublisherServer) WithError(call string, err error) {
	switch call {
	case StreamSend:
		s.OnSend = func(*api.PublisherReply) error { return err }
	case StreamRecv:
		s.OnRecv = func() (*api.PublisherRequest, error) { return nil, err }
	default:
		s.ServerStream.WithError(call, err)
	}
}

// WithEvents creates a Recv method that sends the given open stream message, then sends
// each event before finally sending an io.EOF message to close the stream.
func (s *PublisherServer) WithEvents(info *api.OpenStream, events ...*api.EventWrapper) {
	nsent := -1
	s.OnRecv = func() (*api.PublisherRequest, error) {
		// Ensure that nsent is incremented on each call.
		defer func() { nsent++ }()

		// Send the open stream event if we haven't sent any events yet.
		if nsent < 0 {
			return &api.PublisherRequest{
				Embed: &api.PublisherRequest_OpenStream{
					OpenStream: info,
				},
			}, nil
		}

		// If we've exhausted all the events, send an EOF
		if nsent+1 > len(events) {
			return nil, io.EOF
		}

		// Send the event to the publisher
		return &api.PublisherRequest{
			Embed: &api.PublisherRequest_Event{
				Event: events[nsent],
			},
		}, nil
	}
}

// Capture returns any replies sent by the server on the specified channel.
func (s *PublisherServer) Capture(replies chan<- *api.PublisherReply) {
	s.OnSend = func(msg *api.PublisherReply) error {
		replies <- msg
		return nil
	}
}

// Implements api.Ensign_SubscribeServer for testing the Subscribe streaming RPC.
type SubscribeServer struct {
	ServerStream

	OnSend func(*api.SubscribeReply) error
	OnRecv func() (*api.SubscribeRequest, error)
}

func (s *SubscribeServer) Send(m *api.SubscribeReply) error {
	s.incrCalls(StreamSend)
	if s.OnSend != nil {
		return s.OnSend(m)
	}
	return nil
}

func (s *SubscribeServer) Recv() (*api.SubscribeRequest, error) {
	s.incrCalls(StreamRecv)
	if s.OnRecv != nil {
		return s.OnRecv()
	}
	return nil, nil
}

// WithError ensures that the next call to the specified method returns an error.
func (s *SubscribeServer) WithError(call string, err error) {
	switch call {
	case StreamSend:
		s.OnSend = func(*api.SubscribeReply) error { return err }
	case StreamRecv:
		s.OnRecv = func() (*api.SubscribeRequest, error) { return nil, err }
	default:
		s.ServerStream.WithError(call, err)
	}
}

// WithSubscription creates an object that allows test code to receive events and send
// acks and nacks on the specified subscription channel.
func (s *SubscribeServer) WithSubscription(subscription *api.Subscription) *Subscription {
	sub := &Subscription{
		closed:   false,
		requests: make(chan *api.SubscribeRequest, 1),
		replies:  make(chan *api.SubscribeReply, 1),
	}

	// Add the open stream message to the queue
	sub.requests <- &api.SubscribeRequest{
		Embed: &api.SubscribeRequest_Subscription{
			Subscription: subscription,
		},
	}

	s.OnRecv = func() (*api.SubscribeRequest, error) {
		sub.RLock()
		defer sub.RUnlock()

		if sub.closed {
			return nil, io.EOF
		}

		msg := <-sub.requests
		return msg, nil
	}

	s.OnSend = func(msg *api.SubscribeReply) error {
		sub.RLock()
		defer sub.RUnlock()

		if sub.closed {
			return io.EOF
		}

		sub.replies <- msg
		return nil
	}

	return sub
}

type Subscription struct {
	sync.RWMutex
	closed   bool
	requests chan *api.SubscribeRequest
	replies  chan *api.SubscribeReply
}

func (s *Subscription) Close() {
	s.Lock()
	defer s.Unlock()
	s.closed = true
	close(s.replies)
	close(s.requests)
}

func (s *Subscription) Ready() *api.StreamReady {
	msg := <-s.replies
	return msg.GetReady()
}

func (s *Subscription) Next() *api.EventWrapper {
	msg := <-s.replies
	return msg.GetEvent()
}

func (s *Subscription) Ack(id []byte) {
	s.requests <- &api.SubscribeRequest{
		Embed: &api.SubscribeRequest_Ack{
			Ack: &api.Ack{
				Id:        id,
				Committed: timestamppb.Now(),
			},
		},
	}
}

func (s *Subscription) Nack(id []byte, code api.Nack_Code, msg string) {
	s.requests <- &api.SubscribeRequest{
		Embed: &api.SubscribeRequest_Nack{
			Nack: &api.Nack{
				Id:    id,
				Code:  code,
				Error: msg,
			},
		},
	}
}

// Implements the grpc.ServerStream interface for testing streaming RPCs.
type ServerStream struct {
	Calls map[string]int

	OnSetHeader  func(metadata.MD) error
	OnSendHeader func(metadata.MD) error
	OnSetTrailer func(metadata.MD)
	OnContext    func() context.Context
	OnSendMsg    func(interface{}) error
	OnRecvMsg    func(interface{}) error
}

// WithContext ensures the server stream returns the specified context.
func (s *ServerStream) WithContext(ctx context.Context) {
	s.OnContext = func() context.Context {
		return ctx
	}
}

// WithClaims creates a context with the specified claims on it.
func (s *ServerStream) WithClaims(claims *tokens.Claims) {
	ctx := contexts.WithClaims(context.Background(), claims)
	s.WithContext(ctx)
}

// WithPeer sets the peer on the server context in addition to the claims.
func (s *ServerStream) WithPeer(claims *tokens.Claims, remote *peer.Peer) {
	ctx := contexts.WithClaims(context.Background(), claims)
	ctx = peer.NewContext(ctx, remote)
	s.WithContext(ctx)
}

// WithError ensures that the next call to the specified method returns an error.
func (s *ServerStream) WithError(call string, err error) {
	switch call {
	case StreamSetHeader:
		s.OnSetHeader = func(metadata.MD) error { return err }
	case StreamSendHeader:
		s.OnSendHeader = func(metadata.MD) error { return err }
	case StreamSendMsg:
		s.OnSendMsg = func(interface{}) error { return err }
	case StreamRecvMsg:
		s.OnRecvMsg = func(interface{}) error { return err }
	default:
		panic(fmt.Errorf("unknown call %q", call))
	}
}

func (s *ServerStream) SetHeader(m metadata.MD) error {
	s.incrCalls(StreamSetHeader)
	if s.OnSetHeader != nil {
		return s.OnSetHeader(m)
	}
	return nil
}

func (s *ServerStream) SendHeader(m metadata.MD) error {
	s.incrCalls(StreamSendHeader)
	if s.OnSendHeader != nil {
		return s.OnSendHeader(m)
	}
	return nil
}

func (s *ServerStream) SetTrailer(m metadata.MD) {
	s.incrCalls(StreamSetTrailer)
	if s.OnSetTrailer != nil {
		s.OnSetTrailer(m)
	}
}

func (s *ServerStream) Context() context.Context {
	s.incrCalls(StreamContext)
	if s.OnContext != nil {
		return s.OnContext()
	}
	return context.Background()
}

func (s *ServerStream) SendMsg(m interface{}) error {
	s.incrCalls(StreamSendMsg)
	if s.OnSendMsg != nil {
		return s.OnSendMsg(m)
	}
	return nil
}

func (s *ServerStream) RecvMsg(m interface{}) error {
	s.incrCalls(StreamRecvMsg)
	if s.OnRecvMsg != nil {
		return s.OnRecvMsg(m)
	}
	return nil
}

func (s *ServerStream) incrCalls(call string) {
	if s.Calls == nil {
		s.Calls = make(map[string]int)
	}
	s.Calls[call]++
}
