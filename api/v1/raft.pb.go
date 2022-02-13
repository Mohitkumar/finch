// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: api/v1/raft.proto

package api_v1

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LogItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueName string     `protobuf:"bytes,1,opt,name=queueName,proto3" json:"queueName,omitempty"`
	LogRecord *LogRecord `protobuf:"bytes,2,opt,name=logRecord,proto3" json:"logRecord,omitempty"`
}

func (x *LogItem) Reset() {
	*x = LogItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogItem) ProtoMessage() {}

func (x *LogItem) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogItem.ProtoReflect.Descriptor instead.
func (*LogItem) Descriptor() ([]byte, []int) {
	return file_api_v1_raft_proto_rawDescGZIP(), []int{0}
}

func (x *LogItem) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *LogItem) GetLogRecord() *LogRecord {
	if x != nil {
		return x.LogRecord
	}
	return nil
}

type DataRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Record:
	//	*DataRecord_KvItem
	//	*DataRecord_LogItem
	Record isDataRecord_Record `protobuf_oneof:"Record"`
}

func (x *DataRecord) Reset() {
	*x = DataRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataRecord) ProtoMessage() {}

func (x *DataRecord) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataRecord.ProtoReflect.Descriptor instead.
func (*DataRecord) Descriptor() ([]byte, []int) {
	return file_api_v1_raft_proto_rawDescGZIP(), []int{1}
}

func (m *DataRecord) GetRecord() isDataRecord_Record {
	if m != nil {
		return m.Record
	}
	return nil
}

func (x *DataRecord) GetKvItem() *KVItem {
	if x, ok := x.GetRecord().(*DataRecord_KvItem); ok {
		return x.KvItem
	}
	return nil
}

func (x *DataRecord) GetLogItem() *LogItem {
	if x, ok := x.GetRecord().(*DataRecord_LogItem); ok {
		return x.LogItem
	}
	return nil
}

type isDataRecord_Record interface {
	isDataRecord_Record()
}

type DataRecord_KvItem struct {
	KvItem *KVItem `protobuf:"bytes,1,opt,name=kvItem,proto3,oneof"`
}

type DataRecord_LogItem struct {
	LogItem *LogItem `protobuf:"bytes,2,opt,name=logItem,proto3,oneof"`
}

func (*DataRecord_KvItem) isDataRecord_Record() {}

func (*DataRecord_LogItem) isDataRecord_Record() {}

var File_api_v1_raft_proto protoreflect.FileDescriptor

var file_api_v1_raft_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x07, 0x4c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x28, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0x5f, 0x0a, 0x0a,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x21, 0x0a, 0x06, 0x6b, 0x76,
	0x49, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4b, 0x56, 0x49,
	0x74, 0x65, 0x6d, 0x48, 0x00, 0x52, 0x06, 0x6b, 0x76, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x24, 0x0a,
	0x07, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08,
	0x2e, 0x4c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x48, 0x00, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x49,
	0x74, 0x65, 0x6d, 0x42, 0x08, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x42, 0x1e, 0x5a,
	0x1c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x68, 0x69,
	0x74, 0x6b, 0x75, 0x6d, 0x61, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_raft_proto_rawDescOnce sync.Once
	file_api_v1_raft_proto_rawDescData = file_api_v1_raft_proto_rawDesc
)

func file_api_v1_raft_proto_rawDescGZIP() []byte {
	file_api_v1_raft_proto_rawDescOnce.Do(func() {
		file_api_v1_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_raft_proto_rawDescData)
	})
	return file_api_v1_raft_proto_rawDescData
}

var file_api_v1_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_v1_raft_proto_goTypes = []interface{}{
	(*LogItem)(nil),    // 0: LogItem
	(*DataRecord)(nil), // 1: DataRecord
	(*LogRecord)(nil),  // 2: LogRecord
	(*KVItem)(nil),     // 3: KVItem
}
var file_api_v1_raft_proto_depIdxs = []int32{
	2, // 0: LogItem.logRecord:type_name -> LogRecord
	3, // 1: DataRecord.kvItem:type_name -> KVItem
	0, // 2: DataRecord.logItem:type_name -> LogItem
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_v1_raft_proto_init() }
func file_api_v1_raft_proto_init() {
	if File_api_v1_raft_proto != nil {
		return
	}
	file_api_v1_store_proto_init()
	file_api_v1_log_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogItem); i {
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
		file_api_v1_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataRecord); i {
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
	file_api_v1_raft_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*DataRecord_KvItem)(nil),
		(*DataRecord_LogItem)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_raft_proto_goTypes,
		DependencyIndexes: file_api_v1_raft_proto_depIdxs,
		MessageInfos:      file_api_v1_raft_proto_msgTypes,
	}.Build()
	File_api_v1_raft_proto = out.File
	file_api_v1_raft_proto_rawDesc = nil
	file_api_v1_raft_proto_goTypes = nil
	file_api_v1_raft_proto_depIdxs = nil
}
