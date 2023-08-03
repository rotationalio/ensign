package mock

import (
	"context"
	"fmt"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"google.golang.org/grpc/metadata"
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
