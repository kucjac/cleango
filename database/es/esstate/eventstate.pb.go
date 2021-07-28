// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: eventstate/eventstate.proto

package esstate

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

// EventUnhandled is an event message which states that an event is marked as unhandled.
type EventUnhandled struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType       string `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Timestamp       int64  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	MaxFailures     int32  `protobuf:"varint,3,opt,name=max_failures,json=maxFailures,proto3" json:"max_failures,omitempty"`
	MinFailInterval int64  `protobuf:"varint,4,opt,name=min_fail_interval,json=minFailInterval,proto3" json:"min_fail_interval,omitempty"`
	MaxHandlingTime int64  `protobuf:"varint,5,opt,name=max_handling_time,json=maxHandlingTime,proto3" json:"max_handling_time,omitempty"`
}

func (x *EventUnhandled) Reset() {
	*x = EventUnhandled{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstate_eventstate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventUnhandled) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventUnhandled) ProtoMessage() {}

func (x *EventUnhandled) ProtoReflect() protoreflect.Message {
	mi := &file_eventstate_eventstate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventUnhandled.ProtoReflect.Descriptor instead.
func (*EventUnhandled) Descriptor() ([]byte, []int) {
	return file_eventstate_eventstate_proto_rawDescGZIP(), []int{0}
}

func (x *EventUnhandled) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *EventUnhandled) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *EventUnhandled) GetMaxFailures() int32 {
	if x != nil {
		return x.MaxFailures
	}
	return 0
}

func (x *EventUnhandled) GetMinFailInterval() int64 {
	if x != nil {
		return x.MinFailInterval
	}
	return 0
}

func (x *EventUnhandled) GetMaxHandlingTime() int64 {
	if x != nil {
		return x.MaxHandlingTime
	}
	return 0
}

// EventHandlingStarted is an event message occurred when given handler just
// started handling an event.
type EventHandlingStarted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HandlerName string `protobuf:"bytes,1,opt,name=handler_name,json=handlerName,proto3" json:"handler_name,omitempty"`
}

func (x *EventHandlingStarted) Reset() {
	*x = EventHandlingStarted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstate_eventstate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventHandlingStarted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventHandlingStarted) ProtoMessage() {}

func (x *EventHandlingStarted) ProtoReflect() protoreflect.Message {
	mi := &file_eventstate_eventstate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventHandlingStarted.ProtoReflect.Descriptor instead.
func (*EventHandlingStarted) Descriptor() ([]byte, []int) {
	return file_eventstate_eventstate_proto_rawDescGZIP(), []int{1}
}

func (x *EventHandlingStarted) GetHandlerName() string {
	if x != nil {
		return x.HandlerName
	}
	return ""
}

// EventHandlingFinished is an event message occurred when given handler just
// finished successfully handling an event.
type EventHandlingFinished struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HandlerName string `protobuf:"bytes,1,opt,name=handler_name,json=handlerName,proto3" json:"handler_name,omitempty"`
}

func (x *EventHandlingFinished) Reset() {
	*x = EventHandlingFinished{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstate_eventstate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventHandlingFinished) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventHandlingFinished) ProtoMessage() {}

func (x *EventHandlingFinished) ProtoReflect() protoreflect.Message {
	mi := &file_eventstate_eventstate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventHandlingFinished.ProtoReflect.Descriptor instead.
func (*EventHandlingFinished) Descriptor() ([]byte, []int) {
	return file_eventstate_eventstate_proto_rawDescGZIP(), []int{2}
}

func (x *EventHandlingFinished) GetHandlerName() string {
	if x != nil {
		return x.HandlerName
	}
	return ""
}

// EventHandlingFailed is an event message occurred on a failure when handling given event.
type EventHandlingFailed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HandlerName string `protobuf:"bytes,1,opt,name=handler_name,json=handlerName,proto3" json:"handler_name,omitempty"`
	Err         string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	ErrCode     int32  `protobuf:"varint,3,opt,name=err_code,json=errCode,proto3" json:"err_code,omitempty"`
}

func (x *EventHandlingFailed) Reset() {
	*x = EventHandlingFailed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstate_eventstate_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventHandlingFailed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventHandlingFailed) ProtoMessage() {}

func (x *EventHandlingFailed) ProtoReflect() protoreflect.Message {
	mi := &file_eventstate_eventstate_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventHandlingFailed.ProtoReflect.Descriptor instead.
func (*EventHandlingFailed) Descriptor() ([]byte, []int) {
	return file_eventstate_eventstate_proto_rawDescGZIP(), []int{3}
}

func (x *EventHandlingFailed) GetHandlerName() string {
	if x != nil {
		return x.HandlerName
	}
	return ""
}

func (x *EventHandlingFailed) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *EventHandlingFailed) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

// FailureCountReset resets failure count for given event.
type FailureCountReset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HandlerName string `protobuf:"bytes,1,opt,name=handlerName,proto3" json:"handlerName,omitempty"`
}

func (x *FailureCountReset) Reset() {
	*x = FailureCountReset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventstate_eventstate_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailureCountReset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailureCountReset) ProtoMessage() {}

func (x *FailureCountReset) ProtoReflect() protoreflect.Message {
	mi := &file_eventstate_eventstate_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailureCountReset.ProtoReflect.Descriptor instead.
func (*FailureCountReset) Descriptor() ([]byte, []int) {
	return file_eventstate_eventstate_proto_rawDescGZIP(), []int{4}
}

func (x *FailureCountReset) GetHandlerName() string {
	if x != nil {
		return x.HandlerName
	}
	return ""
}

var File_eventstate_eventstate_proto protoreflect.FileDescriptor

var file_eventstate_eventstate_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xc8, 0x01, 0x0a, 0x0e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x55, 0x6e, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61, 0x78,
	0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x6d, 0x61, 0x78, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11,
	0x6d, 0x69, 0x6e, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x46, 0x61, 0x69, 0x6c,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x2a, 0x0a, 0x11, 0x6d, 0x61, 0x78, 0x5f,
	0x68, 0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0f, 0x6d, 0x61, 0x78, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x14, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x3a, 0x0a, 0x15, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67,
	0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x65, 0x0a, 0x13, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x69, 0x6e, 0x67, 0x46, 0x61, 0x69, 0x6c,
	0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x22, 0x35, 0x0a, 0x11, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x65, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x68, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x63, 0x6a, 0x61, 0x63, 0x2f, 0x63,
	0x6c, 0x65, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f,
	0x65, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eventstate_eventstate_proto_rawDescOnce sync.Once
	file_eventstate_eventstate_proto_rawDescData = file_eventstate_eventstate_proto_rawDesc
)

func file_eventstate_eventstate_proto_rawDescGZIP() []byte {
	file_eventstate_eventstate_proto_rawDescOnce.Do(func() {
		file_eventstate_eventstate_proto_rawDescData = protoimpl.X.CompressGZIP(file_eventstate_eventstate_proto_rawDescData)
	})
	return file_eventstate_eventstate_proto_rawDescData
}

var file_eventstate_eventstate_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_eventstate_eventstate_proto_goTypes = []interface{}{
	(*EventUnhandled)(nil),        // 0: eventstate.EventUnhandled
	(*EventHandlingStarted)(nil),  // 1: eventstate.EventHandlingStarted
	(*EventHandlingFinished)(nil), // 2: eventstate.EventHandlingFinished
	(*EventHandlingFailed)(nil),   // 3: eventstate.EventHandlingFailed
	(*FailureCountReset)(nil),     // 4: eventstate.FailureCountReset
}
var file_eventstate_eventstate_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_eventstate_eventstate_proto_init() }
func file_eventstate_eventstate_proto_init() {
	if File_eventstate_eventstate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eventstate_eventstate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventUnhandled); i {
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
		file_eventstate_eventstate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventHandlingStarted); i {
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
		file_eventstate_eventstate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventHandlingFinished); i {
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
		file_eventstate_eventstate_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventHandlingFailed); i {
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
		file_eventstate_eventstate_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailureCountReset); i {
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
			RawDescriptor: file_eventstate_eventstate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eventstate_eventstate_proto_goTypes,
		DependencyIndexes: file_eventstate_eventstate_proto_depIdxs,
		MessageInfos:      file_eventstate_eventstate_proto_msgTypes,
	}.Build()
	File_eventstate_eventstate_proto = out.File
	file_eventstate_eventstate_proto_rawDesc = nil
	file_eventstate_eventstate_proto_goTypes = nil
	file_eventstate_eventstate_proto_depIdxs = nil
}