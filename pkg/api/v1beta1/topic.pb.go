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

type ShardingStrategy int32

const (
	ShardingStrategy_UNKNOWN             ShardingStrategy = 0
	ShardingStrategy_NO_SHARDING         ShardingStrategy = 1
	ShardingStrategy_CONSISTENT_KEY_HASH ShardingStrategy = 2
	ShardingStrategy_RANDOM              ShardingStrategy = 3
	ShardingStrategy_PUBLISHER_ORDERING  ShardingStrategy = 4
)

// Enum value maps for ShardingStrategy.
var (
	ShardingStrategy_name = map[int32]string{
		0: "UNKNOWN",
		1: "NO_SHARDING",
		2: "CONSISTENT_KEY_HASH",
		3: "RANDOM",
		4: "PUBLISHER_ORDERING",
	}
	ShardingStrategy_value = map[string]int32{
		"UNKNOWN":             0,
		"NO_SHARDING":         1,
		"CONSISTENT_KEY_HASH": 2,
		"RANDOM":              3,
		"PUBLISHER_ORDERING":  4,
	}
)

func (x ShardingStrategy) Enum() *ShardingStrategy {
	p := new(ShardingStrategy)
	*p = x
	return p
}

func (x ShardingStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ShardingStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_ensign_v1beta1_topic_proto_enumTypes[0].Descriptor()
}

func (ShardingStrategy) Type() protoreflect.EnumType {
	return &file_ensign_v1beta1_topic_proto_enumTypes[0]
}

func (x ShardingStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ShardingStrategy.Descriptor instead.
func (ShardingStrategy) EnumDescriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{0}
}

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
	return file_ensign_v1beta1_topic_proto_enumTypes[1].Descriptor()
}

func (TopicMod_Operation) Type() protoreflect.EnumType {
	return &file_ensign_v1beta1_topic_proto_enumTypes[1]
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
	return file_ensign_v1beta1_topic_proto_enumTypes[2].Descriptor()
}

func (TopicTombstone_Status) Type() protoreflect.EnumType {
	return &file_ensign_v1beta1_topic_proto_enumTypes[2]
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

	Id         []byte                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProjectId  []byte                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Name       string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Readonly   bool                   `protobuf:"varint,4,opt,name=readonly,proto3" json:"readonly,omitempty"`
	Offset     uint64                 `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
	Shards     uint32                 `protobuf:"varint,6,opt,name=shards,proto3" json:"shards,omitempty"`
	Placements []*Placement           `protobuf:"bytes,12,rep,name=placements,proto3" json:"placements,omitempty"`
	Types      []*Type                `protobuf:"bytes,13,rep,name=types,proto3" json:"types,omitempty"`
	Created    *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Modified   *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=modified,proto3" json:"modified,omitempty"`
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

func (x *Topic) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Topic) GetProjectId() []byte {
	if x != nil {
		return x.ProjectId
	}
	return nil
}

func (x *Topic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Topic) GetReadonly() bool {
	if x != nil {
		return x.Readonly
	}
	return false
}

func (x *Topic) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Topic) GetShards() uint32 {
	if x != nil {
		return x.Shards
	}
	return 0
}

func (x *Topic) GetPlacements() []*Placement {
	if x != nil {
		return x.Placements
	}
	return nil
}

func (x *Topic) GetTypes() []*Type {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *Topic) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Topic) GetModified() *timestamppb.Timestamp {
	if x != nil {
		return x.Modified
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

// Placement represents the nodes and regions a topic is assigned to for routing.
type Placement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch    uint64           `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Sharding ShardingStrategy `protobuf:"varint,2,opt,name=sharding,proto3,enum=ensign.v1beta1.ShardingStrategy" json:"sharding,omitempty"`
	Regions  []*Region        `protobuf:"bytes,3,rep,name=regions,proto3" json:"regions,omitempty"`
	Nodes    []*Node          `protobuf:"bytes,4,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *Placement) Reset() {
	*x = Placement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Placement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Placement) ProtoMessage() {}

func (x *Placement) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Placement.ProtoReflect.Descriptor instead.
func (*Placement) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{4}
}

func (x *Placement) GetEpoch() uint64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *Placement) GetSharding() ShardingStrategy {
	if x != nil {
		return x.Sharding
	}
	return ShardingStrategy_UNKNOWN
}

func (x *Placement) GetRegions() []*Region {
	if x != nil {
		return x.Regions
	}
	return nil
}

func (x *Placement) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Hostname string  `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Quorum   uint64  `protobuf:"varint,3,opt,name=quorum,proto3" json:"quorum,omitempty"`
	Shard    uint64  `protobuf:"varint,4,opt,name=shard,proto3" json:"shard,omitempty"`
	Region   *Region `protobuf:"bytes,5,opt,name=region,proto3" json:"region,omitempty"`
	Url      string  `protobuf:"bytes,6,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_topic_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_topic_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_topic_proto_rawDescGZIP(), []int{5}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Node) GetQuorum() uint64 {
	if x != nil {
		return x.Quorum
	}
	return 0
}

func (x *Node) GetShard() uint64 {
	if x != nil {
		return x.Shard
	}
	return 0
}

func (x *Node) GetRegion() *Region {
	if x != nil {
		return x.Region
	}
	return nil
}

func (x *Node) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_ensign_v1beta1_topic_proto protoreflect.FileDescriptor

var file_ensign_v1beta1_topic_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1a, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x02, 0x0a, 0x05, 0x54, 0x6f,
	0x70, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x64, 0x6f, 0x6e,
	0x6c, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x6f, 0x6e,
	0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x68,
	0x61, 0x72, 0x64, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2a, 0x0a,
	0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65,
	0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x36, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6d,
	0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22, 0x63, 0x0a, 0x0a, 0x54, 0x6f, 0x70, 0x69, 0x63,
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
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x22,
	0xbd, 0x01, 0x0a, 0x09, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x65, 0x70,
	0x6f, 0x63, 0x68, 0x12, 0x3c, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x08, 0x73, 0x68, 0x61, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x30, 0x0a, 0x07, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x2a, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x22,
	0xa2, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x68, 0x61, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x2a, 0x6d, 0x0a, 0x10, 0x53, 0x68, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x4f, 0x5f, 0x53, 0x48, 0x41, 0x52,
	0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x4e, 0x53, 0x49, 0x53,
	0x54, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x45, 0x59, 0x5f, 0x48, 0x41, 0x53, 0x48, 0x10, 0x02, 0x12,
	0x0a, 0x0a, 0x06, 0x52, 0x41, 0x4e, 0x44, 0x4f, 0x4d, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x50,
	0x55, 0x42, 0x4c, 0x49, 0x53, 0x48, 0x45, 0x52, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x49, 0x4e,
	0x47, 0x10, 0x04, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x65,
	0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
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

var file_ensign_v1beta1_topic_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_ensign_v1beta1_topic_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_ensign_v1beta1_topic_proto_goTypes = []interface{}{
	(ShardingStrategy)(0),         // 0: ensign.v1beta1.ShardingStrategy
	(TopicMod_Operation)(0),       // 1: ensign.v1beta1.TopicMod.Operation
	(TopicTombstone_Status)(0),    // 2: ensign.v1beta1.TopicTombstone.Status
	(*Topic)(nil),                 // 3: ensign.v1beta1.Topic
	(*TopicsPage)(nil),            // 4: ensign.v1beta1.TopicsPage
	(*TopicMod)(nil),              // 5: ensign.v1beta1.TopicMod
	(*TopicTombstone)(nil),        // 6: ensign.v1beta1.TopicTombstone
	(*Placement)(nil),             // 7: ensign.v1beta1.Placement
	(*Node)(nil),                  // 8: ensign.v1beta1.Node
	(*Type)(nil),                  // 9: ensign.v1beta1.Type
	(*timestamppb.Timestamp)(nil), // 10: google.protobuf.Timestamp
	(*Region)(nil),                // 11: ensign.v1beta1.Region
}
var file_ensign_v1beta1_topic_proto_depIdxs = []int32{
	7,  // 0: ensign.v1beta1.Topic.placements:type_name -> ensign.v1beta1.Placement
	9,  // 1: ensign.v1beta1.Topic.types:type_name -> ensign.v1beta1.Type
	10, // 2: ensign.v1beta1.Topic.created:type_name -> google.protobuf.Timestamp
	10, // 3: ensign.v1beta1.Topic.modified:type_name -> google.protobuf.Timestamp
	3,  // 4: ensign.v1beta1.TopicsPage.topics:type_name -> ensign.v1beta1.Topic
	1,  // 5: ensign.v1beta1.TopicMod.operation:type_name -> ensign.v1beta1.TopicMod.Operation
	2,  // 6: ensign.v1beta1.TopicTombstone.state:type_name -> ensign.v1beta1.TopicTombstone.Status
	0,  // 7: ensign.v1beta1.Placement.sharding:type_name -> ensign.v1beta1.ShardingStrategy
	11, // 8: ensign.v1beta1.Placement.regions:type_name -> ensign.v1beta1.Region
	8,  // 9: ensign.v1beta1.Placement.nodes:type_name -> ensign.v1beta1.Node
	11, // 10: ensign.v1beta1.Node.region:type_name -> ensign.v1beta1.Region
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
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
		file_ensign_v1beta1_topic_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Placement); i {
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
		file_ensign_v1beta1_topic_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
			NumEnums:      3,
			NumMessages:   6,
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
