// Code generated by protoc-gen-gogo.
// source: common/types.proto
// DO NOT EDIT!

/*
Package common is a generated protocol buffer package.

It is generated from these files:
	common/types.proto

It has these top-level messages:
	KVPair
	KI64Pair
*/
//nolint: gas
package util

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Define these here for compatibility but use tmlibs/common.KVPair.
type KVPair struct {
	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *KVPair) Reset()                    { *m = KVPair{} }
func (m *KVPair) String() string            { return proto.CompactTextString(m) }
func (*KVPair) ProtoMessage()               {}
func (*KVPair) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{0} }

func (m *KVPair) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KVPair) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Define these here for compatibility but use tmlibs/common.KI64Pair.
type KI64Pair struct {
	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *KI64Pair) Reset()                    { *m = KI64Pair{} }
func (m *KI64Pair) String() string            { return proto.CompactTextString(m) }
func (*KI64Pair) ProtoMessage()               {}
func (*KI64Pair) Descriptor() ([]byte, []int) { return fileDescriptorTypes, []int{1} }

func (m *KI64Pair) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KI64Pair) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*KVPair)(nil), "common.KVPair")
	proto.RegisterType((*KI64Pair)(nil), "common.KI64Pair")
}

func init() { proto.RegisterFile("common/types.proto", fileDescriptorTypes) }

var fileDescriptorTypes = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x83, 0x88, 0x49, 0xe9, 0xa6, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0xa7,
	0xe7, 0xa7, 0xe7, 0xeb, 0x83, 0xa5, 0x93, 0x4a, 0xd3, 0xc0, 0x3c, 0x30, 0x07, 0xcc, 0x82, 0x68,
	0x53, 0x32, 0xe0, 0x62, 0xf3, 0x0e, 0x0b, 0x48, 0xcc, 0x2c, 0x12, 0x12, 0xe0, 0x62, 0xce, 0x4e,
	0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0x31, 0x85, 0x44, 0xb8, 0x58, 0xcb, 0x12,
	0x73, 0x4a, 0x53, 0x25, 0x98, 0xc0, 0x62, 0x10, 0x8e, 0x92, 0x11, 0x17, 0x87, 0xb7, 0xa7, 0x99,
	0x09, 0x31, 0x7a, 0x98, 0xa1, 0x7a, 0x92, 0xd8, 0xc0, 0x96, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x5c, 0xb8, 0x46, 0xc5, 0xb9, 0x00, 0x00, 0x00,
}
