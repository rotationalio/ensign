// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: pagination/pagination.proto

package pagination

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

// Key-Index Cursors are useful for high-performance pagination that do not require
// Postgres Cursors managed by an open transaction. The cursor specifies the current
// page of results so that the next/previous pages can be calculated from the query.
// Cursors also specify an expiration so that a page token cannot be replayed forever.
// Note that Key-Index cursors require the original query to correctly order the index,
// this cursor type assumes that no ordering or filtering has been applied.
//
// The cursor object is serialized and base64 encoded to be sent as a next_page_token
// in a paginated request. Protocol buffers ensures the most compact representation.
type Cursor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The start index is the ID at the beginning of the page and is used for previous
	// page queries, whereas the end index is the last ID on the page and is used to
	// compute the next page for the query. Ensure that IDs are montonically increasing
	// such as autoincrement IDs or ULIDs (do not use UUIDs).
	StartIndex string `protobuf:"bytes,1,opt,name=start_index,json=startIndex,proto3" json:"start_index,omitempty"`
	EndIndex   string `protobuf:"bytes,2,opt,name=end_index,json=endIndex,proto3" json:"end_index,omitempty"`
	// The maximum number of results per page.
	PageSize int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The timestamp when the cursor is no longer valid.
	Expires *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=expires,proto3" json:"expires,omitempty"`
}

func (x *Cursor) Reset() {
	*x = Cursor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagination_pagination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cursor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cursor) ProtoMessage() {}

func (x *Cursor) ProtoReflect() protoreflect.Message {
	mi := &file_pagination_pagination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cursor.ProtoReflect.Descriptor instead.
func (*Cursor) Descriptor() ([]byte, []int) {
	return file_pagination_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *Cursor) GetStartIndex() string {
	if x != nil {
		return x.StartIndex
	}
	return ""
}

func (x *Cursor) GetEndIndex() string {
	if x != nil {
		return x.EndIndex
	}
	return ""
}

func (x *Cursor) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *Cursor) GetExpires() *timestamppb.Timestamp {
	if x != nil {
		return x.Expires
	}
	return nil
}

var File_pagination_pagination_proto protoreflect.FileDescriptor

var file_pagination_pagination_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x99, 0x01, 0x0a, 0x06, 0x43,
	0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x34, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69,
	0x6f, 0x2f, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x74, 0x69,
	0x6c, 0x73, 0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pagination_pagination_proto_rawDescOnce sync.Once
	file_pagination_pagination_proto_rawDescData = file_pagination_pagination_proto_rawDesc
)

func file_pagination_pagination_proto_rawDescGZIP() []byte {
	file_pagination_pagination_proto_rawDescOnce.Do(func() {
		file_pagination_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(file_pagination_pagination_proto_rawDescData)
	})
	return file_pagination_pagination_proto_rawDescData
}

var file_pagination_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pagination_pagination_proto_goTypes = []interface{}{
	(*Cursor)(nil),                // 0: pagination.Cursor
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_pagination_pagination_proto_depIdxs = []int32{
	1, // 0: pagination.Cursor.expires:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pagination_pagination_proto_init() }
func file_pagination_pagination_proto_init() {
	if File_pagination_pagination_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pagination_pagination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cursor); i {
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
			RawDescriptor: file_pagination_pagination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pagination_pagination_proto_goTypes,
		DependencyIndexes: file_pagination_pagination_proto_depIdxs,
		MessageInfos:      file_pagination_pagination_proto_msgTypes,
	}.Build()
	File_pagination_pagination_proto = out.File
	file_pagination_pagination_proto_rawDesc = nil
	file_pagination_pagination_proto_goTypes = nil
	file_pagination_pagination_proto_depIdxs = nil
}
