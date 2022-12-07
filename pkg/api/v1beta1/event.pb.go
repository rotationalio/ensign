// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: ensign/v1beta1/event.proto

package api

import (
	v1beta1 "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
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

// Event is a high level wrapper for a datagram that is totally ordered by the Ensign
// event-driven framework. Events are simply blobs of data and associated metadata that
// can be published by a producer, inserted into a log, and consumed by a subscriber.
// The mimetype of the event allows subscribers to deserialize the data into a specific
// format such as JSON or protocol buffers. The type acts as a key for heterogeneous
// topics and can also be used to lookup schema information for data validation.
// TODO: do we need to allow for event keys or is the type sufficient?
// TODO: how should we implement the event IDs, should we use a time based mechanism like ksuid?
// TODO: is this too nested? should we flatten some of the inner types?
// TODO: do we need generic metadata?
// TODO: what about offset and epoch information?
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TopicId       string                 `protobuf:"bytes,2,opt,name=topic_id,json=topicId,proto3" json:"topic_id,omitempty"`
	Mimetype      v1beta1.MIME           `protobuf:"varint,3,opt,name=mimetype,proto3,enum=mimetype.v1beta1.MIME" json:"mimetype,omitempty"`
	Type          *Type                  `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Key           []byte                 `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	Data          []byte                 `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	Encryption    *Encryption            `protobuf:"bytes,7,opt,name=encryption,proto3" json:"encryption,omitempty"`
	Compression   *Compression           `protobuf:"bytes,8,opt,name=compression,proto3" json:"compression,omitempty"`
	Geography     *Region                `protobuf:"bytes,9,opt,name=geography,proto3" json:"geography,omitempty"`
	Publisher     *Publisher             `protobuf:"bytes,10,opt,name=publisher,proto3" json:"publisher,omitempty"`
	UserDefinedId string                 `protobuf:"bytes,11,opt,name=user_defined_id,json=userDefinedId,proto3" json:"user_defined_id,omitempty"`
	Created       *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Committed     *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=committed,proto3" json:"committed,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetTopicId() string {
	if x != nil {
		return x.TopicId
	}
	return ""
}

func (x *Event) GetMimetype() v1beta1.MIME {
	if x != nil {
		return x.Mimetype
	}
	return v1beta1.MIME(0)
}

func (x *Event) GetType() *Type {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *Event) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Event) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Event) GetEncryption() *Encryption {
	if x != nil {
		return x.Encryption
	}
	return nil
}

func (x *Event) GetCompression() *Compression {
	if x != nil {
		return x.Compression
	}
	return nil
}

func (x *Event) GetGeography() *Region {
	if x != nil {
		return x.Geography
	}
	return nil
}

func (x *Event) GetPublisher() *Publisher {
	if x != nil {
		return x.Publisher
	}
	return nil
}

func (x *Event) GetUserDefinedId() string {
	if x != nil {
		return x.UserDefinedId
	}
	return ""
}

func (x *Event) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Event) GetCommitted() *timestamppb.Timestamp {
	if x != nil {
		return x.Committed
	}
	return nil
}

// An event type is composed of a name and a version so that the type can be looked up
// in the schema registry. The schema can then be used to validate the data inside the
// event. Schemas are optional but types are not unless the mimetype requries a schema
// for deserialization (e.g. protobuf, parquet, avro, etc.).
type Type struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version uint32 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Type) Reset() {
	*x = Type{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type) ProtoMessage() {}

func (x *Type) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type.ProtoReflect.Descriptor instead.
func (*Type) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{1}
}

func (x *Type) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Type) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

// Metadata about the cryptography used to secure the event.
// TODO: should we encrypt each event individually or blocks of events together?
// TODO: this is only partially implemented
type Encryption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm string `protobuf:"bytes,1,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
	KeyId     string `protobuf:"bytes,2,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
}

func (x *Encryption) Reset() {
	*x = Encryption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Encryption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Encryption) ProtoMessage() {}

func (x *Encryption) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Encryption.ProtoReflect.Descriptor instead.
func (*Encryption) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{2}
}

func (x *Encryption) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

func (x *Encryption) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

// Metadata about compression used to reduce the storage size of the event.
// TODO: should we compress each event individually or blocks of events together?
// TODO: this is only partially implemented
type Compression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm string `protobuf:"bytes,2,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
}

func (x *Compression) Reset() {
	*x = Compression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Compression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Compression) ProtoMessage() {}

func (x *Compression) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Compression.ProtoReflect.Descriptor instead.
func (*Compression) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{3}
}

func (x *Compression) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

// Geographic metadata for compliance and region-awareness.
// TODO: this is only partially implemented
type Region struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Region) Reset() {
	*x = Region{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Region) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Region) ProtoMessage() {}

func (x *Region) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Region.ProtoReflect.Descriptor instead.
func (*Region) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{4}
}

func (x *Region) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Information about the publisher of the event for provenance and auditing purposes.
// TODO: this is only partially implemented
type Publisher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Ipaddr   string `protobuf:"bytes,2,opt,name=ipaddr,proto3" json:"ipaddr,omitempty"`
}

func (x *Publisher) Reset() {
	*x = Publisher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ensign_v1beta1_event_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publisher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publisher) ProtoMessage() {}

func (x *Publisher) ProtoReflect() protoreflect.Message {
	mi := &file_ensign_v1beta1_event_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publisher.ProtoReflect.Descriptor instead.
func (*Publisher) Descriptor() ([]byte, []int) {
	return file_ensign_v1beta1_event_proto_rawDescGZIP(), []int{5}
}

func (x *Publisher) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Publisher) GetIpaddr() string {
	if x != nil {
		return x.Ipaddr
	}
	return ""
}

var File_ensign_v1beta1_event_proto protoreflect.FileDescriptor

var file_ensign_v1beta1_event_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1f, 0x6d, 0x69,
	0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d,
	0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8,
	0x04, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6d, 0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4d, 0x49, 0x4d, 0x45, 0x52, 0x08, 0x6d,
	0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x6e, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x45, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x09, 0x67, 0x65, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x67,
	0x65, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x12, 0x37, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6e,
	0x73, 0x69, 0x67, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65,
	0x72, 0x12, 0x26, 0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x72,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x38, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x22, 0x34, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22,
	0x41, 0x0a, 0x0a, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x15, 0x0a, 0x06, 0x6b,
	0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6b, 0x65, 0x79,
	0x49, 0x64, 0x22, 0x2b, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x22,
	0x1c, 0x0a, 0x06, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x40, 0x0a,
	0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x70, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x70, 0x61, 0x64, 0x64, 0x72, 0x42,
	0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x65, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ensign_v1beta1_event_proto_rawDescOnce sync.Once
	file_ensign_v1beta1_event_proto_rawDescData = file_ensign_v1beta1_event_proto_rawDesc
)

func file_ensign_v1beta1_event_proto_rawDescGZIP() []byte {
	file_ensign_v1beta1_event_proto_rawDescOnce.Do(func() {
		file_ensign_v1beta1_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_ensign_v1beta1_event_proto_rawDescData)
	})
	return file_ensign_v1beta1_event_proto_rawDescData
}

var file_ensign_v1beta1_event_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_ensign_v1beta1_event_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: ensign.v1beta1.Event
	(*Type)(nil),                  // 1: ensign.v1beta1.Type
	(*Encryption)(nil),            // 2: ensign.v1beta1.Encryption
	(*Compression)(nil),           // 3: ensign.v1beta1.Compression
	(*Region)(nil),                // 4: ensign.v1beta1.Region
	(*Publisher)(nil),             // 5: ensign.v1beta1.Publisher
	(v1beta1.MIME)(0),             // 6: mimetype.v1beta1.MIME
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_ensign_v1beta1_event_proto_depIdxs = []int32{
	6, // 0: ensign.v1beta1.Event.mimetype:type_name -> mimetype.v1beta1.MIME
	1, // 1: ensign.v1beta1.Event.type:type_name -> ensign.v1beta1.Type
	2, // 2: ensign.v1beta1.Event.encryption:type_name -> ensign.v1beta1.Encryption
	3, // 3: ensign.v1beta1.Event.compression:type_name -> ensign.v1beta1.Compression
	4, // 4: ensign.v1beta1.Event.geography:type_name -> ensign.v1beta1.Region
	5, // 5: ensign.v1beta1.Event.publisher:type_name -> ensign.v1beta1.Publisher
	7, // 6: ensign.v1beta1.Event.created:type_name -> google.protobuf.Timestamp
	7, // 7: ensign.v1beta1.Event.committed:type_name -> google.protobuf.Timestamp
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_ensign_v1beta1_event_proto_init() }
func file_ensign_v1beta1_event_proto_init() {
	if File_ensign_v1beta1_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ensign_v1beta1_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_ensign_v1beta1_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type); i {
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
		file_ensign_v1beta1_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Encryption); i {
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
		file_ensign_v1beta1_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Compression); i {
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
		file_ensign_v1beta1_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Region); i {
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
		file_ensign_v1beta1_event_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publisher); i {
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
			RawDescriptor: file_ensign_v1beta1_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ensign_v1beta1_event_proto_goTypes,
		DependencyIndexes: file_ensign_v1beta1_event_proto_depIdxs,
		MessageInfos:      file_ensign_v1beta1_event_proto_msgTypes,
	}.Build()
	File_ensign_v1beta1_event_proto = out.File
	file_ensign_v1beta1_event_proto_rawDesc = nil
	file_ensign_v1beta1_event_proto_goTypes = nil
	file_ensign_v1beta1_event_proto_depIdxs = nil
}
