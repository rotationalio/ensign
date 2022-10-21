// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: raft/v1beta1/log.proto

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

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Term  uint64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Key   []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_log_proto_rawDescGZIP(), []int{0}
}

func (x *LogEntry) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LogEntry) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *LogEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type LogMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastApplied uint64                 `protobuf:"varint,1,opt,name=last_applied,json=lastApplied,proto3" json:"last_applied,omitempty"`
	CommitIndex uint64                 `protobuf:"varint,2,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"`
	Length      uint64                 `protobuf:"varint,3,opt,name=length,proto3" json:"length,omitempty"`
	Created     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created,proto3" json:"created,omitempty"`
	Modifed     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=modifed,proto3" json:"modifed,omitempty"`
	Snapshot    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
}

func (x *LogMeta) Reset() {
	*x = LogMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogMeta) ProtoMessage() {}

func (x *LogMeta) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogMeta.ProtoReflect.Descriptor instead.
func (*LogMeta) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_log_proto_rawDescGZIP(), []int{1}
}

func (x *LogMeta) GetLastApplied() uint64 {
	if x != nil {
		return x.LastApplied
	}
	return 0
}

func (x *LogMeta) GetCommitIndex() uint64 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

func (x *LogMeta) GetLength() uint64 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *LogMeta) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *LogMeta) GetModifed() *timestamppb.Timestamp {
	if x != nil {
		return x.Modifed
	}
	return nil
}

func (x *LogMeta) GetSnapshot() *timestamppb.Timestamp {
	if x != nil {
		return x.Snapshot
	}
	return nil
}

var File_raft_v1beta1_log_proto protoreflect.FileDescriptor

var file_raft_v1beta1_log_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6c,
	0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x8b, 0x02, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x74,
	0x61, 0x12, 0x21, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x41, 0x70, 0x70,
	0x6c, 0x69, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12,
	0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x65, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x07, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x08, 0x73,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x65,
	0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raft_v1beta1_log_proto_rawDescOnce sync.Once
	file_raft_v1beta1_log_proto_rawDescData = file_raft_v1beta1_log_proto_rawDesc
)

func file_raft_v1beta1_log_proto_rawDescGZIP() []byte {
	file_raft_v1beta1_log_proto_rawDescOnce.Do(func() {
		file_raft_v1beta1_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_v1beta1_log_proto_rawDescData)
	})
	return file_raft_v1beta1_log_proto_rawDescData
}

var file_raft_v1beta1_log_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_raft_v1beta1_log_proto_goTypes = []interface{}{
	(*LogEntry)(nil),              // 0: raft.v1beta1.LogEntry
	(*LogMeta)(nil),               // 1: raft.v1beta1.LogMeta
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_raft_v1beta1_log_proto_depIdxs = []int32{
	2, // 0: raft.v1beta1.LogMeta.created:type_name -> google.protobuf.Timestamp
	2, // 1: raft.v1beta1.LogMeta.modifed:type_name -> google.protobuf.Timestamp
	2, // 2: raft.v1beta1.LogMeta.snapshot:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_raft_v1beta1_log_proto_init() }
func file_raft_v1beta1_log_proto_init() {
	if File_raft_v1beta1_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raft_v1beta1_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
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
		file_raft_v1beta1_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogMeta); i {
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
			RawDescriptor: file_raft_v1beta1_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_raft_v1beta1_log_proto_goTypes,
		DependencyIndexes: file_raft_v1beta1_log_proto_depIdxs,
		MessageInfos:      file_raft_v1beta1_log_proto_msgTypes,
	}.Build()
	File_raft_v1beta1_log_proto = out.File
	file_raft_v1beta1_log_proto_rawDesc = nil
	file_raft_v1beta1_log_proto_goTypes = nil
	file_raft_v1beta1_log_proto_depIdxs = nil
}
