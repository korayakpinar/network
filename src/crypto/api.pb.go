// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: api.proto

package crypto

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

type VerifyPartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pk      []byte `protobuf:"bytes,1,opt,name=pk,proto3" json:"pk,omitempty"`
	GammaG2 []byte `protobuf:"bytes,2,opt,name=gamma_g2,json=gammaG2,proto3" json:"gamma_g2,omitempty"`
	PartDec []byte `protobuf:"bytes,3,opt,name=part_dec,json=partDec,proto3" json:"part_dec,omitempty"`
}

func (x *VerifyPartRequest) Reset() {
	*x = VerifyPartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyPartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyPartRequest) ProtoMessage() {}

func (x *VerifyPartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyPartRequest.ProtoReflect.Descriptor instead.
func (*VerifyPartRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyPartRequest) GetPk() []byte {
	if x != nil {
		return x.Pk
	}
	return nil
}

func (x *VerifyPartRequest) GetGammaG2() []byte {
	if x != nil {
		return x.GammaG2
	}
	return nil
}

func (x *VerifyPartRequest) GetPartDec() []byte {
	if x != nil {
		return x.PartDec
	}
	return nil
}

type IsValidRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pk []byte `protobuf:"bytes,1,opt,name=pk,proto3" json:"pk,omitempty"`
	N  uint64 `protobuf:"varint,2,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *IsValidRequest) Reset() {
	*x = IsValidRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsValidRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsValidRequest) ProtoMessage() {}

func (x *IsValidRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsValidRequest.ProtoReflect.Descriptor instead.
func (*IsValidRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
}

func (x *IsValidRequest) GetPk() []byte {
	if x != nil {
		return x.Pk
	}
	return nil
}

func (x *IsValidRequest) GetN() uint64 {
	if x != nil {
		return x.N
	}
	return 0
}

type DecryptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enc     []byte            `protobuf:"bytes,1,opt,name=enc,proto3" json:"enc,omitempty"`
	Parts   map[uint64][]byte `protobuf:"bytes,2,rep,name=parts,proto3" json:"parts,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	GammaG2 []byte            `protobuf:"bytes,3,opt,name=gamma_g2,json=gammaG2,proto3" json:"gamma_g2,omitempty"`
	Sa1     []byte            `protobuf:"bytes,4,opt,name=sa1,proto3" json:"sa1,omitempty"`
	Sa2     []byte            `protobuf:"bytes,5,opt,name=sa2,proto3" json:"sa2,omitempty"`
	Iv      []byte            `protobuf:"bytes,6,opt,name=iv,proto3" json:"iv,omitempty"`
	T       uint64            `protobuf:"varint,7,opt,name=t,proto3" json:"t,omitempty"`
	N       uint64            `protobuf:"varint,8,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *DecryptRequest) Reset() {
	*x = DecryptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecryptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecryptRequest) ProtoMessage() {}

func (x *DecryptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecryptRequest.ProtoReflect.Descriptor instead.
func (*DecryptRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{2}
}

func (x *DecryptRequest) GetEnc() []byte {
	if x != nil {
		return x.Enc
	}
	return nil
}

func (x *DecryptRequest) GetParts() map[uint64][]byte {
	if x != nil {
		return x.Parts
	}
	return nil
}

func (x *DecryptRequest) GetGammaG2() []byte {
	if x != nil {
		return x.GammaG2
	}
	return nil
}

func (x *DecryptRequest) GetSa1() []byte {
	if x != nil {
		return x.Sa1
	}
	return nil
}

func (x *DecryptRequest) GetSa2() []byte {
	if x != nil {
		return x.Sa2
	}
	return nil
}

func (x *DecryptRequest) GetIv() []byte {
	if x != nil {
		return x.Iv
	}
	return nil
}

func (x *DecryptRequest) GetT() uint64 {
	if x != nil {
		return x.T
	}
	return 0
}

func (x *DecryptRequest) GetN() uint64 {
	if x != nil {
		return x.N
	}
	return 0
}

type EncryptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg []byte `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	T   uint64 `protobuf:"varint,2,opt,name=t,proto3" json:"t,omitempty"`
	N   uint64 `protobuf:"varint,3,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *EncryptRequest) Reset() {
	*x = EncryptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptRequest) ProtoMessage() {}

func (x *EncryptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptRequest.ProtoReflect.Descriptor instead.
func (*EncryptRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{3}
}

func (x *EncryptRequest) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *EncryptRequest) GetT() uint64 {
	if x != nil {
		return x.T
	}
	return 0
}

func (x *EncryptRequest) GetN() uint64 {
	if x != nil {
		return x.N
	}
	return 0
}

type EncryptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enc     []byte `protobuf:"bytes,1,opt,name=enc,proto3" json:"enc,omitempty"`
	Sa1     []byte `protobuf:"bytes,2,opt,name=sa1,proto3" json:"sa1,omitempty"`
	Sa2     []byte `protobuf:"bytes,3,opt,name=sa2,proto3" json:"sa2,omitempty"`
	Iv      []byte `protobuf:"bytes,4,opt,name=iv,proto3" json:"iv,omitempty"`
	GammaG2 []byte `protobuf:"bytes,5,opt,name=gamma_g2,json=gammaG2,proto3" json:"gamma_g2,omitempty"`
}

func (x *EncryptResponse) Reset() {
	*x = EncryptResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptResponse) ProtoMessage() {}

func (x *EncryptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptResponse.ProtoReflect.Descriptor instead.
func (*EncryptResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{4}
}

func (x *EncryptResponse) GetEnc() []byte {
	if x != nil {
		return x.Enc
	}
	return nil
}

func (x *EncryptResponse) GetSa1() []byte {
	if x != nil {
		return x.Sa1
	}
	return nil
}

func (x *EncryptResponse) GetSa2() []byte {
	if x != nil {
		return x.Sa2
	}
	return nil
}

func (x *EncryptResponse) GetIv() []byte {
	if x != nil {
		return x.Iv
	}
	return nil
}

func (x *EncryptResponse) GetGammaG2() []byte {
	if x != nil {
		return x.GammaG2
	}
	return nil
}

type PartDecRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sk      []byte `protobuf:"bytes,1,opt,name=sk,proto3" json:"sk,omitempty"`
	GammaG2 []byte `protobuf:"bytes,2,opt,name=gamma_g2,json=gammaG2,proto3" json:"gamma_g2,omitempty"`
}

func (x *PartDecRequest) Reset() {
	*x = PartDecRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartDecRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartDecRequest) ProtoMessage() {}

func (x *PartDecRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartDecRequest.ProtoReflect.Descriptor instead.
func (*PartDecRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{5}
}

func (x *PartDecRequest) GetSk() []byte {
	if x != nil {
		return x.Sk
	}
	return nil
}

func (x *PartDecRequest) GetGammaG2() []byte {
	if x != nil {
		return x.GammaG2
	}
	return nil
}

type PKRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sk []byte `protobuf:"bytes,1,opt,name=sk,proto3" json:"sk,omitempty"`
	Id uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	N  uint64 `protobuf:"varint,3,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *PKRequest) Reset() {
	*x = PKRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PKRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PKRequest) ProtoMessage() {}

func (x *PKRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PKRequest.ProtoReflect.Descriptor instead.
func (*PKRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{6}
}

func (x *PKRequest) GetSk() []byte {
	if x != nil {
		return x.Sk
	}
	return nil
}

func (x *PKRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PKRequest) GetN() uint64 {
	if x != nil {
		return x.N
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []byte `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{7}
}

func (x *Response) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x11, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x70, 0x6b,
	0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x5f, 0x67, 0x32, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32, 0x12, 0x19, 0x0a, 0x08, 0x70,
	0x61, 0x72, 0x74, 0x5f, 0x64, 0x65, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70,
	0x61, 0x72, 0x74, 0x44, 0x65, 0x63, 0x22, 0x2e, 0x0a, 0x0e, 0x49, 0x73, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x70, 0x6b, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x01, 0x6e, 0x22, 0xf9, 0x01, 0x0a, 0x0e, 0x44, 0x65, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x63,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x65, 0x6e, 0x63, 0x12, 0x30, 0x0a, 0x05, 0x70,
	0x61, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x44, 0x65, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x70, 0x61, 0x72, 0x74, 0x73, 0x12, 0x19, 0x0a,
	0x08, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x5f, 0x67, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x31, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61,
	0x32, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x32, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x76, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x76, 0x12, 0x0c, 0x0a, 0x01,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01, 0x6e, 0x1a, 0x38, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x3e, 0x0a, 0x0e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x0c, 0x0a, 0x01, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x01, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x01, 0x6e, 0x22, 0x72, 0x0a, 0x0f, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x63, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x65, 0x6e, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x31, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x32,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x76, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x76, 0x12, 0x19, 0x0a, 0x08, 0x67,
	0x61, 0x6d, 0x6d, 0x61, 0x5f, 0x67, 0x32, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x67,
	0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32, 0x22, 0x3b, 0x0a, 0x0e, 0x50, 0x61, 0x72, 0x74, 0x44, 0x65,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x73, 0x6b, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x6d, 0x6d,
	0x61, 0x5f, 0x67, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x67, 0x61, 0x6d, 0x6d,
	0x61, 0x47, 0x32, 0x22, 0x39, 0x0a, 0x09, 0x50, 0x4b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x73, 0x6b,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01, 0x6e, 0x22, 0x22,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_proto_goTypes = []any{
	(*VerifyPartRequest)(nil), // 0: VerifyPartRequest
	(*IsValidRequest)(nil),    // 1: IsValidRequest
	(*DecryptRequest)(nil),    // 2: DecryptRequest
	(*EncryptRequest)(nil),    // 3: EncryptRequest
	(*EncryptResponse)(nil),   // 4: EncryptResponse
	(*PartDecRequest)(nil),    // 5: PartDecRequest
	(*PKRequest)(nil),         // 6: PKRequest
	(*Response)(nil),          // 7: Response
	nil,                       // 8: DecryptRequest.PartsEntry
}
var file_api_proto_depIdxs = []int32{
	8, // 0: DecryptRequest.parts:type_name -> DecryptRequest.PartsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*VerifyPartRequest); i {
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
		file_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*IsValidRequest); i {
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
		file_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DecryptRequest); i {
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
		file_api_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EncryptRequest); i {
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
		file_api_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EncryptResponse); i {
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
		file_api_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*PartDecRequest); i {
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
		file_api_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*PKRequest); i {
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
		file_api_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*Response); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		MessageInfos:      file_api_proto_msgTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
