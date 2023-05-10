// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/v1beta1/ensign.proto

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

// EnsignClient is the client API for Ensign service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnsignClient interface {
	// Both the Publish and Subscribe RPCs are bidirectional streaming to allow for acks
	// and nacks of events to be sent between Ensign and the client. The Publish stream
	// is opened and the client sends events and receives acks/nacks -- when the client
	// closes the publish stream, the server sends back information about the current
	// state of the topic. When the Subscribe stream is opened, the client must send an
	// open stream message with the subscription info before receiving events. Once it
	// receives events it must send back acks/nacks up the stream so that Ensign
	// advances the topic offset for the rest of the clients in the group.
	Publish(ctx context.Context, opts ...grpc.CallOption) (Ensign_PublishClient, error)
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (Ensign_SubscribeClient, error)
	// This is a simple topic management interface. Right now we assume that topics are
	// immutable, therefore there is no update topic RPC call. There are two ways to
	// delete a topic - archiving it makes the topic readonly so that no events can be
	// published to it, but it can still be read. Destroying the topic deletes it and
	// removes all of its data, freeing up the topic name to be used again.
	ListTopics(ctx context.Context, in *PageInfo, opts ...grpc.CallOption) (*TopicsPage, error)
	CreateTopic(ctx context.Context, in *Topic, opts ...grpc.CallOption) (*Topic, error)
	RetrieveTopic(ctx context.Context, in *Topic, opts ...grpc.CallOption) (*Topic, error)
	DeleteTopic(ctx context.Context, in *TopicMod, opts ...grpc.CallOption) (*TopicTombstone, error)
	TopicNames(ctx context.Context, in *PageInfo, opts ...grpc.CallOption) (*TopicNamesPage, error)
	TopicExists(ctx context.Context, in *TopicName, opts ...grpc.CallOption) (*TopicExistsInfo, error)
	// Info provides statistics and metrics describing the state of a project
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*ProjectInfo, error)
	// Implements a client-side heartbeat that can also be used by monitoring tools.
	Status(ctx context.Context, in *HealthCheck, opts ...grpc.CallOption) (*ServiceState, error)
}

type ensignClient struct {
	cc grpc.ClientConnInterface
}

func NewEnsignClient(cc grpc.ClientConnInterface) EnsignClient {
	return &ensignClient{cc}
}

func (c *ensignClient) Publish(ctx context.Context, opts ...grpc.CallOption) (Ensign_PublishClient, error) {
	stream, err := c.cc.NewStream(ctx, &Ensign_ServiceDesc.Streams[0], "/ensign.v1beta1.Ensign/Publish", opts...)
	if err != nil {
		return nil, err
	}
	x := &ensignPublishClient{stream}
	return x, nil
}

type Ensign_PublishClient interface {
	Send(*PublisherRequest) error
	Recv() (*PublisherReply, error)
	grpc.ClientStream
}

type ensignPublishClient struct {
	grpc.ClientStream
}

func (x *ensignPublishClient) Send(m *PublisherRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ensignPublishClient) Recv() (*PublisherReply, error) {
	m := new(PublisherReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ensignClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (Ensign_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Ensign_ServiceDesc.Streams[1], "/ensign.v1beta1.Ensign/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &ensignSubscribeClient{stream}
	return x, nil
}

type Ensign_SubscribeClient interface {
	Send(*SubscribeRequest) error
	Recv() (*SubscribeReply, error)
	grpc.ClientStream
}

type ensignSubscribeClient struct {
	grpc.ClientStream
}

func (x *ensignSubscribeClient) Send(m *SubscribeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ensignSubscribeClient) Recv() (*SubscribeReply, error) {
	m := new(SubscribeReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ensignClient) ListTopics(ctx context.Context, in *PageInfo, opts ...grpc.CallOption) (*TopicsPage, error) {
	out := new(TopicsPage)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/ListTopics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) CreateTopic(ctx context.Context, in *Topic, opts ...grpc.CallOption) (*Topic, error) {
	out := new(Topic)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/CreateTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) RetrieveTopic(ctx context.Context, in *Topic, opts ...grpc.CallOption) (*Topic, error) {
	out := new(Topic)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/RetrieveTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) DeleteTopic(ctx context.Context, in *TopicMod, opts ...grpc.CallOption) (*TopicTombstone, error) {
	out := new(TopicTombstone)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/DeleteTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) TopicNames(ctx context.Context, in *PageInfo, opts ...grpc.CallOption) (*TopicNamesPage, error) {
	out := new(TopicNamesPage)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/TopicNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) TopicExists(ctx context.Context, in *TopicName, opts ...grpc.CallOption) (*TopicExistsInfo, error) {
	out := new(TopicExistsInfo)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/TopicExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*ProjectInfo, error) {
	out := new(ProjectInfo)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensignClient) Status(ctx context.Context, in *HealthCheck, opts ...grpc.CallOption) (*ServiceState, error) {
	out := new(ServiceState)
	err := c.cc.Invoke(ctx, "/ensign.v1beta1.Ensign/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnsignServer is the server API for Ensign service.
// All implementations must embed UnimplementedEnsignServer
// for forward compatibility
type EnsignServer interface {
	// Both the Publish and Subscribe RPCs are bidirectional streaming to allow for acks
	// and nacks of events to be sent between Ensign and the client. The Publish stream
	// is opened and the client sends events and receives acks/nacks -- when the client
	// closes the publish stream, the server sends back information about the current
	// state of the topic. When the Subscribe stream is opened, the client must send an
	// open stream message with the subscription info before receiving events. Once it
	// receives events it must send back acks/nacks up the stream so that Ensign
	// advances the topic offset for the rest of the clients in the group.
	Publish(Ensign_PublishServer) error
	Subscribe(Ensign_SubscribeServer) error
	// This is a simple topic management interface. Right now we assume that topics are
	// immutable, therefore there is no update topic RPC call. There are two ways to
	// delete a topic - archiving it makes the topic readonly so that no events can be
	// published to it, but it can still be read. Destroying the topic deletes it and
	// removes all of its data, freeing up the topic name to be used again.
	ListTopics(context.Context, *PageInfo) (*TopicsPage, error)
	CreateTopic(context.Context, *Topic) (*Topic, error)
	RetrieveTopic(context.Context, *Topic) (*Topic, error)
	DeleteTopic(context.Context, *TopicMod) (*TopicTombstone, error)
	TopicNames(context.Context, *PageInfo) (*TopicNamesPage, error)
	TopicExists(context.Context, *TopicName) (*TopicExistsInfo, error)
	// Info provides statistics and metrics describing the state of a project
	Info(context.Context, *InfoRequest) (*ProjectInfo, error)
	// Implements a client-side heartbeat that can also be used by monitoring tools.
	Status(context.Context, *HealthCheck) (*ServiceState, error)
	mustEmbedUnimplementedEnsignServer()
}

// UnimplementedEnsignServer must be embedded to have forward compatible implementations.
type UnimplementedEnsignServer struct {
}

func (UnimplementedEnsignServer) Publish(Ensign_PublishServer) error {
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedEnsignServer) Subscribe(Ensign_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedEnsignServer) ListTopics(context.Context, *PageInfo) (*TopicsPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopics not implemented")
}
func (UnimplementedEnsignServer) CreateTopic(context.Context, *Topic) (*Topic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTopic not implemented")
}
func (UnimplementedEnsignServer) RetrieveTopic(context.Context, *Topic) (*Topic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveTopic not implemented")
}
func (UnimplementedEnsignServer) DeleteTopic(context.Context, *TopicMod) (*TopicTombstone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTopic not implemented")
}
func (UnimplementedEnsignServer) TopicNames(context.Context, *PageInfo) (*TopicNamesPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TopicNames not implemented")
}
func (UnimplementedEnsignServer) TopicExists(context.Context, *TopicName) (*TopicExistsInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TopicExists not implemented")
}
func (UnimplementedEnsignServer) Info(context.Context, *InfoRequest) (*ProjectInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedEnsignServer) Status(context.Context, *HealthCheck) (*ServiceState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedEnsignServer) mustEmbedUnimplementedEnsignServer() {}

// UnsafeEnsignServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnsignServer will
// result in compilation errors.
type UnsafeEnsignServer interface {
	mustEmbedUnimplementedEnsignServer()
}

func RegisterEnsignServer(s grpc.ServiceRegistrar, srv EnsignServer) {
	s.RegisterService(&Ensign_ServiceDesc, srv)
}

func _Ensign_Publish_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EnsignServer).Publish(&ensignPublishServer{stream})
}

type Ensign_PublishServer interface {
	Send(*PublisherReply) error
	Recv() (*PublisherRequest, error)
	grpc.ServerStream
}

type ensignPublishServer struct {
	grpc.ServerStream
}

func (x *ensignPublishServer) Send(m *PublisherReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ensignPublishServer) Recv() (*PublisherRequest, error) {
	m := new(PublisherRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Ensign_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EnsignServer).Subscribe(&ensignSubscribeServer{stream})
}

type Ensign_SubscribeServer interface {
	Send(*SubscribeReply) error
	Recv() (*SubscribeRequest, error)
	grpc.ServerStream
}

type ensignSubscribeServer struct {
	grpc.ServerStream
}

func (x *ensignSubscribeServer) Send(m *SubscribeReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ensignSubscribeServer) Recv() (*SubscribeRequest, error) {
	m := new(SubscribeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Ensign_ListTopics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).ListTopics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/ListTopics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).ListTopics(ctx, req.(*PageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_CreateTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Topic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).CreateTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/CreateTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).CreateTopic(ctx, req.(*Topic))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_RetrieveTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Topic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).RetrieveTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/RetrieveTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).RetrieveTopic(ctx, req.(*Topic))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_DeleteTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicMod)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).DeleteTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/DeleteTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).DeleteTopic(ctx, req.(*TopicMod))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_TopicNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).TopicNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/TopicNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).TopicNames(ctx, req.(*PageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_TopicExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).TopicExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/TopicExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).TopicExists(ctx, req.(*TopicName))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ensign_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheck)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsignServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ensign.v1beta1.Ensign/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsignServer).Status(ctx, req.(*HealthCheck))
	}
	return interceptor(ctx, in, info, handler)
}

// Ensign_ServiceDesc is the grpc.ServiceDesc for Ensign service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ensign_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ensign.v1beta1.Ensign",
	HandlerType: (*EnsignServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTopics",
			Handler:    _Ensign_ListTopics_Handler,
		},
		{
			MethodName: "CreateTopic",
			Handler:    _Ensign_CreateTopic_Handler,
		},
		{
			MethodName: "RetrieveTopic",
			Handler:    _Ensign_RetrieveTopic_Handler,
		},
		{
			MethodName: "DeleteTopic",
			Handler:    _Ensign_DeleteTopic_Handler,
		},
		{
			MethodName: "TopicNames",
			Handler:    _Ensign_TopicNames_Handler,
		},
		{
			MethodName: "TopicExists",
			Handler:    _Ensign_TopicExists_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Ensign_Info_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Ensign_Status_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Publish",
			Handler:       _Ensign_Publish_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _Ensign_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/v1beta1/ensign.proto",
}
