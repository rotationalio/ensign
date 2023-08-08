// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: raft/v1beta1/raft.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Raft_RequestVote_FullMethodName   = "/raft.v1beta1.Raft/RequestVote"
	Raft_AppendEntries_FullMethodName = "/raft.v1beta1.Raft/AppendEntries"
)

// RaftClient is the client API for Raft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftClient interface {
	RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error)
	AppendEntries(ctx context.Context, opts ...grpc.CallOption) (Raft_AppendEntriesClient, error)
}

type raftClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftClient(cc grpc.ClientConnInterface) RaftClient {
	return &raftClient{cc}
}

func (c *raftClient) RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error) {
	out := new(VoteReply)
	err := c.cc.Invoke(ctx, Raft_RequestVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) AppendEntries(ctx context.Context, opts ...grpc.CallOption) (Raft_AppendEntriesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Raft_ServiceDesc.Streams[0], Raft_AppendEntries_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &raftAppendEntriesClient{stream}
	return x, nil
}

type Raft_AppendEntriesClient interface {
	Send(*AppendRequest) error
	Recv() (*AppendReply, error)
	grpc.ClientStream
}

type raftAppendEntriesClient struct {
	grpc.ClientStream
}

func (x *raftAppendEntriesClient) Send(m *AppendRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *raftAppendEntriesClient) Recv() (*AppendReply, error) {
	m := new(AppendReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RaftServer is the server API for Raft service.
// All implementations must embed UnimplementedRaftServer
// for forward compatibility
type RaftServer interface {
	RequestVote(context.Context, *VoteRequest) (*VoteReply, error)
	AppendEntries(Raft_AppendEntriesServer) error
	mustEmbedUnimplementedRaftServer()
}

// UnimplementedRaftServer must be embedded to have forward compatible implementations.
type UnimplementedRaftServer struct {
}

func (UnimplementedRaftServer) RequestVote(context.Context, *VoteRequest) (*VoteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedRaftServer) AppendEntries(Raft_AppendEntriesServer) error {
	return status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedRaftServer) mustEmbedUnimplementedRaftServer() {}

// UnsafeRaftServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftServer will
// result in compilation errors.
type UnsafeRaftServer interface {
	mustEmbedUnimplementedRaftServer()
}

func RegisterRaftServer(s grpc.ServiceRegistrar, srv RaftServer) {
	s.RegisterService(&Raft_ServiceDesc, srv)
}

func _Raft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Raft_RequestVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).RequestVote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_AppendEntries_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RaftServer).AppendEntries(&raftAppendEntriesServer{stream})
}

type Raft_AppendEntriesServer interface {
	Send(*AppendReply) error
	Recv() (*AppendRequest, error)
	grpc.ServerStream
}

type raftAppendEntriesServer struct {
	grpc.ServerStream
}

func (x *raftAppendEntriesServer) Send(m *AppendReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *raftAppendEntriesServer) Recv() (*AppendRequest, error) {
	m := new(AppendRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Raft_ServiceDesc is the grpc.ServiceDesc for Raft service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Raft_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft.v1beta1.Raft",
	HandlerType: (*RaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _Raft_RequestVote_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AppendEntries",
			Handler:       _Raft_AppendEntries_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "raft/v1beta1/raft.proto",
}
