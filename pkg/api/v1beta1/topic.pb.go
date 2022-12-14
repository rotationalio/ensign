// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: ensign/v1beta1/topic.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TopicMod_Operation int32

const (
	TopicMod_NOOP    TopicMod_Operation = 0
	TopicMod_ARCHIVE TopicMod_Operation = 1 // makes the topic readonly
	TopicMod_DESTROY TopicMod_Operation = 2 // deletes the topic and removes all of its data
)

// Enum value maps for TopicMod_Operation.
var (
	TopicMod_Operation_name = map[int32]string{
		0: "NOOP",
		1: "ARCHIVE",
		2: "DESTROY",
	}
	TopicMod_Operation_value = map[string]int32{
		"NOOP":    0,
		"ARCHIVE": 1,
		"DESTROY": 2,
	}
)

func (x TopicMod_Operation) Enum() *TopicMod_Operation {
	p := new(TopicMod_Operation)
	*p = x
	return p
}

func (x TopicMod_Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TopicMod_Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_ensign_v1beta1_topic_proto_enumTypes[0].Descriptor()
}

func (TopicMod_Operation) Type() protoreflect.EnumType {
	return &file_ensign_v1beta1_topic_proto_enumTypes[0]
}

func (x TopicMod_Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TopicMod_Operation.Descriptor instead.
func (TopicMod_Operation) EnumDescriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{2, 0}
}

type TopicTombstone_Status int32

const (
	TopicTombstone_UNKNOWN  TopicTombstone_Status = 0
	TopicTombstone_READONLY TopicTombstone_Status = 1
	TopicTombstone_DELETING TopicTombstone_Status = 2
)

// Enum value maps for TopicTombstone_Status.
var (
	TopicTombstone_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "READONLY",
		2: "DELETING",
	}
	TopicTombstone_Status_value = map[string]int32{
		"UNKNOWN":  0,
		"READONLY": 1,
		"DELETING": 2,
	}
)

func (x TopicTombstone_Status) Enum() *TopicTombstone_Status {
	p := new(TopicTombstone_Status)
	*p = x
	return p
}

func (x TopicTombstone_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TopicTombstone_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_ensign_v1beta1_topic_proto_enumTypes[1].Descriptor()
}

func (TopicTombstone_Status) Type() protoreflect.EnumType {
	return &file_ensign_v1beta1_topic_proto_enumTypes[1]
}

func (x TopicTombstone_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TopicTombstone_Status.Descriptor instead.
func (TopicTombstone_Status) EnumDescriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{3, 0}
}

// Topics are collections of related events and the events inside of a topic are totally
// ordered by ID and their log index. Topics must define the event types and regions
// that they are operated on, which will allow Ensign to determine how to distribute the
// topic over multiple nodes. Users must use the topic ID to connect to a publish or
// subscribe stream. Users can create and delete topics, but for the current
// implementation, topics are immutable -- meaning that they cannot be changed. Topics
// can be deleted in two ways: they can be archived (making them readonly) or they can
// be destroyed, which removes the name of the topic and all the events in the topic.
type Topic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Types    []*Type                `protobuf:"bytes,3,rep,name=types,proto3" json:"types,omitempty"`
	Regions  []*Region              `protobuf:"bytes,4,rep,name=regions,proto3" json:"regions,omitempty"`
	Readonly bool                   `protobuf:"varint,14,opt,name=readonly,proto3" json:"readonly,omitempty"`
	Created  *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *Topic) Reset() {
	*x = Topic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Topic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Topic) ProtoMessage() {}

func (x *Topic) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Topic.ProtoReflect.Descriptor instead.
func (*Topic) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{0}
}

func (x *Topic) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Topic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Topic) GetTypes() []*Type {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *Topic) GetRegions() []*Region {
	if x != nil {
		return x.Regions
	}
	return nil
}

func (x *Topic) GetReadonly() bool {
	if x != nil {
		return x.Readonly
	}
	return false
}

func (x *Topic) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

// A list of paginated topics the user can use to identify topic ids to subscribe to.
type TopicsPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topics        []*Topic `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
	NextPageToken string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *TopicsPage) Reset() {
	*x = TopicsPage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicsPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicsPage) ProtoMessage() {}

func (x *TopicsPage) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicsPage.ProtoReflect.Descriptor instead.
func (*TopicsPage) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{1}
}

func (x *TopicsPage) GetTopics() []*Topic {
	if x != nil {
		return x.Topics
	}
	return nil
}

func (x *TopicsPage) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// A topic modification operation to archive or destroy the topic.
type TopicMod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Operation TopicMod_Operation `protobuf:"varint,2,opt,name=operation,proto3,enum=ensign.v1beta1.TopicMod_Operation" json:"operation,omitempty"`
}

func (x *TopicMod) Reset() {
	*x = TopicMod{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicMod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicMod) ProtoMessage() {}

func (x *TopicMod) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicMod.ProtoReflect.Descriptor instead.
func (*TopicMod) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{2}
}

func (x *TopicMod) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TopicMod) GetOperation() TopicMod_Operation {
	if x != nil {
		return x.Operation
	}
	return TopicMod_NOOP
}

// A temporary representation of the topic state, e.g. was it modified to be readonly
// or is it in the process of being deleted. Once deleted the topic is permenantly gone.
type TopicTombstone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	State TopicTombstone_Status `protobuf:"varint,2,opt,name=state,proto3,enum=ensign.v1beta1.TopicTombstone_Status" json:"state,omitempty"`
}

func (x *TopicTombstone) Reset() {
	*x = TopicTombstone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopicTombstone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicTombstone) ProtoMessage() {}

func (x *TopicTombstone) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicTombstone.ProtoReflect.Descriptor instead.
func (*TopicTombstone) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{3}
}

func (x *TopicTombstone) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TopicTombstone) GetState() TopicTombstone_Status {
	if x != nil {
		return x.State
	}
	return TopicTombstone_UNKNOWN
}

var File_ensign_v1beta1_topic_proto protoreflect.FileDescriptor

var file_ensign_v1beta1_topic_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1a, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x01, 0x0a, 0x05, 0x54, 0x6f,
	0x70, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x07, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x72, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x64, 0x6f, 0x6e, 0x6c,
	0x79, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x6f, 0x6e, 0x6c,
	0x79, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x63, 0x0a, 0x0a, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x73, 0x50, 0x61, 0x67, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x06, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e,
	0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x8d, 0x01, 0x0a,
	0x08, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x4d, 0x6f, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x40, 0x0a, 0x09, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x65,
	0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x6f,
	0x70, 0x69, 0x63, 0x4d, 0x6f, 0x64, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2f, 0x0a, 0x09, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4f, 0x50,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x52, 0x43, 0x48, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12,
	0x0b, 0x0a, 0x07, 0x44, 0x45, 0x53, 0x54, 0x52, 0x4f, 0x59, 0x10, 0x02, 0x22, 0x90, 0x01, 0x0a,
	0x0e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x54, 0x6f, 0x6d, 0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x3b, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25,
	0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e,
	0x54, 0x6f, 0x70, 0x69, 0x63, 0x54, 0x6f, 0x6d, 0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x31, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x41, 0x44, 0x4f, 0x4e, 0x4c, 0x59, 0x10,
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x42,
	0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x65, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ensign_v1beta1_topic_proto_rawDescOnce sync.Once
	file_ensign_v1beta1_topic_proto_rawDescData = file_ensign_v1beta1_topic_proto_rawDesc
)

func file_ensign_v1beta1_topic_proto_rawDescGZIP() []byte {
	file_ensign_v1beta1_topic_proto_rawDescOnce.Do(func() {
		file_ensign_v1beta1_topic_proto_rawDescData = protoimpl.X.CompressGZIP(file_ensign_v1beta1_topic_proto_rawDescData)
	})
	return file_ensign_v1beta1_topic_proto_rawDescData
}

var file_ensign_v1beta1_topic_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_ensign_v1beta1_topic_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_ensign_v1beta1_topic_proto_goTypes = []interface{}{
	(TopicMod_Operation)(0),       // 0: ensign.v1beta1.TopicMod.Operation
	(TopicTombstone_Status)(0),    // 1: ensign.v1beta1.TopicTombstone.Status
	(*Topic)(nil),                 // 2: ensign.v1beta1.Topic
	(*TopicsPage)(nil),            // 3: ensign.v1beta1.TopicsPage
	(*TopicMod)(nil),              // 4: ensign.v1beta1.TopicMod
	(*TopicTombstone)(nil),        // 5: ensign.v1beta1.TopicTombstone
	(*Type)(nil),                  // 6: ensign.v1beta1.Type
	(*Region)(nil),                // 7: ensign.v1beta1.Region
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_ensign_v1beta1_topic_proto_depIdxs = []int32{
	6, // 0: ensign.v1beta1.Topic.types:type_name -> ensign.v1beta1.Type
	7, // 1: ensign.v1beta1.Topic.regions:type_name -> ensign.v1beta1.Region
	8, // 2: ensign.v1beta1.Topic.created:type_name -> google.protobuf.Timestamp
	2, // 3: ensign.v1beta1.TopicsPage.topics:type_name -> ensign.v1beta1.Topic
	0, // 4: ensign.v1beta1.TopicMod.operation:type_name -> ensign.v1beta1.TopicMod.Operation
	1, // 5: ensign.v1beta1.TopicTombstone.state:type_name -> ensign.v1beta1.TopicTombstone.Status
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_ensign_v1beta1_topic_proto_init() }
func file_ensign_v1beta1_topic_proto_init() {
	if File_ensign_v1beta1_topic_proto != nil {
		return
	}
	file_ensign_v1beta1_event_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ensign_v1beta1_topic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Topic); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ensign_v1beta1_topic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicsPage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ensign_v1beta1_topic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicMod); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ensign_v1beta1_topic_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopicTombstone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ensign_v1beta1_topic_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ensign_v1beta1_topic_proto_goTypes,
		DependencyIndexes: file_ensign_v1beta1_topic_proto_depIdxs,
		EnumInfos:         file_ensign_v1beta1_topic_proto_enumTypes,
		MessageInfos:      file_ensign_v1beta1_topic_proto_msgTypes,
	}.Build()
	File_ensign_v1beta1_topic_proto = out.File
	file_ensign_v1beta1_topic_proto_rawDesc = nil
	file_ensign_v1beta1_topic_proto_goTypes = nil
	file_ensign_v1beta1_topic_proto_depIdxs = nil
}
