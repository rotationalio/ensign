// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: quarterdeck/query/v1beta1/query.proto

package query

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

// A collection of statements that can be executed independently or inside of a single
// transaction. If the transaction flag is true, then all statements are executed inside
// of a transaction and a single result returned. Otherwise all statements are executed
// independently and a result for each statement is returned.
type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transaction bool         `protobuf:"varint,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	Statements  []*Statement `protobuf:"bytes,2,rep,name=statements,proto3" json:"statements,omitempty"`
}

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[0]
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
	return file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP(), []int{0}
}

func (x *Query) GetTransaction() bool {
	if x != nil {
		return x.Transaction
	}
	return false
}

func (x *Query) GetStatements() []*Statement {
	if x != nil {
		return x.Statements
	}
	return nil
}

// A single SQL statement that is parameterized by ? placeholders along with the values
// that should be passed in a secure fashion to those placeholders.
type Statement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sql        string       `protobuf:"bytes,1,opt,name=sql,proto3" json:"sql,omitempty"`
	Parameters []*Parameter `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *Statement) Reset() {
	*x = Statement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statement) ProtoMessage() {}

func (x *Statement) ProtoReflect() protoreflect.Message {
	mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statement.ProtoReflect.Descriptor instead.
func (*Statement) Descriptor() ([]byte, []int) {
	return file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP(), []int{1}
}

func (x *Statement) GetSql() string {
	if x != nil {
		return x.Sql
	}
	return ""
}

func (x *Statement) GetParameters() []*Parameter {
	if x != nil {
		return x.Parameters
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
		mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Parameter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Parameter) ProtoMessage() {}

func (x *Parameter) ProtoReflect() protoreflect.Message {
	mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[2]
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
	return file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP(), []int{2}
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

// Result holds the results of an Exec query against the database.
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastInsertId int64  `protobuf:"varint,1,opt,name=last_insert_id,json=lastInsertId,proto3" json:"last_insert_id,omitempty"`
	RowsAffected int64  `protobuf:"varint,2,opt,name=rows_affected,json=rowsAffected,proto3" json:"rows_affected,omitempty"`
	Error        string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP(), []int{3}
}

func (x *Result) GetLastInsertId() int64 {
	if x != nil {
		return x.LastInsertId
	}
	return 0
}

func (x *Result) GetRowsAffected() int64 {
	if x != nil {
		return x.RowsAffected
	}
	return 0
}

func (x *Result) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Results returns one or more results for a query.
type Results struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []*Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *Results) Reset() {
	*x = Results{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Results) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Results) ProtoMessage() {}

func (x *Results) ProtoReflect() protoreflect.Message {
	mi := &file_quarterdeck_query_v1beta1_query_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Results.ProtoReflect.Descriptor instead.
func (*Results) Descriptor() ([]byte, []int) {
	return file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP(), []int{4}
}

func (x *Results) GetResults() []*Result {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_quarterdeck_query_v1beta1_query_proto protoreflect.FileDescriptor

var file_quarterdeck_query_v1beta1_query_proto_rawDesc = []byte{
	0x0a, 0x25, 0x71, 0x75, 0x61, 0x72, 0x74, 0x65, 0x72, 0x64, 0x65, 0x63, 0x6b, 0x2f, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x71, 0x75, 0x61, 0x72, 0x74, 0x65, 0x72,
	0x64, 0x65, 0x63, 0x6b, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x22, 0x6f, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x44, 0x0a,
	0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x71, 0x75, 0x61, 0x72, 0x74, 0x65, 0x72, 0x64, 0x65, 0x63, 0x6b, 0x2e,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x22, 0x63, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x71, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73,
	0x71, 0x6c, 0x12, 0x44, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x71, 0x75, 0x61, 0x72, 0x74, 0x65, 0x72,
	0x64, 0x65, 0x63, 0x6b, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x0a, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0x78, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x01, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12,
	0x48, 0x00, 0x52, 0x01, 0x69, 0x12, 0x0e, 0x0a, 0x01, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x00, 0x52, 0x01, 0x64, 0x12, 0x0e, 0x0a, 0x01, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x00, 0x52, 0x01, 0x62, 0x12, 0x0e, 0x0a, 0x01, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x00, 0x52, 0x01, 0x79, 0x12, 0x0e, 0x0a, 0x01, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x01, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x69, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x24, 0x0a, 0x0e,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x6f, 0x77, 0x73, 0x5f, 0x61, 0x66, 0x66, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x6f, 0x77, 0x73, 0x41,
	0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x46, 0x0a,
	0x07, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x3b, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x71, 0x75, 0x61, 0x72,
	0x74, 0x65, 0x72, 0x64, 0x65, 0x63, 0x6b, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x6f,
	0x2f, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x71, 0x75, 0x61, 0x72,
	0x74, 0x65, 0x72, 0x64, 0x65, 0x63, 0x6b, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x2f,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_quarterdeck_query_v1beta1_query_proto_rawDescOnce sync.Once
	file_quarterdeck_query_v1beta1_query_proto_rawDescData = file_quarterdeck_query_v1beta1_query_proto_rawDesc
)

func file_quarterdeck_query_v1beta1_query_proto_rawDescGZIP() []byte {
	file_quarterdeck_query_v1beta1_query_proto_rawDescOnce.Do(func() {
		file_quarterdeck_query_v1beta1_query_proto_rawDescData = protoimpl.X.CompressGZIP(file_quarterdeck_query_v1beta1_query_proto_rawDescData)
	})
	return file_quarterdeck_query_v1beta1_query_proto_rawDescData
}

var file_quarterdeck_query_v1beta1_query_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_quarterdeck_query_v1beta1_query_proto_goTypes = []interface{}{
	(*Query)(nil),     // 0: quarterdeck.query.v1beta1.Query
	(*Statement)(nil), // 1: quarterdeck.query.v1beta1.Statement
	(*Parameter)(nil), // 2: quarterdeck.query.v1beta1.Parameter
	(*Result)(nil),    // 3: quarterdeck.query.v1beta1.Result
	(*Results)(nil),   // 4: quarterdeck.query.v1beta1.Results
}
var file_quarterdeck_query_v1beta1_query_proto_depIdxs = []int32{
	1, // 0: quarterdeck.query.v1beta1.Query.statements:type_name -> quarterdeck.query.v1beta1.Statement
	2, // 1: quarterdeck.query.v1beta1.Statement.parameters:type_name -> quarterdeck.query.v1beta1.Parameter
	3, // 2: quarterdeck.query.v1beta1.Results.results:type_name -> quarterdeck.query.v1beta1.Result
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_quarterdeck_query_v1beta1_query_proto_init() }
func file_quarterdeck_query_v1beta1_query_proto_init() {
	if File_quarterdeck_query_v1beta1_query_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_quarterdeck_query_v1beta1_query_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_quarterdeck_query_v1beta1_query_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statement); i {
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
		file_quarterdeck_query_v1beta1_query_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_quarterdeck_query_v1beta1_query_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
		file_quarterdeck_query_v1beta1_query_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Results); i {
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
	file_quarterdeck_query_v1beta1_query_proto_msgTypes[2].OneofWrappers = []interface{}{
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
			RawDescriptor: file_quarterdeck_query_v1beta1_query_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_quarterdeck_query_v1beta1_query_proto_goTypes,
		DependencyIndexes: file_quarterdeck_query_v1beta1_query_proto_depIdxs,
		MessageInfos:      file_quarterdeck_query_v1beta1_query_proto_msgTypes,
	}.Build()
	File_quarterdeck_query_v1beta1_query_proto = out.File
	file_quarterdeck_query_v1beta1_query_proto_rawDesc = nil
	file_quarterdeck_query_v1beta1_query_proto_goTypes = nil
	file_quarterdeck_query_v1beta1_query_proto_depIdxs = nil
}
