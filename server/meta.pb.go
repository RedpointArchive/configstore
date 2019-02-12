// Code generated by protoc-gen-go. DO NOT EDIT.
// source: meta.proto

package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ValueType int32

const (
	ValueType_Double    ValueType = 0
	ValueType_Int64     ValueType = 1
	ValueType_String    ValueType = 2
	ValueType_Timestamp ValueType = 3
	ValueType_Boolean   ValueType = 4
	ValueType_Bytes     ValueType = 5
	ValueType_Key_      ValueType = 6
)

var ValueType_name = map[int32]string{
	0: "Double",
	1: "Int64",
	2: "String",
	3: "Timestamp",
	4: "Boolean",
	5: "Bytes",
	6: "Key_",
}
var ValueType_value = map[string]int32{
	"Double":    0,
	"Int64":     1,
	"String":    2,
	"Timestamp": 3,
	"Boolean":   4,
	"Bytes":     5,
	"Key_":      6,
}

func (x ValueType) String() string {
	return proto.EnumName(ValueType_name, int32(x))
}
func (ValueType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{0}
}

type FieldEditorInfoType int32

const (
	FieldEditorInfoType_Default  FieldEditorInfoType = 0
	FieldEditorInfoType_Password FieldEditorInfoType = 1
	FieldEditorInfoType_Lookup   FieldEditorInfoType = 2
)

var FieldEditorInfoType_name = map[int32]string{
	0: "Default",
	1: "Password",
	2: "Lookup",
}
var FieldEditorInfoType_value = map[string]int32{
	"Default":  0,
	"Password": 1,
	"Lookup":   2,
}

func (x FieldEditorInfoType) String() string {
	return proto.EnumName(FieldEditorInfoType_name, int32(x))
}
func (FieldEditorInfoType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{1}
}

type Key struct {
	Val                  string   `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	IsSet                bool     `protobuf:"varint,2,opt,name=isSet,proto3" json:"isSet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Key) Reset()         { *m = Key{} }
func (m *Key) String() string { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()    {}
func (*Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{0}
}
func (m *Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Key.Unmarshal(m, b)
}
func (m *Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Key.Marshal(b, m, deterministic)
}
func (dst *Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Key.Merge(dst, src)
}
func (m *Key) XXX_Size() int {
	return xxx_messageInfo_Key.Size(m)
}
func (m *Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Key proto.InternalMessageInfo

func (m *Key) GetVal() string {
	if m != nil {
		return m.Val
	}
	return ""
}

func (m *Key) GetIsSet() bool {
	if m != nil {
		return m.IsSet
	}
	return false
}

type Value struct {
	Id                   int32     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 ValueType `protobuf:"varint,2,opt,name=type,proto3,enum=meta.ValueType" json:"type,omitempty"`
	DoubleValue          float64   `protobuf:"fixed64,3,opt,name=doubleValue,proto3" json:"doubleValue,omitempty"`
	Int64Value           int64     `protobuf:"varint,4,opt,name=int64Value,proto3" json:"int64Value,omitempty"`
	StringValue          string    `protobuf:"bytes,5,opt,name=stringValue,proto3" json:"stringValue,omitempty"`
	TimestampValue       []byte    `protobuf:"bytes,6,opt,name=timestampValue,proto3" json:"timestampValue,omitempty"`
	BooleanValue         bool      `protobuf:"varint,7,opt,name=booleanValue,proto3" json:"booleanValue,omitempty"`
	BytesValue           []byte    `protobuf:"bytes,8,opt,name=bytesValue,proto3" json:"bytesValue,omitempty"`
	KeyValue             *Key      `protobuf:"bytes,9,opt,name=keyValue,proto3" json:"keyValue,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{1}
}
func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (dst *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(dst, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

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

func (m *Value) GetBytesValue() []byte {
	if m != nil {
		return m.BytesValue
	}
	return nil
}

func (m *Value) GetKeyValue() *Key {
	if m != nil {
		return m.KeyValue
	}
	return nil
}

type Field struct {
	Id                   int32            `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type                 ValueType        `protobuf:"varint,3,opt,name=type,proto3,enum=meta.ValueType" json:"type,omitempty"`
	Comment              string           `protobuf:"bytes,4,opt,name=comment,proto3" json:"comment,omitempty"`
	Editor               *FieldEditorInfo `protobuf:"bytes,5,opt,name=editor,proto3" json:"editor,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Field) Reset()         { *m = Field{} }
func (m *Field) String() string { return proto.CompactTextString(m) }
func (*Field) ProtoMessage()    {}
func (*Field) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{2}
}
func (m *Field) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Field.Unmarshal(m, b)
}
func (m *Field) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Field.Marshal(b, m, deterministic)
}
func (dst *Field) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Field.Merge(dst, src)
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

func (m *Field) GetEditor() *FieldEditorInfo {
	if m != nil {
		return m.Editor
	}
	return nil
}

type FieldEditorInfo struct {
	DisplayName          string              `protobuf:"bytes,1,opt,name=displayName,proto3" json:"displayName,omitempty"`
	Type                 FieldEditorInfoType `protobuf:"varint,2,opt,name=type,proto3,enum=meta.FieldEditorInfoType" json:"type,omitempty"`
	Readonly             bool                `protobuf:"varint,3,opt,name=readonly,proto3" json:"readonly,omitempty"`
	ForeignType          string              `protobuf:"bytes,4,opt,name=foreignType,proto3" json:"foreignType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *FieldEditorInfo) Reset()         { *m = FieldEditorInfo{} }
func (m *FieldEditorInfo) String() string { return proto.CompactTextString(m) }
func (*FieldEditorInfo) ProtoMessage()    {}
func (*FieldEditorInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{3}
}
func (m *FieldEditorInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldEditorInfo.Unmarshal(m, b)
}
func (m *FieldEditorInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldEditorInfo.Marshal(b, m, deterministic)
}
func (dst *FieldEditorInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldEditorInfo.Merge(dst, src)
}
func (m *FieldEditorInfo) XXX_Size() int {
	return xxx_messageInfo_FieldEditorInfo.Size(m)
}
func (m *FieldEditorInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldEditorInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FieldEditorInfo proto.InternalMessageInfo

func (m *FieldEditorInfo) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *FieldEditorInfo) GetType() FieldEditorInfoType {
	if m != nil {
		return m.Type
	}
	return FieldEditorInfoType_Default
}

func (m *FieldEditorInfo) GetReadonly() bool {
	if m != nil {
		return m.Readonly
	}
	return false
}

func (m *FieldEditorInfo) GetForeignType() string {
	if m != nil {
		return m.ForeignType
	}
	return ""
}

type KindEditor struct {
	Singular             string   `protobuf:"bytes,1,opt,name=singular,proto3" json:"singular,omitempty"`
	Plural               string   `protobuf:"bytes,2,opt,name=plural,proto3" json:"plural,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KindEditor) Reset()         { *m = KindEditor{} }
func (m *KindEditor) String() string { return proto.CompactTextString(m) }
func (*KindEditor) ProtoMessage()    {}
func (*KindEditor) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{4}
}
func (m *KindEditor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KindEditor.Unmarshal(m, b)
}
func (m *KindEditor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KindEditor.Marshal(b, m, deterministic)
}
func (dst *KindEditor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KindEditor.Merge(dst, src)
}
func (m *KindEditor) XXX_Size() int {
	return xxx_messageInfo_KindEditor.Size(m)
}
func (m *KindEditor) XXX_DiscardUnknown() {
	xxx_messageInfo_KindEditor.DiscardUnknown(m)
}

var xxx_messageInfo_KindEditor proto.InternalMessageInfo

func (m *KindEditor) GetSingular() string {
	if m != nil {
		return m.Singular
	}
	return ""
}

func (m *KindEditor) GetPlural() string {
	if m != nil {
		return m.Plural
	}
	return ""
}

type Kind struct {
	Name                 string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Fields               []*Field    `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	Editor               *KindEditor `protobuf:"bytes,3,opt,name=editor,proto3" json:"editor,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Kind) Reset()         { *m = Kind{} }
func (m *Kind) String() string { return proto.CompactTextString(m) }
func (*Kind) ProtoMessage()    {}
func (*Kind) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{5}
}
func (m *Kind) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Kind.Unmarshal(m, b)
}
func (m *Kind) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Kind.Marshal(b, m, deterministic)
}
func (dst *Kind) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Kind.Merge(dst, src)
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

func (m *Kind) GetEditor() *KindEditor {
	if m != nil {
		return m.Editor
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
	return fileDescriptor_meta_217312ed94b5bec3, []int{6}
}
func (m *Schema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Schema.Unmarshal(m, b)
}
func (m *Schema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Schema.Marshal(b, m, deterministic)
}
func (dst *Schema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Schema.Merge(dst, src)
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
	return fileDescriptor_meta_217312ed94b5bec3, []int{7}
}
func (m *GetSchemaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSchemaRequest.Unmarshal(m, b)
}
func (m *GetSchemaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSchemaRequest.Marshal(b, m, deterministic)
}
func (dst *GetSchemaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSchemaRequest.Merge(dst, src)
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
	return fileDescriptor_meta_217312ed94b5bec3, []int{8}
}
func (m *GetSchemaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSchemaResponse.Unmarshal(m, b)
}
func (m *GetSchemaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSchemaResponse.Marshal(b, m, deterministic)
}
func (dst *GetSchemaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSchemaResponse.Merge(dst, src)
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

type MetaListEntitiesRequest struct {
	Start                []byte   `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	Limit                uint32   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	KindName             string   `protobuf:"bytes,3,opt,name=kindName,proto3" json:"kindName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetaListEntitiesRequest) Reset()         { *m = MetaListEntitiesRequest{} }
func (m *MetaListEntitiesRequest) String() string { return proto.CompactTextString(m) }
func (*MetaListEntitiesRequest) ProtoMessage()    {}
func (*MetaListEntitiesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{9}
}
func (m *MetaListEntitiesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetaListEntitiesRequest.Unmarshal(m, b)
}
func (m *MetaListEntitiesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetaListEntitiesRequest.Marshal(b, m, deterministic)
}
func (dst *MetaListEntitiesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaListEntitiesRequest.Merge(dst, src)
}
func (m *MetaListEntitiesRequest) XXX_Size() int {
	return xxx_messageInfo_MetaListEntitiesRequest.Size(m)
}
func (m *MetaListEntitiesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaListEntitiesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MetaListEntitiesRequest proto.InternalMessageInfo

func (m *MetaListEntitiesRequest) GetStart() []byte {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *MetaListEntitiesRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *MetaListEntitiesRequest) GetKindName() string {
	if m != nil {
		return m.KindName
	}
	return ""
}

type MetaListEntitiesResponse struct {
	Next                 []byte        `protobuf:"bytes,1,opt,name=next,proto3" json:"next,omitempty"`
	MoreResults          bool          `protobuf:"varint,2,opt,name=moreResults,proto3" json:"moreResults,omitempty"`
	Entities             []*MetaEntity `protobuf:"bytes,3,rep,name=entities,proto3" json:"entities,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *MetaListEntitiesResponse) Reset()         { *m = MetaListEntitiesResponse{} }
func (m *MetaListEntitiesResponse) String() string { return proto.CompactTextString(m) }
func (*MetaListEntitiesResponse) ProtoMessage()    {}
func (*MetaListEntitiesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{10}
}
func (m *MetaListEntitiesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetaListEntitiesResponse.Unmarshal(m, b)
}
func (m *MetaListEntitiesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetaListEntitiesResponse.Marshal(b, m, deterministic)
}
func (dst *MetaListEntitiesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaListEntitiesResponse.Merge(dst, src)
}
func (m *MetaListEntitiesResponse) XXX_Size() int {
	return xxx_messageInfo_MetaListEntitiesResponse.Size(m)
}
func (m *MetaListEntitiesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaListEntitiesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MetaListEntitiesResponse proto.InternalMessageInfo

func (m *MetaListEntitiesResponse) GetNext() []byte {
	if m != nil {
		return m.Next
	}
	return nil
}

func (m *MetaListEntitiesResponse) GetMoreResults() bool {
	if m != nil {
		return m.MoreResults
	}
	return false
}

func (m *MetaListEntitiesResponse) GetEntities() []*MetaEntity {
	if m != nil {
		return m.Entities
	}
	return nil
}

type MetaEntity struct {
	Key                  *Key     `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values               []*Value `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetaEntity) Reset()         { *m = MetaEntity{} }
func (m *MetaEntity) String() string { return proto.CompactTextString(m) }
func (*MetaEntity) ProtoMessage()    {}
func (*MetaEntity) Descriptor() ([]byte, []int) {
	return fileDescriptor_meta_217312ed94b5bec3, []int{11}
}
func (m *MetaEntity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetaEntity.Unmarshal(m, b)
}
func (m *MetaEntity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetaEntity.Marshal(b, m, deterministic)
}
func (dst *MetaEntity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaEntity.Merge(dst, src)
}
func (m *MetaEntity) XXX_Size() int {
	return xxx_messageInfo_MetaEntity.Size(m)
}
func (m *MetaEntity) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaEntity.DiscardUnknown(m)
}

var xxx_messageInfo_MetaEntity proto.InternalMessageInfo

func (m *MetaEntity) GetKey() *Key {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *MetaEntity) GetValues() []*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterType((*Key)(nil), "meta.Key")
	proto.RegisterType((*Value)(nil), "meta.Value")
	proto.RegisterType((*Field)(nil), "meta.Field")
	proto.RegisterType((*FieldEditorInfo)(nil), "meta.FieldEditorInfo")
	proto.RegisterType((*KindEditor)(nil), "meta.KindEditor")
	proto.RegisterType((*Kind)(nil), "meta.Kind")
	proto.RegisterType((*Schema)(nil), "meta.Schema")
	proto.RegisterType((*GetSchemaRequest)(nil), "meta.GetSchemaRequest")
	proto.RegisterType((*GetSchemaResponse)(nil), "meta.GetSchemaResponse")
	proto.RegisterType((*MetaListEntitiesRequest)(nil), "meta.MetaListEntitiesRequest")
	proto.RegisterType((*MetaListEntitiesResponse)(nil), "meta.MetaListEntitiesResponse")
	proto.RegisterType((*MetaEntity)(nil), "meta.MetaEntity")
	proto.RegisterEnum("meta.ValueType", ValueType_name, ValueType_value)
	proto.RegisterEnum("meta.FieldEditorInfoType", FieldEditorInfoType_name, FieldEditorInfoType_value)
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
	MetaList(ctx context.Context, in *MetaListEntitiesRequest, opts ...grpc.CallOption) (*MetaListEntitiesResponse, error)
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

func (c *configstoreMetaServiceClient) MetaList(ctx context.Context, in *MetaListEntitiesRequest, opts ...grpc.CallOption) (*MetaListEntitiesResponse, error) {
	out := new(MetaListEntitiesResponse)
	err := c.cc.Invoke(ctx, "/meta.ConfigstoreMetaService/MetaList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigstoreMetaServiceServer is the server API for ConfigstoreMetaService service.
type ConfigstoreMetaServiceServer interface {
	GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error)
	MetaList(context.Context, *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error)
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

func _ConfigstoreMetaService_MetaList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetaListEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigstoreMetaServiceServer).MetaList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.ConfigstoreMetaService/MetaList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigstoreMetaServiceServer).MetaList(ctx, req.(*MetaListEntitiesRequest))
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
		{
			MethodName: "MetaList",
			Handler:    _ConfigstoreMetaService_MetaList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "meta.proto",
}

func init() { proto.RegisterFile("meta.proto", fileDescriptor_meta_217312ed94b5bec3) }

var fileDescriptor_meta_217312ed94b5bec3 = []byte{
	// 784 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0x5f, 0x8f, 0xdb, 0x44,
	0x10, 0xaf, 0xe3, 0xd8, 0x67, 0x4f, 0xd2, 0xab, 0x59, 0xe0, 0x1a, 0x0e, 0x51, 0x59, 0xe6, 0x8f,
	0xac, 0x8a, 0xeb, 0x43, 0x40, 0x48, 0x48, 0x15, 0x42, 0xa5, 0x05, 0x9d, 0x52, 0x2a, 0xb4, 0xa9,
	0x78, 0x3d, 0x6d, 0x2e, 0x93, 0xb0, 0x8a, 0xed, 0x35, 0xde, 0xf5, 0x81, 0x5f, 0xf8, 0x1e, 0xbc,
	0xc1, 0x67, 0xe0, 0x0b, 0xa2, 0xfd, 0x13, 0x9f, 0x2f, 0x17, 0x78, 0xf3, 0xfc, 0x66, 0x76, 0x7e,
	0xb3, 0xbf, 0xd9, 0x19, 0x03, 0x94, 0xa8, 0xd8, 0xb3, 0xba, 0x11, 0x4a, 0x90, 0xb1, 0xfe, 0xce,
	0x2e, 0xc0, 0x5f, 0x60, 0x47, 0x12, 0xf0, 0x6f, 0x58, 0x31, 0xf3, 0x52, 0x2f, 0x8f, 0xa9, 0xfe,
	0x24, 0xef, 0x41, 0xc0, 0xe5, 0x12, 0xd5, 0x6c, 0x94, 0x7a, 0x79, 0x44, 0xad, 0x91, 0xfd, 0x33,
	0x82, 0xe0, 0x67, 0x56, 0xb4, 0x48, 0x4e, 0x61, 0xc4, 0xd7, 0xe6, 0x40, 0x40, 0x47, 0x7c, 0x4d,
	0x3e, 0x86, 0xb1, 0xea, 0x6a, 0x34, 0xe1, 0xa7, 0xf3, 0x47, 0xcf, 0x0c, 0x93, 0x09, 0x7d, 0xdb,
	0xd5, 0x48, 0x8d, 0x93, 0xa4, 0x30, 0x59, 0x8b, 0x76, 0x55, 0xa0, 0x71, 0xcc, 0xfc, 0xd4, 0xcb,
	0x3d, 0x3a, 0x84, 0xc8, 0x13, 0x00, 0x5e, 0xa9, 0xaf, 0xbe, 0xb4, 0x01, 0xe3, 0xd4, 0xcb, 0x7d,
	0x3a, 0x40, 0x74, 0x06, 0xa9, 0x1a, 0x5e, 0x6d, 0x6d, 0x40, 0x60, 0x0a, 0x1e, 0x42, 0xe4, 0x33,
	0x38, 0x55, 0xbc, 0x44, 0xa9, 0x58, 0x59, 0xdb, 0xa0, 0x30, 0xf5, 0xf2, 0x29, 0x3d, 0x40, 0x49,
	0x06, 0xd3, 0x95, 0x10, 0x05, 0xb2, 0xca, 0x46, 0x9d, 0x98, 0x7b, 0xde, 0xc1, 0x74, 0x35, 0xab,
	0x4e, 0xa1, 0xb4, 0x11, 0x91, 0xc9, 0x33, 0x40, 0xc8, 0xa7, 0x10, 0xed, 0xb0, 0xb3, 0xde, 0x38,
	0xf5, 0xf2, 0xc9, 0x3c, 0xb6, 0x17, 0x5f, 0x60, 0x47, 0x7b, 0x57, 0xf6, 0xa7, 0x07, 0xc1, 0xf7,
	0x1c, 0x8b, 0xf5, 0x3d, 0xd5, 0x08, 0x8c, 0x2b, 0x56, 0x5a, 0xd5, 0x62, 0x6a, 0xbe, 0x7b, 0x25,
	0xfd, 0xff, 0x53, 0x72, 0x06, 0x27, 0xd7, 0xa2, 0x2c, 0xb1, 0x52, 0x46, 0xa4, 0x98, 0xee, 0x4d,
	0x72, 0x01, 0x21, 0xae, 0xb9, 0x12, 0x8d, 0x11, 0x67, 0x32, 0x7f, 0xdf, 0x26, 0x30, 0xfc, 0xaf,
	0x8c, 0xe3, 0xb2, 0xda, 0x08, 0xea, 0x82, 0xb2, 0xbf, 0x3c, 0x78, 0x74, 0xe0, 0x33, 0x6d, 0xe2,
	0xb2, 0x2e, 0x58, 0xf7, 0x46, 0x17, 0x67, 0x5f, 0xc5, 0x10, 0x22, 0x17, 0x77, 0xba, 0xfd, 0xc1,
	0x51, 0x8a, 0x41, 0xb5, 0xe7, 0x10, 0x35, 0xc8, 0xd6, 0xa2, 0x2a, 0x3a, 0x73, 0xad, 0x88, 0xf6,
	0xb6, 0x26, 0xdb, 0x88, 0x06, 0xf9, 0xb6, 0xd2, 0x07, 0xdc, 0x6d, 0x86, 0x50, 0xf6, 0x2d, 0xc0,
	0x82, 0x57, 0x2e, 0xb3, 0xce, 0x25, 0x79, 0xb5, 0x6d, 0x0b, 0xd6, 0xb8, 0xca, 0x7a, 0x9b, 0x9c,
	0x41, 0x58, 0x17, 0x6d, 0xc3, 0x0a, 0x27, 0xa8, 0xb3, 0x32, 0x0e, 0x63, 0x9d, 0xa1, 0x97, 0xdb,
	0xbb, 0x23, 0x77, 0xb8, 0xd1, 0x85, 0xcb, 0xd9, 0x28, 0xf5, 0xf3, 0xc9, 0x7c, 0x32, 0xb8, 0x0c,
	0x75, 0x2e, 0x92, 0xf7, 0xa2, 0xfa, 0x46, 0xd4, 0xc4, 0xb5, 0xb9, 0x2f, 0xab, 0xd7, 0xf3, 0x1b,
	0x08, 0x97, 0xd7, 0xbf, 0x60, 0xc9, 0x8e, 0x92, 0xa5, 0x10, 0xec, 0x78, 0xd5, 0x73, 0xc1, 0x6d,
	0x1a, 0x6a, 0x1d, 0x19, 0x81, 0xe4, 0x07, 0x54, 0x36, 0x05, 0xc5, 0x5f, 0x5b, 0x94, 0x2a, 0xfb,
	0x1a, 0xde, 0x19, 0x60, 0xb2, 0x16, 0x95, 0x44, 0xf2, 0x09, 0x84, 0xd2, 0x20, 0x86, 0x60, 0x32,
	0x9f, 0xda, 0x5c, 0x2e, 0xca, 0xf9, 0x32, 0x06, 0x8f, 0x7f, 0x44, 0xc5, 0x5e, 0x73, 0xa9, 0x5e,
	0x55, 0x8a, 0x2b, 0x8e, 0xd2, 0x65, 0xd5, 0x13, 0x2e, 0x15, 0x6b, 0x94, 0x39, 0x3f, 0xa5, 0xd6,
	0xd0, 0x68, 0xc1, 0x4b, 0x6e, 0xe7, 0xfe, 0x21, 0xb5, 0x86, 0x16, 0x5d, 0x97, 0x67, 0x9e, 0x83,
	0x6f, 0x45, 0xdf, 0xdb, 0xd9, 0x1f, 0x30, 0xbb, 0x4f, 0xe1, 0x8a, 0xd4, 0x1a, 0xe0, 0xef, 0x7b,
	0x0a, 0xf3, 0xad, 0x1b, 0x5e, 0x8a, 0x06, 0x29, 0xca, 0xb6, 0x50, 0xd2, 0xed, 0x97, 0x21, 0x44,
	0x3e, 0x87, 0x08, 0x5d, 0xa6, 0x99, 0x6f, 0x84, 0x72, 0x7a, 0x6b, 0x1e, 0xc3, 0xd1, 0xd1, 0x3e,
	0x22, 0x7b, 0x03, 0x70, 0x8b, 0x93, 0x0f, 0xc1, 0xdf, 0x61, 0xe7, 0x34, 0x19, 0x4c, 0xa3, 0x46,
	0x75, 0xaf, 0x6f, 0xf4, 0x20, 0x1d, 0xf4, 0xda, 0x0c, 0x17, 0x75, 0xae, 0xa7, 0x57, 0x10, 0xf7,
	0xd3, 0x46, 0x00, 0xc2, 0x97, 0x66, 0x3d, 0x25, 0x0f, 0x48, 0x0c, 0xc1, 0xa5, 0xde, 0x44, 0x89,
	0xa7, 0xe1, 0xa5, 0xd9, 0x39, 0xc9, 0x88, 0x3c, 0x84, 0xf8, 0xed, 0x7e, 0xb5, 0x24, 0x3e, 0x99,
	0xc0, 0xc9, 0x0b, 0xbb, 0x43, 0x92, 0xb1, 0x3e, 0xf2, 0x42, 0xaf, 0x8b, 0x24, 0x20, 0x11, 0x8c,
	0x17, 0xd8, 0x5d, 0x25, 0xe1, 0xd3, 0xe7, 0xf0, 0xee, 0x91, 0x51, 0xd1, 0x07, 0x5f, 0xe2, 0x86,
	0xb5, 0x85, 0x4a, 0x1e, 0x90, 0x29, 0x44, 0x3f, 0x31, 0x29, 0x7f, 0x13, 0xcd, 0xda, 0xd2, 0xbd,
	0x16, 0x62, 0xd7, 0xd6, 0xc9, 0x68, 0xfe, 0xb7, 0x07, 0x67, 0xdf, 0x89, 0x6a, 0xc3, 0xb7, 0x52,
	0x89, 0x06, 0xf5, 0xd5, 0x97, 0xd8, 0xdc, 0xf0, 0x6b, 0x24, 0xcf, 0x21, 0xee, 0xdf, 0x09, 0x39,
	0xb3, 0x77, 0x3b, 0x7c, 0x4c, 0xe7, 0x8f, 0xef, 0xe1, 0xae, 0x57, 0x97, 0x10, 0xed, 0xfb, 0x48,
	0x3e, 0xba, 0xd5, 0xfb, 0xc8, 0xd3, 0x39, 0x7f, 0xf2, 0x5f, 0x6e, 0x9b, 0x6a, 0x15, 0x9a, 0x5f,
	0xcc, 0x17, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x5b, 0xc8, 0x23, 0x70, 0x06, 0x00, 0x00,
}

