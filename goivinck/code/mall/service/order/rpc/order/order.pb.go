// Code generated by protoc-gen-go. DO NOT EDIT.
// source: order.proto

package order

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
	Pid                  int64    `protobuf:"varint,2,opt,name=Pid,proto3" json:"Pid,omitempty"`
	Amount               int64    `protobuf:"varint,3,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Status               int64    `protobuf:"varint,4,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{0}
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

func (m *CreateRequest) GetPid() int64 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *CreateRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *CreateRequest) GetStatus() int64 {
	if m != nil {
		return m.Status
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
	return fileDescriptor_order_bd099b846d1810ce, []int{1}
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

type UpdateRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Uid                  int64    `protobuf:"varint,2,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Pid                  int64    `protobuf:"varint,3,opt,name=Pid,proto3" json:"Pid,omitempty"`
	Amount               int64    `protobuf:"varint,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Status               int64    `protobuf:"varint,5,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{2}
}
func (m *UpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRequest.Unmarshal(m, b)
}
func (m *UpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRequest.Merge(dst, src)
}
func (m *UpdateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRequest.Size(m)
}
func (m *UpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRequest proto.InternalMessageInfo

func (m *UpdateRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *UpdateRequest) GetPid() int64 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *UpdateRequest) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *UpdateRequest) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type UpdateResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateResponse) Reset()         { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()    {}
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{3}
}
func (m *UpdateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateResponse.Unmarshal(m, b)
}
func (m *UpdateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateResponse.Marshal(b, m, deterministic)
}
func (dst *UpdateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateResponse.Merge(dst, src)
}
func (m *UpdateResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateResponse.Size(m)
}
func (m *UpdateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateResponse proto.InternalMessageInfo

type RemoveRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveRequest) Reset()         { *m = RemoveRequest{} }
func (m *RemoveRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveRequest) ProtoMessage()    {}
func (*RemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{4}
}
func (m *RemoveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveRequest.Unmarshal(m, b)
}
func (m *RemoveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveRequest.Marshal(b, m, deterministic)
}
func (dst *RemoveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveRequest.Merge(dst, src)
}
func (m *RemoveRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveRequest.Size(m)
}
func (m *RemoveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveRequest proto.InternalMessageInfo

func (m *RemoveRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type RemoveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveResponse) Reset()         { *m = RemoveResponse{} }
func (m *RemoveResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveResponse) ProtoMessage()    {}
func (*RemoveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{5}
}
func (m *RemoveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveResponse.Unmarshal(m, b)
}
func (m *RemoveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveResponse.Marshal(b, m, deterministic)
}
func (dst *RemoveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveResponse.Merge(dst, src)
}
func (m *RemoveResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveResponse.Size(m)
}
func (m *RemoveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveResponse proto.InternalMessageInfo

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
	return fileDescriptor_order_bd099b846d1810ce, []int{6}
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
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uid                  int64    `protobuf:"varint,2,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Pid                  int64    `protobuf:"varint,3,opt,name=Pid,proto3" json:"Pid,omitempty"`
	Amount               int64    `protobuf:"varint,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Status               int64    `protobuf:"varint,5,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DetailResponse) Reset()         { *m = DetailResponse{} }
func (m *DetailResponse) String() string { return proto.CompactTextString(m) }
func (*DetailResponse) ProtoMessage()    {}
func (*DetailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{7}
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

func (m *DetailResponse) GetPid() int64 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *DetailResponse) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *DetailResponse) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type ListRequest struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{8}
}
func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (dst *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(dst, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

func (m *ListRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type ListResponse struct {
	Data                 []*DetailResponse `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}
func (*ListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{9}
}
func (m *ListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListResponse.Unmarshal(m, b)
}
func (m *ListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListResponse.Marshal(b, m, deterministic)
}
func (dst *ListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListResponse.Merge(dst, src)
}
func (m *ListResponse) XXX_Size() int {
	return xxx_messageInfo_ListResponse.Size(m)
}
func (m *ListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListResponse proto.InternalMessageInfo

func (m *ListResponse) GetData() []*DetailResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

// 订单支付
type PaidRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaidRequest) Reset()         { *m = PaidRequest{} }
func (m *PaidRequest) String() string { return proto.CompactTextString(m) }
func (*PaidRequest) ProtoMessage()    {}
func (*PaidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{10}
}
func (m *PaidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaidRequest.Unmarshal(m, b)
}
func (m *PaidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaidRequest.Marshal(b, m, deterministic)
}
func (dst *PaidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaidRequest.Merge(dst, src)
}
func (m *PaidRequest) XXX_Size() int {
	return xxx_messageInfo_PaidRequest.Size(m)
}
func (m *PaidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PaidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PaidRequest proto.InternalMessageInfo

func (m *PaidRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type PaidResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaidResponse) Reset()         { *m = PaidResponse{} }
func (m *PaidResponse) String() string { return proto.CompactTextString(m) }
func (*PaidResponse) ProtoMessage()    {}
func (*PaidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_bd099b846d1810ce, []int{11}
}
func (m *PaidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaidResponse.Unmarshal(m, b)
}
func (m *PaidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaidResponse.Marshal(b, m, deterministic)
}
func (dst *PaidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaidResponse.Merge(dst, src)
}
func (m *PaidResponse) XXX_Size() int {
	return xxx_messageInfo_PaidResponse.Size(m)
}
func (m *PaidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PaidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PaidResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CreateRequest)(nil), "OrderClient.CreateRequest")
	proto.RegisterType((*CreateResponse)(nil), "OrderClient.CreateResponse")
	proto.RegisterType((*UpdateRequest)(nil), "OrderClient.UpdateRequest")
	proto.RegisterType((*UpdateResponse)(nil), "OrderClient.UpdateResponse")
	proto.RegisterType((*RemoveRequest)(nil), "OrderClient.RemoveRequest")
	proto.RegisterType((*RemoveResponse)(nil), "OrderClient.RemoveResponse")
	proto.RegisterType((*DetailRequest)(nil), "OrderClient.DetailRequest")
	proto.RegisterType((*DetailResponse)(nil), "OrderClient.DetailResponse")
	proto.RegisterType((*ListRequest)(nil), "OrderClient.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "OrderClient.ListResponse")
	proto.RegisterType((*PaidRequest)(nil), "OrderClient.PaidRequest")
	proto.RegisterType((*PaidResponse)(nil), "OrderClient.PaidResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error)
	Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Paid(ctx context.Context, in *PaidRequest, opts ...grpc.CallOption) (*PaidResponse, error)
}

type orderClient struct {
	cc *grpc.ClientConn
}

func NewOrderClient(cc *grpc.ClientConn) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error) {
	out := new(RemoveResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error) {
	out := new(DetailResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/Detail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) Paid(ctx context.Context, in *PaidRequest, opts ...grpc.CallOption) (*PaidResponse, error) {
	out := new(PaidResponse)
	err := c.cc.Invoke(ctx, "/OrderClient.Order/Paid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
type OrderServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Remove(context.Context, *RemoveRequest) (*RemoveResponse, error)
	Detail(context.Context, *DetailRequest) (*DetailResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Paid(context.Context, *PaidRequest) (*PaidResponse, error)
}

func RegisterOrderServer(s *grpc.Server, srv OrderServer) {
	s.RegisterService(&_Order_serviceDesc, srv)
}

func _Order_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/Detail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).Detail(ctx, req.(*DetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_Paid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).Paid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderClient.Order/Paid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).Paid(ctx, req.(*PaidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Order_serviceDesc = grpc.ServiceDesc{
	ServiceName: "OrderClient.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Order_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Order_Update_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Order_Remove_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _Order_Detail_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Order_List_Handler,
		},
		{
			MethodName: "Paid",
			Handler:    _Order_Paid_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}

func init() { proto.RegisterFile("order.proto", fileDescriptor_order_bd099b846d1810ce) }

var fileDescriptor_order_bd099b846d1810ce = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x41, 0x4f, 0xc2, 0x30,
	0x14, 0x0e, 0xdb, 0x98, 0xc9, 0x1b, 0x2c, 0xa4, 0x07, 0x53, 0x47, 0x0c, 0xa4, 0x27, 0x4e, 0x98,
	0xe0, 0xd1, 0x83, 0x41, 0xbc, 0x90, 0x98, 0x48, 0x66, 0xb8, 0x78, 0xab, 0xb6, 0x87, 0x26, 0x40,
	0x71, 0x2d, 0xfe, 0x0b, 0xff, 0xb3, 0xe9, 0x4a, 0x65, 0xad, 0xcc, 0x9b, 0xb7, 0xf5, 0xf5, 0xfb,
	0xbe, 0xf7, 0xf5, 0x7d, 0x2f, 0x83, 0x4c, 0x56, 0x8c, 0x57, 0xd3, 0x7d, 0x25, 0xb5, 0x44, 0xd9,
	0xb3, 0x39, 0x2c, 0x36, 0x82, 0xef, 0x34, 0x79, 0x87, 0xfe, 0xa2, 0xe2, 0x54, 0xf3, 0x92, 0x7f,
	0x1c, 0xb8, 0xd2, 0x68, 0x00, 0xf1, 0x5a, 0x30, 0xdc, 0x19, 0x77, 0x26, 0x71, 0x69, 0x3e, 0x4d,
	0x65, 0x25, 0x18, 0x8e, 0x6c, 0x65, 0x25, 0x18, 0xba, 0x84, 0x74, 0xbe, 0x95, 0x87, 0x9d, 0xc6,
	0x71, 0x5d, 0x3c, 0x9e, 0x4c, 0xfd, 0x45, 0x53, 0x7d, 0x50, 0x38, 0xb1, 0x75, 0x7b, 0x22, 0x63,
	0xc8, 0x5d, 0x13, 0xb5, 0x97, 0x3b, 0xc5, 0x51, 0x0e, 0xd1, 0xd2, 0x35, 0x89, 0x96, 0x8c, 0x28,
	0xe8, 0xaf, 0xf7, 0xac, 0x61, 0x23, 0x00, 0x38, 0x5b, 0xd1, 0x2f, 0x5b, 0xf1, 0x39, 0x5b, 0x49,
	0x8b, 0xad, 0xae, 0x67, 0x6b, 0x00, 0xb9, 0x6b, 0x6a, 0x6d, 0x91, 0x11, 0xf4, 0x4b, 0xbe, 0x95,
	0x9f, 0x6d, 0x36, 0x0c, 0xc5, 0x01, 0x4e, 0x94, 0x47, 0xae, 0xa9, 0xd8, 0xb4, 0x51, 0x34, 0xe4,
	0x0e, 0x70, 0x7a, 0xfc, 0xcf, 0x84, 0x23, 0xf1, 0x3f, 0x6f, 0x1b, 0x41, 0xf6, 0x24, 0x94, 0x6e,
	0x4d, 0x95, 0xdc, 0x43, 0xcf, 0x02, 0x8e, 0xa6, 0x6e, 0x20, 0x61, 0x54, 0x53, 0xdc, 0x19, 0xc7,
	0x93, 0x6c, 0x36, 0x9c, 0x36, 0x96, 0x64, 0xea, 0xfb, 0x2f, 0x6b, 0x20, 0xb9, 0x86, 0x6c, 0x45,
	0x05, 0x6b, 0x7b, 0x76, 0x0e, 0x3d, 0x7b, 0x6d, 0x49, 0xb3, 0xaf, 0x18, 0xba, 0xb5, 0x26, 0x9a,
	0x43, 0x6a, 0xb7, 0x01, 0x15, 0x5e, 0x17, 0x6f, 0x0f, 0x8b, 0xe1, 0xd9, 0xbb, 0xa3, 0xd9, 0x39,
	0xa4, 0x36, 0xb9, 0x40, 0xc2, 0xdb, 0xa1, 0x40, 0xc2, 0x8f, 0xda, 0x48, 0xd8, 0x24, 0x03, 0x09,
	0x2f, 0xff, 0x40, 0xc2, 0x8f, 0xde, 0x48, 0xd8, 0xc9, 0x04, 0x12, 0xde, 0x3e, 0x14, 0x7f, 0x8d,
	0x12, 0xdd, 0x41, 0x62, 0x52, 0x40, 0xd8, 0x03, 0x35, 0x92, 0x2b, 0xae, 0xce, 0xdc, 0x9c, 0xc8,
	0x66, 0xc4, 0x01, 0xb9, 0x11, 0x4a, 0x40, 0x6e, 0xe6, 0xf1, 0x70, 0xf1, 0xda, 0xad, 0x7f, 0x0a,
	0x6f, 0x69, 0xfd, 0x57, 0xb8, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x9b, 0xd2, 0x33, 0x24,
	0x04, 0x00, 0x00,
}
