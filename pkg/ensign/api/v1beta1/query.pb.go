// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: api/v1beta1/query.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Query represents a single EnSQL query with associated placeholder parameters.
type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query  string       `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Params []*Parameter `protobuf:"bytes,2,rep,name=params,proto3" json:"params,omitempty"`
}

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1beta1_query_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1beta1_query_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_api_v1beta1_query_proto_rawDescGZIP(), []int{0}
}

func (x *Query) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *Query) GetParams() []*Parameter {
	if x != nil {
		return x.Params
	}
	return nil
}

// Parameter holds a primitive value for passing as a placeholder to a sqlite query.
type Parameter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*Parameter_I
	//	*Parameter_D
	//	*Parameter_B
	//	*Parameter_Y
	//	*Parameter_S
	Value isParameter_Value `protobuf_oneof:"value"`
	Name  string            `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Parameter) Reset() {
	*x = Parameter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1beta1_query_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Parameter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Parameter) ProtoMessage() {}

func (x *Parameter) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1beta1_query_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Parameter.ProtoReflect.Descriptor instead.
func (*Parameter) Descriptor() ([]byte, []int) {
	return file_api_v1beta1_query_proto_rawDescGZIP(), []int{1}
}

func (m *Parameter) GetValue() isParameter_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Parameter) GetI() int64 {
	if x, ok := x.GetValue().(*Parameter_I); ok {
		return x.I
	}
	return 0
}

func (x *Parameter) GetD() float64 {
	if x, ok := x.GetValue().(*Parameter_D); ok {
		return x.D
	}
	return 0
}

func (x *Parameter) GetB() bool {
	if x, ok := x.GetValue().(*Parameter_B); ok {
		return x.B
	}
	return false
}

func (x *Parameter) GetY() []byte {
	if x, ok := x.GetValue().(*Parameter_Y); ok {
		return x.Y
	}
	return nil
}

func (x *Parameter) GetS() string {
	if x, ok := x.GetValue().(*Parameter_S); ok {
		return x.S
	}
	return ""
}

func (x *Parameter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type isParameter_Value interface {
	isParameter_Value()
}

type Parameter_I struct {
	I int64 `protobuf:"zigzag64,1,opt,name=i,proto3,oneof"`
}

type Parameter_D struct {
	D float64 `protobuf:"fixed64,2,opt,name=d,proto3,oneof"`
}

type Parameter_B struct {
	B bool `protobuf:"varint,3,opt,name=b,proto3,oneof"`
}

type Parameter_Y struct {
	Y []byte `protobuf:"bytes,4,opt,name=y,proto3,oneof"`
}

type Parameter_S struct {
	S string `protobuf:"bytes,5,opt,name=s,proto3,oneof"`
}

func (*Parameter_I) isParameter_Value() {}

func (*Parameter_D) isParameter_Value() {}

func (*Parameter_B) isParameter_Value() {}

func (*Parameter_Y) isParameter_Value() {}

func (*Parameter_S) isParameter_Value() {}

// Explanation returns information about the plan for executing a query and approximate
// results or errors that might be returned.
type QueryExplanation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueryExplanation) Reset() {
	*x = QueryExplanation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1beta1_query_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryExplanation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryExplanation) ProtoMessage() {}

func (x *QueryExplanation) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1beta1_query_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryExplanation.ProtoReflect.Descriptor instead.
func (*QueryExplanation) Descriptor() ([]byte, []int) {
	return file_api_v1beta1_query_proto_rawDescGZIP(), []int{2}
}

var File_api_v1beta1_query_proto protoreflect.FileDescriptor

var file_api_v1beta1_query_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x22, 0x50, 0x0a, 0x05, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x31, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x78, 0x0a, 0x09, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x01, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x12, 0x48, 0x00, 0x52, 0x01, 0x69, 0x12, 0x0e, 0x0a, 0x01, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x01, 0x64, 0x12, 0x0e, 0x0a, 0x01, 0x62, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x01, 0x62, 0x12, 0x0e, 0x0a, 0x01, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x01, 0x79, 0x12, 0x0e, 0x0a, 0x01, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x01, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x07, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x51, 0x75, 0x65, 0x72, 0x79, 0x45, 0x78,
	0x70, 0x6c, 0x61, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_v1beta1_query_proto_rawDescOnce sync.Once
	file_api_v1beta1_query_proto_rawDescData = file_api_v1beta1_query_proto_rawDesc
)

func file_api_v1beta1_query_proto_rawDescGZIP() []byte {
	file_api_v1beta1_query_proto_rawDescOnce.Do(func() {
		file_api_v1beta1_query_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1beta1_query_proto_rawDescData)
	})
	return file_api_v1beta1_query_proto_rawDescData
}

var file_api_v1beta1_query_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_v1beta1_query_proto_goTypes = []interface{}{
	(*Query)(nil),            // 0: ensign.v1beta1.Query
	(*Parameter)(nil),        // 1: ensign.v1beta1.Parameter
	(*QueryExplanation)(nil), // 2: ensign.v1beta1.QueryExplanation
}
var file_api_v1beta1_query_proto_depIdxs = []int32{
	1, // 0: ensign.v1beta1.Query.params:type_name -> ensign.v1beta1.Parameter
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_v1beta1_query_proto_init() }
func file_api_v1beta1_query_proto_init() {
	if File_api_v1beta1_query_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1beta1_query_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Query); i {
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
		file_api_v1beta1_query_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Parameter); i {
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
		file_api_v1beta1_query_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryExplanation); i {
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
	file_api_v1beta1_query_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Parameter_I)(nil),
		(*Parameter_D)(nil),
		(*Parameter_B)(nil),
		(*Parameter_Y)(nil),
		(*Parameter_S)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1beta1_query_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1beta1_query_proto_goTypes,
		DependencyIndexes: file_api_v1beta1_query_proto_depIdxs,
		MessageInfos:      file_api_v1beta1_query_proto_msgTypes,
	}.Build()
	File_api_v1beta1_query_proto = out.File
	file_api_v1beta1_query_proto_rawDesc = nil
	file_api_v1beta1_query_proto_goTypes = nil
	file_api_v1beta1_query_proto_depIdxs = nil
}