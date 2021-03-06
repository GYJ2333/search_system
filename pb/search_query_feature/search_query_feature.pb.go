// Code generated by protoc-gen-go. DO NOT EDIT.
// source: search_query_feature.proto

package search_query_feature

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type SetType int32

const (
	SetType_TYPE_UNKNOWN SetType = 0
	SetType_TYPE_UPDATE  SetType = 1
	SetType_TYPE_ADD     SetType = 2
	SetType_TYPE_DELETE  SetType = 3
)

var SetType_name = map[int32]string{
	0: "TYPE_UNKNOWN",
	1: "TYPE_UPDATE",
	2: "TYPE_ADD",
	3: "TYPE_DELETE",
}

var SetType_value = map[string]int32{
	"TYPE_UNKNOWN": 0,
	"TYPE_UPDATE":  1,
	"TYPE_ADD":     2,
	"TYPE_DELETE":  3,
}

func (x SetType) String() string {
	return proto.EnumName(SetType_name, int32(x))
}

func (SetType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{0}
}

type ResponseHeader struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Err                  string   `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseHeader) Reset()         { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()    {}
func (*ResponseHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{0}
}

func (m *ResponseHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseHeader.Unmarshal(m, b)
}
func (m *ResponseHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseHeader.Marshal(b, m, deterministic)
}
func (m *ResponseHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseHeader.Merge(m, src)
}
func (m *ResponseHeader) XXX_Size() int {
	return xxx_messageInfo_ResponseHeader.Size(m)
}
func (m *ResponseHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseHeader proto.InternalMessageInfo

func (m *ResponseHeader) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ResponseHeader) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type Query struct {
	QueryId              string            `protobuf:"bytes,1,opt,name=query_id,json=queryId,proto3" json:"query_id,omitempty"`
	QueryName            string            `protobuf:"bytes,2,opt,name=query_name,json=queryName,proto3" json:"query_name,omitempty"`
	Kind                 string            `protobuf:"bytes,3,opt,name=kind,proto3" json:"kind,omitempty"`
	Feature              string            `protobuf:"bytes,4,opt,name=feature,proto3" json:"feature,omitempty"`
	Ext                  map[string]string `protobuf:"bytes,16,rep,name=ext,proto3" json:"ext,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{1}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetQueryId() string {
	if m != nil {
		return m.QueryId
	}
	return ""
}

func (m *Query) GetQueryName() string {
	if m != nil {
		return m.QueryName
	}
	return ""
}

func (m *Query) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Query) GetFeature() string {
	if m != nil {
		return m.Feature
	}
	return ""
}

func (m *Query) GetExt() map[string]string {
	if m != nil {
		return m.Ext
	}
	return nil
}

type OfflineRequest struct {
	Type                 SetType           `protobuf:"varint,1,opt,name=type,proto3,enum=search_query_feature.SetType" json:"type,omitempty"`
	Querys               []*Query          `protobuf:"bytes,2,rep,name=querys,proto3" json:"querys,omitempty"`
	Ext                  map[string]string `protobuf:"bytes,16,rep,name=ext,proto3" json:"ext,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *OfflineRequest) Reset()         { *m = OfflineRequest{} }
func (m *OfflineRequest) String() string { return proto.CompactTextString(m) }
func (*OfflineRequest) ProtoMessage()    {}
func (*OfflineRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{2}
}

func (m *OfflineRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OfflineRequest.Unmarshal(m, b)
}
func (m *OfflineRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OfflineRequest.Marshal(b, m, deterministic)
}
func (m *OfflineRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OfflineRequest.Merge(m, src)
}
func (m *OfflineRequest) XXX_Size() int {
	return xxx_messageInfo_OfflineRequest.Size(m)
}
func (m *OfflineRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OfflineRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OfflineRequest proto.InternalMessageInfo

func (m *OfflineRequest) GetType() SetType {
	if m != nil {
		return m.Type
	}
	return SetType_TYPE_UNKNOWN
}

func (m *OfflineRequest) GetQuerys() []*Query {
	if m != nil {
		return m.Querys
	}
	return nil
}

func (m *OfflineRequest) GetExt() map[string]string {
	if m != nil {
		return m.Ext
	}
	return nil
}

type Status struct {
	QueryId              string   `protobuf:"bytes,1,opt,name=query_id,json=queryId,proto3" json:"query_id,omitempty"`
	Ok                   bool     `protobuf:"varint,2,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg                  string   `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{3}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetQueryId() string {
	if m != nil {
		return m.QueryId
	}
	return ""
}

func (m *Status) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *Status) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type OfflineResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	QueryStatus          []*Status       `protobuf:"bytes,2,rep,name=query_status,json=queryStatus,proto3" json:"query_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *OfflineResponse) Reset()         { *m = OfflineResponse{} }
func (m *OfflineResponse) String() string { return proto.CompactTextString(m) }
func (*OfflineResponse) ProtoMessage()    {}
func (*OfflineResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{4}
}

func (m *OfflineResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OfflineResponse.Unmarshal(m, b)
}
func (m *OfflineResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OfflineResponse.Marshal(b, m, deterministic)
}
func (m *OfflineResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OfflineResponse.Merge(m, src)
}
func (m *OfflineResponse) XXX_Size() int {
	return xxx_messageInfo_OfflineResponse.Size(m)
}
func (m *OfflineResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OfflineResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OfflineResponse proto.InternalMessageInfo

func (m *OfflineResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *OfflineResponse) GetQueryStatus() []*Status {
	if m != nil {
		return m.QueryStatus
	}
	return nil
}

type OnlineRequest struct {
	UserId               string            `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Features             []string          `protobuf:"bytes,1,rep,name=features,proto3" json:"features,omitempty"`
	Ext                  map[string]string `protobuf:"bytes,16,rep,name=ext,proto3" json:"ext,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *OnlineRequest) Reset()         { *m = OnlineRequest{} }
func (m *OnlineRequest) String() string { return proto.CompactTextString(m) }
func (*OnlineRequest) ProtoMessage()    {}
func (*OnlineRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{5}
}

func (m *OnlineRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OnlineRequest.Unmarshal(m, b)
}
func (m *OnlineRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OnlineRequest.Marshal(b, m, deterministic)
}
func (m *OnlineRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OnlineRequest.Merge(m, src)
}
func (m *OnlineRequest) XXX_Size() int {
	return xxx_messageInfo_OnlineRequest.Size(m)
}
func (m *OnlineRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OnlineRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OnlineRequest proto.InternalMessageInfo

func (m *OnlineRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *OnlineRequest) GetFeatures() []string {
	if m != nil {
		return m.Features
	}
	return nil
}

func (m *OnlineRequest) GetExt() map[string]string {
	if m != nil {
		return m.Ext
	}
	return nil
}

type OnlineResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	QueryIds             []string        `protobuf:"bytes,2,rep,name=query_ids,json=queryIds,proto3" json:"query_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *OnlineResponse) Reset()         { *m = OnlineResponse{} }
func (m *OnlineResponse) String() string { return proto.CompactTextString(m) }
func (*OnlineResponse) ProtoMessage()    {}
func (*OnlineResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{6}
}

func (m *OnlineResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OnlineResponse.Unmarshal(m, b)
}
func (m *OnlineResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OnlineResponse.Marshal(b, m, deterministic)
}
func (m *OnlineResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OnlineResponse.Merge(m, src)
}
func (m *OnlineResponse) XXX_Size() int {
	return xxx_messageInfo_OnlineResponse.Size(m)
}
func (m *OnlineResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OnlineResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OnlineResponse proto.InternalMessageInfo

func (m *OnlineResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *OnlineResponse) GetQueryIds() []string {
	if m != nil {
		return m.QueryIds
	}
	return nil
}

type ChoseRequest struct {
	UserId               string            `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	QueryId              string            `protobuf:"bytes,2,opt,name=query_id,json=queryId,proto3" json:"query_id,omitempty"`
	Ext                  map[string]string `protobuf:"bytes,16,rep,name=ext,proto3" json:"ext,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ChoseRequest) Reset()         { *m = ChoseRequest{} }
func (m *ChoseRequest) String() string { return proto.CompactTextString(m) }
func (*ChoseRequest) ProtoMessage()    {}
func (*ChoseRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{7}
}

func (m *ChoseRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChoseRequest.Unmarshal(m, b)
}
func (m *ChoseRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChoseRequest.Marshal(b, m, deterministic)
}
func (m *ChoseRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChoseRequest.Merge(m, src)
}
func (m *ChoseRequest) XXX_Size() int {
	return xxx_messageInfo_ChoseRequest.Size(m)
}
func (m *ChoseRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChoseRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChoseRequest proto.InternalMessageInfo

func (m *ChoseRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ChoseRequest) GetQueryId() string {
	if m != nil {
		return m.QueryId
	}
	return ""
}

func (m *ChoseRequest) GetExt() map[string]string {
	if m != nil {
		return m.Ext
	}
	return nil
}

type ChoseResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ChoseResponse) Reset()         { *m = ChoseResponse{} }
func (m *ChoseResponse) String() string { return proto.CompactTextString(m) }
func (*ChoseResponse) ProtoMessage()    {}
func (*ChoseResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9fe36960e0ca820, []int{8}
}

func (m *ChoseResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChoseResponse.Unmarshal(m, b)
}
func (m *ChoseResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChoseResponse.Marshal(b, m, deterministic)
}
func (m *ChoseResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChoseResponse.Merge(m, src)
}
func (m *ChoseResponse) XXX_Size() int {
	return xxx_messageInfo_ChoseResponse.Size(m)
}
func (m *ChoseResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChoseResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChoseResponse proto.InternalMessageInfo

func (m *ChoseResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterEnum("search_query_feature.SetType", SetType_name, SetType_value)
	proto.RegisterType((*ResponseHeader)(nil), "search_query_feature.ResponseHeader")
	proto.RegisterType((*Query)(nil), "search_query_feature.Query")
	proto.RegisterMapType((map[string]string)(nil), "search_query_feature.Query.ExtEntry")
	proto.RegisterType((*OfflineRequest)(nil), "search_query_feature.OfflineRequest")
	proto.RegisterMapType((map[string]string)(nil), "search_query_feature.OfflineRequest.ExtEntry")
	proto.RegisterType((*Status)(nil), "search_query_feature.Status")
	proto.RegisterType((*OfflineResponse)(nil), "search_query_feature.OfflineResponse")
	proto.RegisterType((*OnlineRequest)(nil), "search_query_feature.OnlineRequest")
	proto.RegisterMapType((map[string]string)(nil), "search_query_feature.OnlineRequest.ExtEntry")
	proto.RegisterType((*OnlineResponse)(nil), "search_query_feature.OnlineResponse")
	proto.RegisterType((*ChoseRequest)(nil), "search_query_feature.ChoseRequest")
	proto.RegisterMapType((map[string]string)(nil), "search_query_feature.ChoseRequest.ExtEntry")
	proto.RegisterType((*ChoseResponse)(nil), "search_query_feature.ChoseResponse")
}

func init() { proto.RegisterFile("search_query_feature.proto", fileDescriptor_a9fe36960e0ca820) }

var fileDescriptor_a9fe36960e0ca820 = []byte{
	// 610 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0x49, 0x7f, 0x92, 0xb3, 0xae, 0x8b, 0xac, 0x49, 0x84, 0x8c, 0x49, 0x55, 0x36, 0xa4,
	0x8a, 0x9f, 0x49, 0x6c, 0xd2, 0x84, 0x10, 0x30, 0x4d, 0x34, 0xc0, 0x18, 0x74, 0xc5, 0x2d, 0x42,
	0x5c, 0x55, 0x61, 0x71, 0xd7, 0xaa, 0x6b, 0xd2, 0xc5, 0x0e, 0x5a, 0x2f, 0x79, 0x03, 0xde, 0x08,
	0x89, 0x77, 0xe1, 0x19, 0xb8, 0x45, 0xb1, 0xdd, 0x2a, 0x95, 0xd2, 0x6e, 0x17, 0xe5, 0xce, 0xe7,
	0xc7, 0x9f, 0xbf, 0xf3, 0x7d, 0xa7, 0x29, 0x38, 0x8c, 0xfa, 0xf1, 0x79, 0xbf, 0x7b, 0x95, 0xd0,
	0x78, 0xd2, 0xed, 0x51, 0x9f, 0x27, 0x31, 0xdd, 0x1b, 0xc7, 0x11, 0x8f, 0xf0, 0x66, 0x5e, 0xcd,
	0x3d, 0x84, 0x2a, 0xa1, 0x6c, 0x1c, 0x85, 0x8c, 0xbe, 0xa3, 0x7e, 0x40, 0x63, 0x8c, 0xa1, 0x70,
	0x1e, 0x05, 0xd4, 0x46, 0x35, 0x54, 0x5f, 0x27, 0xe2, 0x8c, 0x2d, 0xd0, 0x69, 0x1c, 0xdb, 0x7a,
	0x0d, 0xd5, 0x4d, 0x92, 0x1e, 0xdd, 0x3f, 0x08, 0x8a, 0x9f, 0x52, 0x24, 0x7c, 0x0f, 0x0c, 0x09,
	0x39, 0x08, 0xc4, 0x1d, 0x93, 0x94, 0x45, 0x7c, 0x12, 0xe0, 0x6d, 0x00, 0x59, 0x0a, 0xfd, 0x11,
	0xb5, 0x35, 0x51, 0x34, 0x45, 0xa6, 0xe9, 0x8f, 0x68, 0xfa, 0xd2, 0x70, 0x10, 0x06, 0x0a, 0x56,
	0x9c, 0xb1, 0x0d, 0x65, 0x45, 0xcd, 0x2e, 0x48, 0x30, 0x15, 0xe2, 0x43, 0xd0, 0xe9, 0x35, 0xb7,
	0xad, 0x9a, 0x5e, 0x5f, 0xdb, 0xdf, 0xdd, 0xcb, 0x9d, 0x54, 0x30, 0xda, 0xf3, 0xae, 0xb9, 0x17,
	0xf2, 0x78, 0x42, 0xd2, 0x0b, 0xce, 0x21, 0x18, 0xd3, 0x44, 0x3a, 0xc7, 0x90, 0x4e, 0x14, 0xcd,
	0xf4, 0x88, 0x37, 0xa1, 0xf8, 0xdd, 0xbf, 0x4c, 0xa6, 0xec, 0x64, 0xf0, 0x5c, 0x7b, 0x86, 0xde,
	0x17, 0x8c, 0xa2, 0x65, 0xb9, 0x7f, 0x11, 0x54, 0xcf, 0x7a, 0xbd, 0xcb, 0x41, 0x48, 0x09, 0xbd,
	0x4a, 0x28, 0xe3, 0xf8, 0x29, 0x14, 0xf8, 0x64, 0x2c, 0x05, 0xaa, 0xee, 0x6f, 0xe7, 0x33, 0x69,
	0x53, 0xde, 0x99, 0x8c, 0x29, 0x11, 0xad, 0xf8, 0x00, 0x4a, 0xa2, 0xcc, 0x6c, 0x4d, 0xd0, 0xdf,
	0x5a, 0x42, 0x9f, 0xa8, 0x56, 0x7c, 0x94, 0x1d, 0xf8, 0x49, 0xfe, 0x8d, 0x79, 0x6a, 0x2b, 0x9b,
	0x5c, 0xb7, 0x2c, 0xd7, 0x83, 0x52, 0x9b, 0xfb, 0x3c, 0x61, 0xcb, 0x1c, 0xae, 0x82, 0x16, 0x0d,
	0x05, 0x82, 0x41, 0xb4, 0x68, 0x98, 0x3e, 0x33, 0x62, 0x17, 0xd3, 0x45, 0x19, 0xb1, 0x0b, 0xf7,
	0x27, 0x82, 0x8d, 0x19, 0x4b, 0xb9, 0x68, 0xf8, 0x05, 0x94, 0xfa, 0x62, 0xd9, 0x04, 0xdc, 0x42,
	0x37, 0xe7, 0x17, 0x93, 0xa8, 0x3b, 0xf8, 0x08, 0x2a, 0xb2, 0x8f, 0x09, 0x7a, 0x4a, 0xd2, 0xfb,
	0x0b, 0x7c, 0x10, 0x3d, 0x64, 0x4d, 0x64, 0x65, 0xe0, 0xfe, 0x46, 0xb0, 0x7e, 0x16, 0x66, 0x2d,
	0xbd, 0x0b, 0xe5, 0x84, 0xd1, 0x38, 0x1d, 0x50, 0xaa, 0x51, 0x4a, 0xc3, 0x93, 0x00, 0x3b, 0x60,
	0x28, 0x24, 0x66, 0xa3, 0x9a, 0x5e, 0x37, 0xc9, 0x2c, 0xc6, 0xaf, 0xb2, 0xfe, 0x3c, 0x5e, 0xe0,
	0x4f, 0xf8, 0xbf, 0xec, 0x19, 0x42, 0x75, 0x0a, 0xbe, 0x12, 0x55, 0xb7, 0xc0, 0x9c, 0x9a, 0x2c,
	0x25, 0x35, 0x89, 0xa1, 0x5c, 0x66, 0xee, 0x2f, 0x04, 0x95, 0xd7, 0xfd, 0x88, 0xe5, 0x09, 0x86,
	0xe6, 0x04, 0xcb, 0xee, 0x8a, 0x36, 0xbf, 0x2b, 0x2f, 0xb3, 0x7a, 0x3d, 0xca, 0x27, 0x97, 0x7d,
	0x64, 0xa5, 0x72, 0x7d, 0x84, 0x75, 0x85, 0xbd, 0x0a, 0xb5, 0x1e, 0x9e, 0x42, 0x59, 0xfd, 0xc2,
	0xb1, 0x05, 0x95, 0xce, 0xd7, 0x96, 0xd7, 0xfd, 0xdc, 0x3c, 0x6d, 0x9e, 0x7d, 0x69, 0x5a, 0x77,
	0xf0, 0x06, 0xac, 0xc9, 0x4c, 0xab, 0x71, 0xdc, 0xf1, 0x2c, 0x84, 0x2b, 0x60, 0x88, 0xc4, 0x71,
	0xa3, 0x61, 0x69, 0xb3, 0x72, 0xc3, 0xfb, 0xe0, 0x75, 0x3c, 0x4b, 0xdf, 0xff, 0xa1, 0x41, 0xf9,
	0x8d, 0xfa, 0xca, 0xb5, 0x40, 0x7f, 0x4b, 0x39, 0xde, 0xb9, 0xc5, 0x3a, 0x39, 0xbb, 0xcb, 0x9b,
	0xd4, 0xa0, 0x04, 0xf4, 0x36, 0xe5, 0x78, 0xf7, 0x36, 0x1f, 0x10, 0xe7, 0xc1, 0x0d, 0x5d, 0x0a,
	0xb3, 0x05, 0x45, 0xa1, 0x26, 0x76, 0x6f, 0xb6, 0xd1, 0xd9, 0x59, 0xda, 0x23, 0x11, 0xbf, 0x95,
	0xc4, 0x9f, 0xd4, 0xc1, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x19, 0x82, 0x59, 0x29, 0xc2, 0x06,
	0x00, 0x00,
}
