// Code generated by protoc-gen-go. DO NOT EDIT.
// source: meta.proto

package main

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ValueType int32

const (
	ValueType_Double    ValueType = 0
	ValueType_Int64     ValueType = 1
	ValueType_String    ValueType = 2
	ValueType_Timestamp ValueType = 3
	ValueType_Boolean   ValueType = 4
)

var ValueType_name = map[int32]string{
	0: "Double",
	1: "Int64",
	2: "String",
	3: "Timestamp",
	4: "Boolean",
}

var ValueType_value = map[string]int32{
	"Double":    0,
	"Int64":     1,
	"String":    2,
	"Timestamp": 3,
	"Boolean":   4,
}

func (x ValueType) String() string {
	return proto.EnumName(ValueType_name, int32(x))
}

func (ValueType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{0}
}

type Value struct {
	Type                 ValueType `protobuf:"varint,1,opt,name=type,proto3,enum=meta.ValueType" json:"type,omitempty"`
	DoubleValue          float64   `protobuf:"fixed64,2,opt,name=doubleValue,proto3" json:"doubleValue,omitempty"`
	Int64Value           int64     `protobuf:"varint,3,opt,name=int64Value,proto3" json:"int64Value,omitempty"`
	StringValue          string    `protobuf:"bytes,4,opt,name=stringValue,proto3" json:"stringValue,omitempty"`
	TimestampValue       []byte    `protobuf:"bytes,5,opt,name=timestampValue,proto3" json:"timestampValue,omitempty"`
	BooleanValue         bool      `protobuf:"varint,6,opt,name=booleanValue,proto3" json:"booleanValue,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{0}
}

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetType() ValueType {
	if m != nil {
		return m.Type
	}
	return ValueType_Double
}

func (m *Value) GetDoubleValue() float64 {
	if m != nil {
		return m.DoubleValue
	}
	return 0
}

func (m *Value) GetInt64Value() int64 {
	if m != nil {
		return m.Int64Value
	}
	return 0
}

func (m *Value) GetStringValue() string {
	if m != nil {
		return m.StringValue
	}
	return ""
}

func (m *Value) GetTimestampValue() []byte {
	if m != nil {
		return m.TimestampValue
	}
	return nil
}

func (m *Value) GetBooleanValue() bool {
	if m != nil {
		return m.BooleanValue
	}
	return false
}

type Field struct {
	Id                   int32     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type                 ValueType `protobuf:"varint,3,opt,name=type,proto3,enum=meta.ValueType" json:"type,omitempty"`
	Comment              string    `protobuf:"bytes,4,opt,name=comment,proto3" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Field) Reset()         { *m = Field{} }
func (m *Field) String() string { return proto.CompactTextString(m) }
func (*Field) ProtoMessage()    {}
func (*Field) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{1}
}

func (m *Field) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Field.Unmarshal(m, b)
}
func (m *Field) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Field.Marshal(b, m, deterministic)
}
func (m *Field) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Field.Merge(m, src)
}
func (m *Field) XXX_Size() int {
	return xxx_messageInfo_Field.Size(m)
}
func (m *Field) XXX_DiscardUnknown() {
	xxx_messageInfo_Field.DiscardUnknown(m)
}

var xxx_messageInfo_Field proto.InternalMessageInfo

func (m *Field) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Field) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Field) GetType() ValueType {
	if m != nil {
		return m.Type
	}
	return ValueType_Double
}

func (m *Field) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

type Kind struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Fields               []*Field `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Kind) Reset()         { *m = Kind{} }
func (m *Kind) String() string { return proto.CompactTextString(m) }
func (*Kind) ProtoMessage()    {}
func (*Kind) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{2}
}

func (m *Kind) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Kind.Unmarshal(m, b)
}
func (m *Kind) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Kind.Marshal(b, m, deterministic)
}
func (m *Kind) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Kind.Merge(m, src)
}
func (m *Kind) XXX_Size() int {
	return xxx_messageInfo_Kind.Size(m)
}
func (m *Kind) XXX_DiscardUnknown() {
	xxx_messageInfo_Kind.DiscardUnknown(m)
}

var xxx_messageInfo_Kind proto.InternalMessageInfo

func (m *Kind) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Kind) GetFields() []*Field {
	if m != nil {
		return m.Fields
	}
	return nil
}

type Schema struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Kinds                []*Kind  `protobuf:"bytes,2,rep,name=kinds,proto3" json:"kinds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Schema) Reset()         { *m = Schema{} }
func (m *Schema) String() string { return proto.CompactTextString(m) }
func (*Schema) ProtoMessage()    {}
func (*Schema) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{3}
}

func (m *Schema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Schema.Unmarshal(m, b)
}
func (m *Schema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Schema.Marshal(b, m, deterministic)
}
func (m *Schema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Schema.Merge(m, src)
}
func (m *Schema) XXX_Size() int {
	return xxx_messageInfo_Schema.Size(m)
}
func (m *Schema) XXX_DiscardUnknown() {
	xxx_messageInfo_Schema.DiscardUnknown(m)
}

var xxx_messageInfo_Schema proto.InternalMessageInfo

func (m *Schema) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Schema) GetKinds() []*Kind {
	if m != nil {
		return m.Kinds
	}
	return nil
}

type GetSchemaRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSchemaRequest) Reset()         { *m = GetSchemaRequest{} }
func (m *GetSchemaRequest) String() string { return proto.CompactTextString(m) }
func (*GetSchemaRequest) ProtoMessage()    {}
func (*GetSchemaRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{4}
}

func (m *GetSchemaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSchemaRequest.Unmarshal(m, b)
}
func (m *GetSchemaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSchemaRequest.Marshal(b, m, deterministic)
}
func (m *GetSchemaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSchemaRequest.Merge(m, src)
}
func (m *GetSchemaRequest) XXX_Size() int {
	return xxx_messageInfo_GetSchemaRequest.Size(m)
}
func (m *GetSchemaRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSchemaRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSchemaRequest proto.InternalMessageInfo

type GetSchemaResponse struct {
	Schema               *Schema  `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSchemaResponse) Reset()         { *m = GetSchemaResponse{} }
func (m *GetSchemaResponse) String() string { return proto.CompactTextString(m) }
func (*GetSchemaResponse) ProtoMessage()    {}
func (*GetSchemaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{5}
}

func (m *GetSchemaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSchemaResponse.Unmarshal(m, b)
}
func (m *GetSchemaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSchemaResponse.Marshal(b, m, deterministic)
}
func (m *GetSchemaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSchemaResponse.Merge(m, src)
}
func (m *GetSchemaResponse) XXX_Size() int {
	return xxx_messageInfo_GetSchemaResponse.Size(m)
}
func (m *GetSchemaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSchemaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSchemaResponse proto.InternalMessageInfo

func (m *GetSchemaResponse) GetSchema() *Schema {
	if m != nil {
		return m.Schema
	}
	return nil
}

func init() {
	proto.RegisterEnum("meta.ValueType", ValueType_name, ValueType_value)
	proto.RegisterType((*Value)(nil), "meta.Value")
	proto.RegisterType((*Field)(nil), "meta.Field")
	proto.RegisterType((*Kind)(nil), "meta.Kind")
	proto.RegisterType((*Schema)(nil), "meta.Schema")
	proto.RegisterType((*GetSchemaRequest)(nil), "meta.GetSchemaRequest")
	proto.RegisterType((*GetSchemaResponse)(nil), "meta.GetSchemaResponse")
}

func init() { proto.RegisterFile("meta.proto", fileDescriptor_3b5ea8fe65782bcc) }

var fileDescriptor_3b5ea8fe65782bcc = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4d, 0x6f, 0xd4, 0x30,
	0x14, 0xc4, 0xf9, 0x2a, 0x79, 0x59, 0x96, 0xf0, 0x0e, 0x65, 0xc5, 0x01, 0x59, 0x01, 0xa1, 0x88,
	0x43, 0x0f, 0x0b, 0x42, 0x42, 0x42, 0x20, 0x01, 0x02, 0x01, 0xe2, 0xe2, 0xad, 0x7a, 0xcf, 0x6e,
	0x5e, 0x8b, 0x21, 0xb1, 0x43, 0xec, 0x45, 0xea, 0x7f, 0xe5, 0xc7, 0xa0, 0xd8, 0x69, 0x48, 0x17,
	0xd4, 0x9b, 0x33, 0x33, 0x1a, 0xcf, 0x78, 0x02, 0xd0, 0x92, 0xad, 0x4e, 0xba, 0x5e, 0x5b, 0x8d,
	0xd1, 0x70, 0x2e, 0x7e, 0x33, 0x88, 0xcf, 0xaa, 0x66, 0x4f, 0xf8, 0x08, 0x22, 0x7b, 0xd9, 0xd1,
	0x8a, 0x71, 0x56, 0x2e, 0xd7, 0x77, 0x4f, 0x9c, 0xd4, 0x51, 0xa7, 0x97, 0x1d, 0x09, 0x47, 0x22,
	0x87, 0xac, 0xd6, 0xfb, 0x6d, 0x43, 0x8e, 0x58, 0x05, 0x9c, 0x95, 0x4c, 0xcc, 0x21, 0x7c, 0x08,
	0x20, 0x95, 0x7d, 0xf1, 0xdc, 0x0b, 0x42, 0xce, 0xca, 0x50, 0xcc, 0x90, 0xc1, 0xc1, 0xd8, 0x5e,
	0xaa, 0x0b, 0x2f, 0x88, 0x38, 0x2b, 0x53, 0x31, 0x87, 0xf0, 0x09, 0x2c, 0xad, 0x6c, 0xc9, 0xd8,
	0xaa, 0xed, 0xbc, 0x28, 0xe6, 0xac, 0x5c, 0x88, 0x03, 0x14, 0x0b, 0x58, 0x6c, 0xb5, 0x6e, 0xa8,
	0x52, 0x5e, 0x95, 0x70, 0x56, 0xde, 0x16, 0xd7, 0xb0, 0xe2, 0x3b, 0xc4, 0x1f, 0x24, 0x35, 0x35,
	0x2e, 0x21, 0x90, 0xb5, 0xeb, 0x16, 0x8b, 0x40, 0xd6, 0x88, 0x10, 0xa9, 0xaa, 0xf5, 0x0d, 0x52,
	0xe1, 0xce, 0xd3, 0x0b, 0x84, 0x37, 0xbd, 0xc0, 0x0a, 0x8e, 0x76, 0xba, 0x6d, 0x49, 0xd9, 0x31,
	0xfb, 0xd5, 0x67, 0xf1, 0x06, 0xa2, 0x2f, 0x52, 0xfd, 0xb5, 0x66, 0xd7, 0xac, 0x93, 0xf3, 0x21,
	0x87, 0x59, 0x05, 0x3c, 0x2c, 0xb3, 0x75, 0xe6, 0xcd, 0x5d, 0x36, 0x31, 0x52, 0xc5, 0x6b, 0x48,
	0x36, 0xbb, 0x6f, 0xd4, 0x56, 0xff, 0xb5, 0xe0, 0x10, 0xff, 0x90, 0x6a, 0x72, 0x00, 0xef, 0x30,
	0xdc, 0x28, 0x3c, 0x51, 0x20, 0xe4, 0x1f, 0xc9, 0x7a, 0x0b, 0x41, 0x3f, 0xf7, 0x64, 0x6c, 0xf1,
	0x12, 0xee, 0xcd, 0x30, 0xd3, 0x69, 0x65, 0x08, 0x1f, 0x43, 0x62, 0x1c, 0xe2, 0x2e, 0xc8, 0xd6,
	0x0b, 0xef, 0x35, 0xaa, 0x46, 0xee, 0xe9, 0x67, 0x48, 0xa7, 0xf2, 0x08, 0x90, 0xbc, 0x77, 0x2b,
	0xe7, 0xb7, 0x30, 0x85, 0xf8, 0xd3, 0x30, 0x68, 0xce, 0x06, 0x78, 0xe3, 0xa6, 0xcb, 0x03, 0xbc,
	0x03, 0xe9, 0xe9, 0xd5, 0x42, 0x79, 0x88, 0x19, 0x1c, 0xbd, 0xf5, 0x53, 0xe4, 0xd1, 0xfa, 0x0c,
	0x8e, 0xdf, 0x69, 0x75, 0x2e, 0x2f, 0x8c, 0xd5, 0x3d, 0x7d, 0x25, 0x5b, 0x6d, 0xa8, 0xff, 0x25,
	0x77, 0x84, 0xaf, 0x20, 0x9d, 0x02, 0xe2, 0xb1, 0x0f, 0x72, 0xd8, 0xe2, 0xc1, 0xfd, 0x7f, 0x70,
	0xdf, 0x64, 0x9b, 0xb8, 0x7f, 0xf9, 0xd9, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbb, 0x80, 0xb5,
	0x4a, 0xd9, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConfigstoreMetaServiceClient is the client API for ConfigstoreMetaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigstoreMetaServiceClient interface {
	GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error)
}

type configstoreMetaServiceClient struct {
	cc *grpc.ClientConn
}

func NewConfigstoreMetaServiceClient(cc *grpc.ClientConn) ConfigstoreMetaServiceClient {
	return &configstoreMetaServiceClient{cc}
}

func (c *configstoreMetaServiceClient) GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error) {
	out := new(GetSchemaResponse)
	err := c.cc.Invoke(ctx, "/meta.ConfigstoreMetaService/GetSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigstoreMetaServiceServer is the server API for ConfigstoreMetaService service.
type ConfigstoreMetaServiceServer interface {
	GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error)
}

func RegisterConfigstoreMetaServiceServer(s *grpc.Server, srv ConfigstoreMetaServiceServer) {
	s.RegisterService(&_ConfigstoreMetaService_serviceDesc, srv)
}

func _ConfigstoreMetaService_GetSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigstoreMetaServiceServer).GetSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.ConfigstoreMetaService/GetSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigstoreMetaServiceServer).GetSchema(ctx, req.(*GetSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConfigstoreMetaService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "meta.ConfigstoreMetaService",
	HandlerType: (*ConfigstoreMetaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSchema",
			Handler:    _ConfigstoreMetaService_GetSchema_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "meta.proto",
}

