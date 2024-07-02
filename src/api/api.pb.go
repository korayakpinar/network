// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: src/api/api.proto

package api

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

// DecryptData function messages
type DecryptDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enc string   `protobuf:"bytes,1,opt,name=enc,proto3" json:"enc,omitempty"`
	Pks []string `protobuf:"bytes,2,rep,name=pks,proto3" json:"pks,omitempty"`
	Sa1 string   `protobuf:"bytes,3,opt,name=sa1,proto3" json:"sa1,omitempty"`
	Sa2 string   `protobuf:"bytes,4,opt,name=sa2,proto3" json:"sa2,omitempty"`
	Iv  string   `protobuf:"bytes,5,opt,name=iv,proto3" json:"iv,omitempty"`
	T   uint64   `protobuf:"varint,6,opt,name=t,proto3" json:"t,omitempty"`
	N   uint64   `protobuf:"varint,7,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *DecryptDataRequest) Reset() {
	*x = DecryptDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecryptDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecryptDataRequest) ProtoMessage() {}

func (x *DecryptDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecryptDataRequest.ProtoReflect.Descriptor instead.
func (*DecryptDataRequest) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{0}
}

func (x *DecryptDataRequest) GetEnc() string {
	if x != nil {
		return x.Enc
	}
	return ""
}

func (x *DecryptDataRequest) GetPks() []string {
	if x != nil {
		return x.Pks
	}
	return nil
}

func (x *DecryptDataRequest) GetSa1() string {
	if x != nil {
		return x.Sa1
	}
	return ""
}

func (x *DecryptDataRequest) GetSa2() string {
	if x != nil {
		return x.Sa2
	}
	return ""
}

func (x *DecryptDataRequest) GetIv() string {
	if x != nil {
		return x.Iv
	}
	return ""
}

func (x *DecryptDataRequest) GetT() uint64 {
	if x != nil {
		return x.T
	}
	return 0
}

func (x *DecryptDataRequest) GetN() uint64 {
	if x != nil {
		return x.N
	}
	return 0
}

type DecryptDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *DecryptDataResponse) Reset() {
	*x = DecryptDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecryptDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecryptDataResponse) ProtoMessage() {}

func (x *DecryptDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecryptDataResponse.ProtoReflect.Descriptor instead.
func (*DecryptDataResponse) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{1}
}

func (x *DecryptDataResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

// PartialDecrypt function messages
type PartialDecryptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GammaG2 string `protobuf:"bytes,1,opt,name=gammaG2,proto3" json:"gammaG2,omitempty"`
}

func (x *PartialDecryptRequest) Reset() {
	*x = PartialDecryptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialDecryptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialDecryptRequest) ProtoMessage() {}

func (x *PartialDecryptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialDecryptRequest.ProtoReflect.Descriptor instead.
func (*PartialDecryptRequest) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{2}
}

func (x *PartialDecryptRequest) GetGammaG2() string {
	if x != nil {
		return x.GammaG2
	}
	return ""
}

type PartialDecryptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *PartialDecryptResponse) Reset() {
	*x = PartialDecryptResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialDecryptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialDecryptResponse) ProtoMessage() {}

func (x *PartialDecryptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialDecryptResponse.ProtoReflect.Descriptor instead.
func (*PartialDecryptResponse) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{3}
}

func (x *PartialDecryptResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

// VerifyPart function messages
type VerifyPartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pk      string `protobuf:"bytes,1,opt,name=pk,proto3" json:"pk,omitempty"`
	GammaG2 string `protobuf:"bytes,2,opt,name=gammaG2,proto3" json:"gammaG2,omitempty"`
	PartDec string `protobuf:"bytes,3,opt,name=partDec,proto3" json:"partDec,omitempty"`
}

func (x *VerifyPartRequest) Reset() {
	*x = VerifyPartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyPartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyPartRequest) ProtoMessage() {}

func (x *VerifyPartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[4]
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
	return file_src_api_api_proto_rawDescGZIP(), []int{4}
}

func (x *VerifyPartRequest) GetPk() string {
	if x != nil {
		return x.Pk
	}
	return ""
}

func (x *VerifyPartRequest) GetGammaG2() string {
	if x != nil {
		return x.GammaG2
	}
	return ""
}

func (x *VerifyPartRequest) GetPartDec() string {
	if x != nil {
		return x.PartDec
	}
	return ""
}

type VerifyPartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *VerifyPartResponse) Reset() {
	*x = VerifyPartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_api_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyPartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyPartResponse) ProtoMessage() {}

func (x *VerifyPartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_api_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyPartResponse.ProtoReflect.Descriptor instead.
func (*VerifyPartResponse) Descriptor() ([]byte, []int) {
	return file_src_api_api_proto_rawDescGZIP(), []int{5}
}

func (x *VerifyPartResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_src_api_api_proto protoreflect.FileDescriptor

var file_src_api_api_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x72, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x88, 0x01, 0x0a, 0x12, 0x44, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x63, 0x12, 0x10, 0x0a, 0x03,
	0x70, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x70, 0x6b, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x61, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x61, 0x31,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73,
	0x61, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x76, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01, 0x74,
	0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01, 0x6e, 0x22, 0x2d,
	0x0a, 0x13, 0x44, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x31, 0x0a,
	0x15, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x44, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47,
	0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32,
	0x22, 0x30, 0x0a, 0x16, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x44, 0x65, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x57, 0x0a, 0x11, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x61, 0x72, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x70, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61,
	0x47, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47,
	0x32, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x72, 0x74, 0x44, 0x65, 0x63, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x72, 0x74, 0x44, 0x65, 0x63, 0x22, 0x2c, 0x0a, 0x12, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x09, 0x5a, 0x07, 0x73, 0x72, 0x63,
	0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_api_api_proto_rawDescOnce sync.Once
	file_src_api_api_proto_rawDescData = file_src_api_api_proto_rawDesc
)

func file_src_api_api_proto_rawDescGZIP() []byte {
	file_src_api_api_proto_rawDescOnce.Do(func() {
		file_src_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_api_api_proto_rawDescData)
	})
	return file_src_api_api_proto_rawDescData
}

var file_src_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_src_api_api_proto_goTypes = []any{
	(*DecryptDataRequest)(nil),     // 0: DecryptDataRequest
	(*DecryptDataResponse)(nil),    // 1: DecryptDataResponse
	(*PartialDecryptRequest)(nil),  // 2: PartialDecryptRequest
	(*PartialDecryptResponse)(nil), // 3: PartialDecryptResponse
	(*VerifyPartRequest)(nil),      // 4: VerifyPartRequest
	(*VerifyPartResponse)(nil),     // 5: VerifyPartResponse
}
var file_src_api_api_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_src_api_api_proto_init() }
func file_src_api_api_proto_init() {
	if File_src_api_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_api_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DecryptDataRequest); i {
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
		file_src_api_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DecryptDataResponse); i {
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
		file_src_api_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PartialDecryptRequest); i {
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
		file_src_api_api_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*PartialDecryptResponse); i {
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
		file_src_api_api_proto_msgTypes[4].Exporter = func(v any, i int) any {
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
		file_src_api_api_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*VerifyPartResponse); i {
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
			RawDescriptor: file_src_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_src_api_api_proto_goTypes,
		DependencyIndexes: file_src_api_api_proto_depIdxs,
		MessageInfos:      file_src_api_api_proto_msgTypes,
	}.Build()
	File_src_api_api_proto = out.File
	file_src_api_api_proto_rawDesc = nil
	file_src_api_api_proto_goTypes = nil
	file_src_api_api_proto_depIdxs = nil
}
