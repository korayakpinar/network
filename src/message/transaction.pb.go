// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: src/message/transaction.proto

package message

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

type EncryptedTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *TransactionHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Body   *TransactionBody   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *EncryptedTransaction) Reset() {
	*x = EncryptedTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_message_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedTransaction) ProtoMessage() {}

func (x *EncryptedTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_src_message_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedTransaction.ProtoReflect.Descriptor instead.
func (*EncryptedTransaction) Descriptor() ([]byte, []int) {
	return file_src_message_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *EncryptedTransaction) GetHeader() *TransactionHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *EncryptedTransaction) GetBody() *TransactionBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type TransactionHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash    string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	GammaG2 []byte   `protobuf:"bytes,2,opt,name=gammaG2,proto3" json:"gammaG2,omitempty"`
	PkIDs   []uint64 `protobuf:"varint,3,rep,packed,name=pkIDs,proto3" json:"pkIDs,omitempty"` // List of public key IDs
}

func (x *TransactionHeader) Reset() {
	*x = TransactionHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_message_transaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionHeader) ProtoMessage() {}

func (x *TransactionHeader) ProtoReflect() protoreflect.Message {
	mi := &file_src_message_transaction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionHeader.ProtoReflect.Descriptor instead.
func (*TransactionHeader) Descriptor() ([]byte, []int) {
	return file_src_message_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *TransactionHeader) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TransactionHeader) GetGammaG2() []byte {
	if x != nil {
		return x.GammaG2
	}
	return nil
}

func (x *TransactionHeader) GetPkIDs() []uint64 {
	if x != nil {
		return x.PkIDs
	}
	return nil
}

type TransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sa1     []byte `protobuf:"bytes,1,opt,name=sa1,proto3" json:"sa1,omitempty"`         // String array of length 2
	Sa2     []byte `protobuf:"bytes,2,opt,name=sa2,proto3" json:"sa2,omitempty"`         // String array of length 6
	Iv      []byte `protobuf:"bytes,3,opt,name=iv,proto3" json:"iv,omitempty"`           // Initialization vector for AES256, 16 bytes long
	EncText []byte `protobuf:"bytes,4,opt,name=encText,proto3" json:"encText,omitempty"` // Encrypted text
	T       uint64 `protobuf:"varint,5,opt,name=t,proto3" json:"t,omitempty"`            // Unsigned integer field
}

func (x *TransactionBody) Reset() {
	*x = TransactionBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_message_transaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionBody) ProtoMessage() {}

func (x *TransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_src_message_transaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionBody.ProtoReflect.Descriptor instead.
func (*TransactionBody) Descriptor() ([]byte, []int) {
	return file_src_message_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionBody) GetSa1() []byte {
	if x != nil {
		return x.Sa1
	}
	return nil
}

func (x *TransactionBody) GetSa2() []byte {
	if x != nil {
		return x.Sa2
	}
	return nil
}

func (x *TransactionBody) GetIv() []byte {
	if x != nil {
		return x.Iv
	}
	return nil
}

func (x *TransactionBody) GetEncText() []byte {
	if x != nil {
		return x.EncText
	}
	return nil
}

func (x *TransactionBody) GetT() uint64 {
	if x != nil {
		return x.T
	}
	return 0
}

var File_src_message_transaction_proto protoreflect.FileDescriptor

var file_src_message_transaction_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x72, 0x63, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x68, 0x0a, 0x14, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x57, 0x0a, 0x11, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x67, 0x61, 0x6d, 0x6d, 0x61, 0x47, 0x32, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x6b, 0x49, 0x44, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x05, 0x70, 0x6b, 0x49,
	0x44, 0x73, 0x22, 0x6d, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x31, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x32, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x61, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x76, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x63,
	0x54, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x65, 0x6e, 0x63, 0x54,
	0x65, 0x78, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x01,
	0x74, 0x42, 0x0d, 0x5a, 0x0b, 0x73, 0x72, 0x63, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_message_transaction_proto_rawDescOnce sync.Once
	file_src_message_transaction_proto_rawDescData = file_src_message_transaction_proto_rawDesc
)

func file_src_message_transaction_proto_rawDescGZIP() []byte {
	file_src_message_transaction_proto_rawDescOnce.Do(func() {
		file_src_message_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_message_transaction_proto_rawDescData)
	})
	return file_src_message_transaction_proto_rawDescData
}

var file_src_message_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_src_message_transaction_proto_goTypes = []any{
	(*EncryptedTransaction)(nil), // 0: EncryptedTransaction
	(*TransactionHeader)(nil),    // 1: TransactionHeader
	(*TransactionBody)(nil),      // 2: TransactionBody
}
var file_src_message_transaction_proto_depIdxs = []int32{
	1, // 0: EncryptedTransaction.header:type_name -> TransactionHeader
	2, // 1: EncryptedTransaction.body:type_name -> TransactionBody
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_src_message_transaction_proto_init() }
func file_src_message_transaction_proto_init() {
	if File_src_message_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_message_transaction_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*EncryptedTransaction); i {
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
		file_src_message_transaction_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TransactionHeader); i {
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
		file_src_message_transaction_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TransactionBody); i {
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
			RawDescriptor: file_src_message_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_src_message_transaction_proto_goTypes,
		DependencyIndexes: file_src_message_transaction_proto_depIdxs,
		MessageInfos:      file_src_message_transaction_proto_msgTypes,
	}.Build()
	File_src_message_transaction_proto = out.File
	file_src_message_transaction_proto_rawDesc = nil
	file_src_message_transaction_proto_goTypes = nil
	file_src_message_transaction_proto_depIdxs = nil
}
