// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.19.6
// source: chat_db/chat-db.proto

package chat_db

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReadMostRecentResponse struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	ChannelId          int64                  `protobuf:"varint,1,opt,name=channelId,proto3" json:"channelId,omitempty"`
	Message            string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	UserId             int64                  `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	TimePostedUnixTime int64                  `protobuf:"varint,4,opt,name=timePostedUnixTime,proto3" json:"timePostedUnixTime,omitempty"`
	Offset             int64                  `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *ReadMostRecentResponse) Reset() {
	*x = ReadMostRecentResponse{}
	mi := &file_chat_db_chat_db_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadMostRecentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMostRecentResponse) ProtoMessage() {}

func (x *ReadMostRecentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chat_db_chat_db_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMostRecentResponse.ProtoReflect.Descriptor instead.
func (*ReadMostRecentResponse) Descriptor() ([]byte, []int) {
	return file_chat_db_chat_db_proto_rawDescGZIP(), []int{0}
}

func (x *ReadMostRecentResponse) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *ReadMostRecentResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ReadMostRecentResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReadMostRecentResponse) GetTimePostedUnixTime() int64 {
	if x != nil {
		return x.TimePostedUnixTime
	}
	return 0
}

func (x *ReadMostRecentResponse) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type ReadMostRecentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     int64                  `protobuf:"varint,1,opt,name=channelId,proto3" json:"channelId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadMostRecentRequest) Reset() {
	*x = ReadMostRecentRequest{}
	mi := &file_chat_db_chat_db_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadMostRecentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMostRecentRequest) ProtoMessage() {}

func (x *ReadMostRecentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_db_chat_db_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMostRecentRequest.ProtoReflect.Descriptor instead.
func (*ReadMostRecentRequest) Descriptor() ([]byte, []int) {
	return file_chat_db_chat_db_proto_rawDescGZIP(), []int{1}
}

func (x *ReadMostRecentRequest) GetChannelId() int64 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

var File_chat_db_chat_db_proto protoreflect.FileDescriptor

const file_chat_db_chat_db_proto_rawDesc = "" +
	"\n" +
	"\x15chat_db/chat-db.proto\x12\achat_db\"\xb0\x01\n" +
	"\x16ReadMostRecentResponse\x12\x1c\n" +
	"\tchannelId\x18\x01 \x01(\x03R\tchannelId\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12\x16\n" +
	"\x06userId\x18\x03 \x01(\x03R\x06userId\x12.\n" +
	"\x12timePostedUnixTime\x18\x04 \x01(\x03R\x12timePostedUnixTime\x12\x16\n" +
	"\x06offset\x18\x05 \x01(\x03R\x06offset\"5\n" +
	"\x15ReadMostRecentRequest\x12\x1c\n" +
	"\tchannelId\x18\x01 \x01(\x03R\tchannelId2_\n" +
	"\x06chatdb\x12U\n" +
	"\x0eReadMostRecent\x12\x1e.chat_db.ReadMostRecentRequest\x1a\x1f.chat_db.ReadMostRecentResponse\"\x000\x01B5Z3github.com/WadeCappa/real_time_chat/chat-db/chat-dbb\x06proto3"

var (
	file_chat_db_chat_db_proto_rawDescOnce sync.Once
	file_chat_db_chat_db_proto_rawDescData []byte
)

func file_chat_db_chat_db_proto_rawDescGZIP() []byte {
	file_chat_db_chat_db_proto_rawDescOnce.Do(func() {
		file_chat_db_chat_db_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_chat_db_chat_db_proto_rawDesc), len(file_chat_db_chat_db_proto_rawDesc)))
	})
	return file_chat_db_chat_db_proto_rawDescData
}

var file_chat_db_chat_db_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_chat_db_chat_db_proto_goTypes = []any{
	(*ReadMostRecentResponse)(nil), // 0: chat_db.ReadMostRecentResponse
	(*ReadMostRecentRequest)(nil),  // 1: chat_db.ReadMostRecentRequest
}
var file_chat_db_chat_db_proto_depIdxs = []int32{
	1, // 0: chat_db.chatdb.ReadMostRecent:input_type -> chat_db.ReadMostRecentRequest
	0, // 1: chat_db.chatdb.ReadMostRecent:output_type -> chat_db.ReadMostRecentResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_chat_db_chat_db_proto_init() }
func file_chat_db_chat_db_proto_init() {
	if File_chat_db_chat_db_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_chat_db_chat_db_proto_rawDesc), len(file_chat_db_chat_db_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_db_chat_db_proto_goTypes,
		DependencyIndexes: file_chat_db_chat_db_proto_depIdxs,
		MessageInfos:      file_chat_db_chat_db_proto_msgTypes,
	}.Build()
	File_chat_db_chat_db_proto = out.File
	file_chat_db_chat_db_proto_goTypes = nil
	file_chat_db_chat_db_proto_depIdxs = nil
}
