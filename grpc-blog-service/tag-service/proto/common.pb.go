// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/common.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

type Pager struct {
	Page                 int64    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize             int64    `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	TotalRows            int64    `protobuf:"varint,3,opt,name=total_rows,json=totalRows,proto3" json:"total_rows,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pager) Reset()         { *m = Pager{} }
func (m *Pager) String() string { return proto.CompactTextString(m) }
func (*Pager) ProtoMessage()    {}
func (*Pager) Descriptor() ([]byte, []int) {
	return fileDescriptor_1747d3070a2311a0, []int{0}
}

func (m *Pager) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pager.Unmarshal(m, b)
}
func (m *Pager) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pager.Marshal(b, m, deterministic)
}
func (m *Pager) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pager.Merge(m, src)
}
func (m *Pager) XXX_Size() int {
	return xxx_messageInfo_Pager.Size(m)
}
func (m *Pager) XXX_DiscardUnknown() {
	xxx_messageInfo_Pager.DiscardUnknown(m)
}

var xxx_messageInfo_Pager proto.InternalMessageInfo

func (m *Pager) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *Pager) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *Pager) GetTotalRows() int64 {
	if m != nil {
		return m.TotalRows
	}
	return 0
}

type Error struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Detail               *any.Any `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_1747d3070a2311a0, []int{1}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Error) GetDetail() *any.Any {
	if m != nil {
		return m.Detail
	}
	return nil
}

func init() {
	proto.RegisterType((*Pager)(nil), "proto.Pager")
	proto.RegisterType((*Error)(nil), "proto.Error")
}

func init() {
	proto.RegisterFile("proto/common.proto", fileDescriptor_1747d3070a2311a0)
}

var fileDescriptor_1747d3070a2311a0 = []byte{
	// 255 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8e, 0x41, 0x4b, 0x03, 0x31,
	0x10, 0x85, 0xa9, 0x75, 0xab, 0x8d, 0xb7, 0xe0, 0x61, 0x55, 0x04, 0xe9, 0xc9, 0x83, 0xbb, 0x01,
	0xfd, 0x05, 0x0a, 0x9e, 0x7a, 0x50, 0xd7, 0x83, 0xe0, 0xa5, 0x64, 0xd3, 0x31, 0x06, 0x36, 0x3b,
	0x61, 0x92, 0xb5, 0x6c, 0x7f, 0xbd, 0x64, 0xd6, 0xe2, 0xe9, 0xbd, 0x37, 0x33, 0x7c, 0x6f, 0x84,
	0x0c, 0x84, 0x09, 0x95, 0x41, 0xef, 0xb1, 0xaf, 0x39, 0xc8, 0x82, 0xe5, 0xf2, 0xc2, 0x22, 0xda,
	0x0e, 0x14, 0xa7, 0x76, 0xf8, 0x52, 0xba, 0x1f, 0xa7, 0x8b, 0xd5, 0x87, 0x28, 0x5e, 0xb5, 0x05,
	0x92, 0x52, 0x1c, 0x07, 0x6d, 0xa1, 0x9c, 0xdd, 0xcc, 0x6e, 0xe7, 0x0d, 0x7b, 0x79, 0x25, 0x96,
	0x59, 0x37, 0xd1, 0xed, 0xa1, 0x3c, 0xe2, 0xc5, 0x69, 0x1e, 0xbc, 0xbb, 0x3d, 0xc8, 0x6b, 0x21,
	0x12, 0x26, 0xdd, 0x6d, 0x08, 0x77, 0xb1, 0x9c, 0xf3, 0x76, 0xc9, 0x93, 0x06, 0x77, 0x71, 0x65,
	0x44, 0xf1, 0x4c, 0x84, 0x0c, 0x36, 0xb8, 0x9d, 0xc0, 0x45, 0xc3, 0x5e, 0x96, 0xe2, 0xc4, 0x43,
	0x8c, 0xb9, 0x2f, 0x63, 0x97, 0xcd, 0x21, 0xca, 0x3b, 0xb1, 0xd8, 0x42, 0xd2, 0xae, 0x63, 0xe2,
	0xd9, 0xfd, 0x79, 0x3d, 0xfd, 0x5e, 0x1f, 0x7e, 0xaf, 0x1f, 0xfb, 0xb1, 0xf9, 0xbb, 0x79, 0x7a,
	0xfb, 0x7c, 0xb1, 0x2e, 0x7d, 0x0f, 0x6d, 0x6d, 0xd0, 0xab, 0x38, 0x78, 0x0f, 0xb4, 0x5e, 0x2b,
	0x8b, 0x55, 0xae, 0xa9, 0x62, 0xef, 0x42, 0x80, 0x54, 0x75, 0xae, 0x25, 0x4d, 0xa3, 0xb2, 0x14,
	0x4c, 0xd5, 0x76, 0x68, 0xab, 0x08, 0xf4, 0xe3, 0x0c, 0xa8, 0xa4, 0xff, 0xfd, 0xd4, 0xb0, 0x60,
	0x79, 0xf8, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x59, 0x4f, 0x13, 0xe0, 0x4f, 0x01, 0x00, 0x00,
}
