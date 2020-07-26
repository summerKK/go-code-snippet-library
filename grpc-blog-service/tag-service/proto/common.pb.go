// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/common.proto

package proto

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

func init() {
	proto.RegisterType((*Pager)(nil), "proto.Pager")
}

func init() {
	proto.RegisterFile("proto/common.proto", fileDescriptor_1747d3070a2311a0)
}

var fileDescriptor_1747d3070a2311a0 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8e, 0xb1, 0x4b, 0xc5, 0x30,
	0x10, 0x87, 0xa9, 0xb5, 0x62, 0x33, 0x66, 0x2a, 0x88, 0x20, 0x4e, 0x2e, 0x69, 0x06, 0xff, 0x03,
	0xd7, 0x0e, 0x6a, 0x1d, 0x04, 0x97, 0x92, 0xc4, 0x23, 0x2f, 0xd0, 0xf4, 0xc2, 0x25, 0x7d, 0xe5,
	0xf5, 0xaf, 0x7f, 0xf4, 0x96, 0x37, 0x7d, 0xbf, 0xfb, 0xbe, 0xe5, 0x84, 0x4c, 0x84, 0x05, 0xb5,
	0xc3, 0x18, 0x71, 0xe9, 0xf9, 0x90, 0x0d, 0xe3, 0xf5, 0x57, 0x34, 0x5f, 0xc6, 0x03, 0x49, 0x29,
	0xee, 0x93, 0xf1, 0xd0, 0x55, 0x2f, 0xd5, 0x5b, 0x3d, 0xf2, 0x96, 0x4f, 0xa2, 0x3d, 0x38, 0xe5,
	0xb0, 0x43, 0x77, 0xc7, 0xe1, 0xf1, 0x10, 0x3f, 0x61, 0x07, 0xf9, 0x2c, 0x44, 0xc1, 0x62, 0xe6,
	0x89, 0x70, 0xcb, 0x5d, 0xcd, 0xb5, 0x65, 0x33, 0xe2, 0x96, 0x3f, 0xbe, 0xff, 0x3e, 0x7d, 0x28,
	0xa7, 0xd5, 0xf6, 0x0e, 0xa3, 0xce, 0x6b, 0x8c, 0x40, 0xc3, 0xa0, 0x3d, 0x2a, 0x87, 0xff, 0xa0,
	0xf2, 0x12, 0x52, 0x82, 0xa2, 0xe6, 0x60, 0xc9, 0xd0, 0x45, 0x7b, 0x4a, 0x4e, 0xd9, 0x19, 0xbd,
	0xca, 0x40, 0xe7, 0xe0, 0x40, 0x17, 0x73, 0xdb, 0xfc, 0xab, 0x7d, 0x60, 0xbc, 0x5f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0x00, 0xa7, 0x38, 0xe6, 0xcf, 0x00, 0x00, 0x00,
}
