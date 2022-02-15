// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pay.proto

package pay

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

type CreateRequest struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Oid                  int64    `protobuf:"varint,2,opt,name=Oid,proto3" json:"Oid,omitempty"`
	Amount               int64    `protobuf:"varint,3,opt,name=Amount,proto3" json:"Amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{0}
}
func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (dst *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(dst, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *CreateRequest) GetOid() int64 {
	if m != nil {
		return m.Oid
	}
	return 0
}

func (m *CreateRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type CreateResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateResponse) Reset()         { *m = CreateResponse{} }
func (m *CreateResponse) String() string { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()    {}
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{1}
}
func (m *CreateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateResponse.Unmarshal(m, b)
}
func (m *CreateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateResponse.Marshal(b, m, deterministic)
}
func (dst *CreateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateResponse.Merge(dst, src)
}
func (m *CreateResponse) XXX_Size() int {
	return xxx_messageInfo_CreateResponse.Size(m)
}
func (m *CreateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateResponse proto.InternalMessageInfo

func (m *CreateResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type DetailRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DetailRequest) Reset()         { *m = DetailRequest{} }
func (m *DetailRequest) String() string { return proto.CompactTextString(m) }
func (*DetailRequest) ProtoMessage()    {}
func (*DetailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{2}
}
func (m *DetailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetailRequest.Unmarshal(m, b)
}
func (m *DetailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetailRequest.Marshal(b, m, deterministic)
}
func (dst *DetailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetailRequest.Merge(dst, src)
}
func (m *DetailRequest) XXX_Size() int {
	return xxx_messageInfo_DetailRequest.Size(m)
}
func (m *DetailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DetailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DetailRequest proto.InternalMessageInfo

func (m *DetailRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type DetailResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Uid                  int64    `protobuf:"varint,2,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Oid                  int64    `protobuf:"varint,3,opt,name=Oid,proto3" json:"Oid,omitempty"`
	Amount               int64    `protobuf:"varint,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Source               int64    `protobuf:"varint,5,opt,name=Source,proto3" json:"Source,omitempty"`
	Status               int64    `protobuf:"varint,6,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DetailResponse) Reset()         { *m = DetailResponse{} }
func (m *DetailResponse) String() string { return proto.CompactTextString(m) }
func (*DetailResponse) ProtoMessage()    {}
func (*DetailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{3}
}
func (m *DetailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetailResponse.Unmarshal(m, b)
}
func (m *DetailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetailResponse.Marshal(b, m, deterministic)
}
func (dst *DetailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetailResponse.Merge(dst, src)
}
func (m *DetailResponse) XXX_Size() int {
	return xxx_messageInfo_DetailResponse.Size(m)
}
func (m *DetailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DetailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DetailResponse proto.InternalMessageInfo

func (m *DetailResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *DetailResponse) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *DetailResponse) GetOid() int64 {
	if m != nil {
		return m.Oid
	}
	return 0
}

func (m *DetailResponse) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *DetailResponse) GetSource() int64 {
	if m != nil {
		return m.Source
	}
	return 0
}

func (m *DetailResponse) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type CallbackRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uid                  int64    `protobuf:"varint,2,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Oid                  int64    `protobuf:"varint,3,opt,name=Oid,proto3" json:"Oid,omitempty"`
	Amount               int64    `protobuf:"varint,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Source               int64    `protobuf:"varint,5,opt,name=Source,proto3" json:"Source,omitempty"`
	Status               int64    `protobuf:"varint,6,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallbackRequest) Reset()         { *m = CallbackRequest{} }
func (m *CallbackRequest) String() string { return proto.CompactTextString(m) }
func (*CallbackRequest) ProtoMessage()    {}
func (*CallbackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{4}
}
func (m *CallbackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallbackRequest.Unmarshal(m, b)
}
func (m *CallbackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallbackRequest.Marshal(b, m, deterministic)
}
func (dst *CallbackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallbackRequest.Merge(dst, src)
}
func (m *CallbackRequest) XXX_Size() int {
	return xxx_messageInfo_CallbackRequest.Size(m)
}
func (m *CallbackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CallbackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CallbackRequest proto.InternalMessageInfo

func (m *CallbackRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CallbackRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *CallbackRequest) GetOid() int64 {
	if m != nil {
		return m.Oid
	}
	return 0
}

func (m *CallbackRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *CallbackRequest) GetSource() int64 {
	if m != nil {
		return m.Source
	}
	return 0
}

func (m *CallbackRequest) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type CallbackResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallbackResponse) Reset()         { *m = CallbackResponse{} }
func (m *CallbackResponse) String() string { return proto.CompactTextString(m) }
func (*CallbackResponse) ProtoMessage()    {}
func (*CallbackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_pay_d6253a25ab586d79, []int{5}
}
func (m *CallbackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallbackResponse.Unmarshal(m, b)
}
func (m *CallbackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallbackResponse.Marshal(b, m, deterministic)
}
func (dst *CallbackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallbackResponse.Merge(dst, src)
}
func (m *CallbackResponse) XXX_Size() int {
	return xxx_messageInfo_CallbackResponse.Size(m)
}
func (m *CallbackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CallbackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CallbackResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CreateRequest)(nil), "PayClient.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "PayClient.CreateResponse")
	proto.RegisterType((*DetailRequest)(nil), "PayClient.DetailRequest")
	proto.RegisterType((*DetailResponse)(nil), "PayClient.DetailResponse")
	proto.RegisterType((*CallbackRequest)(nil), "PayClient.CallbackRequest")
	proto.RegisterType((*CallbackResponse)(nil), "PayClient.CallbackResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PayClient is the client API for Pay service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PayClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error)
	Callback(ctx context.Context, in *CallbackRequest, opts ...grpc.CallOption) (*CallbackResponse, error)
}

type payClient struct {
	cc *grpc.ClientConn
}

func NewPayClient(cc *grpc.ClientConn) PayClient {
	return &payClient{cc}
}

func (c *payClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/PayClient.Pay/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *payClient) Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error) {
	out := new(DetailResponse)
	err := c.cc.Invoke(ctx, "/PayClient.Pay/Detail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *payClient) Callback(ctx context.Context, in *CallbackRequest, opts ...grpc.CallOption) (*CallbackResponse, error) {
	out := new(CallbackResponse)
	err := c.cc.Invoke(ctx, "/PayClient.Pay/Callback", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PayServer is the server API for Pay service.
type PayServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Detail(context.Context, *DetailRequest) (*DetailResponse, error)
	Callback(context.Context, *CallbackRequest) (*CallbackResponse, error)
}

func RegisterPayServer(s *grpc.Server, srv PayServer) {
	s.RegisterService(&_Pay_serviceDesc, srv)
}

func _Pay_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PayClient.Pay/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pay_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PayClient.Pay/Detail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayServer).Detail(ctx, req.(*DetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pay_Callback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayServer).Callback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PayClient.Pay/Callback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayServer).Callback(ctx, req.(*CallbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pay_serviceDesc = grpc.ServiceDesc{
	ServiceName: "PayClient.Pay",
	HandlerType: (*PayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Pay_Create_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _Pay_Detail_Handler,
		},
		{
			MethodName: "Callback",
			Handler:    _Pay_Callback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pay.proto",
}

func init() { proto.RegisterFile("pay.proto", fileDescriptor_pay_d6253a25ab586d79) }

var fileDescriptor_pay_d6253a25ab586d79 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x92, 0xb1, 0x4e, 0xc3, 0x30,
	0x18, 0x84, 0x65, 0x9b, 0x5a, 0xf4, 0x97, 0x1a, 0x2a, 0x0f, 0xc8, 0x84, 0x81, 0x2a, 0x13, 0x53,
	0x06, 0x98, 0x3b, 0x40, 0x58, 0x2a, 0x86, 0x56, 0x45, 0x2c, 0x6c, 0x6e, 0xe3, 0xc1, 0x22, 0xc4,
	0x21, 0x71, 0x86, 0xbc, 0x03, 0xbc, 0x17, 0x8f, 0x85, 0x12, 0xc7, 0x25, 0x2e, 0x99, 0xbb, 0xe5,
	0x3f, 0xdb, 0xf7, 0x7f, 0xba, 0x0b, 0x4c, 0x0b, 0xd1, 0xc4, 0x45, 0xa9, 0x8d, 0x66, 0xd3, 0x8d,
	0x68, 0x92, 0x4c, 0xc9, 0xdc, 0x44, 0xcf, 0x30, 0x4b, 0x4a, 0x29, 0x8c, 0xdc, 0xca, 0xcf, 0x5a,
	0x56, 0x86, 0xcd, 0x81, 0xbc, 0xaa, 0x94, 0xa3, 0x05, 0xba, 0x25, 0xdb, 0xf6, 0xb3, 0x55, 0xd6,
	0x2a, 0xe5, 0xd8, 0x2a, 0x6b, 0x95, 0xb2, 0x4b, 0xa0, 0x0f, 0x1f, 0xba, 0xce, 0x0d, 0x27, 0x9d,
	0xd8, 0x4f, 0xd1, 0x02, 0x02, 0x67, 0x56, 0x15, 0x3a, 0xaf, 0x24, 0x0b, 0x00, 0xaf, 0x9c, 0x19,
	0x5e, 0xa5, 0xd1, 0x0d, 0xcc, 0x9e, 0xa4, 0x11, 0x2a, 0x73, 0xeb, 0x8e, 0x2f, 0x7c, 0x21, 0x08,
	0xdc, 0x8d, 0x71, 0x0f, 0x47, 0x88, 0xff, 0x11, 0x92, 0x31, 0xc2, 0xb3, 0x21, 0x61, 0xab, 0xbf,
	0xe8, 0xba, 0xdc, 0x4b, 0x3e, 0xb1, 0xba, 0x9d, 0x3a, 0xdd, 0x08, 0x53, 0x57, 0x9c, 0xf6, 0x7a,
	0x37, 0x45, 0xdf, 0x08, 0x2e, 0x12, 0x91, 0x65, 0x3b, 0xb1, 0x7f, 0x1f, 0x20, 0x1f, 0x02, 0xc2,
	0xea, 0xb4, 0x3c, 0x0c, 0xe6, 0x7f, 0x38, 0x36, 0x9f, 0xbb, 0x1f, 0x04, 0x64, 0x23, 0x1a, 0xb6,
	0x04, 0x6a, 0xd3, 0x67, 0x3c, 0x3e, 0x14, 0x1c, 0x7b, 0xed, 0x86, 0x57, 0x23, 0x27, 0x7d, 0xcc,
	0x4b, 0xa0, 0x36, 0x78, 0xef, 0xb9, 0xd7, 0x96, 0xf7, 0xfc, 0xa8, 0xa5, 0x04, 0xce, 0x1d, 0x19,
	0x0b, 0x87, 0x5b, 0xfc, 0xf4, 0xc2, 0xeb, 0xd1, 0x33, 0x6b, 0xf2, 0x38, 0x79, 0x23, 0x85, 0x68,
	0x76, 0xb4, 0xfb, 0x4d, 0xef, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x18, 0x22, 0x28, 0x52, 0xb3,
	0x02, 0x00, 0x00,
}
