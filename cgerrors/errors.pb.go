// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: errors.proto

package cgerrors

import (
	"reflect"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ErrorCode is a code that defines errors specification.
type ErrorCode int32

const (
	ErrorCode_OK                 ErrorCode = 0
	ErrorCode_Canceled           ErrorCode = 1
	ErrorCode_Unknown            ErrorCode = 2
	ErrorCode_InvalidArgument    ErrorCode = 3
	ErrorCode_DeadlineExceeded   ErrorCode = 4
	ErrorCode_NotFound           ErrorCode = 5
	ErrorCode_AlreadyExists      ErrorCode = 6
	ErrorCode_PermissionDenied   ErrorCode = 7
	ErrorCode_ResourceExhausted  ErrorCode = 8
	ErrorCode_FailedPrecondition ErrorCode = 9
	ErrorCode_Aborted            ErrorCode = 10
	ErrorCode_OutOfRange         ErrorCode = 11
	ErrorCode_Unimplemented      ErrorCode = 12
	ErrorCode_Internal           ErrorCode = 13
	ErrorCode_Unavailable        ErrorCode = 14
	ErrorCode_DataLoss           ErrorCode = 15
	ErrorCode_Unauthenticated    ErrorCode = 16
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0:  "OK",
		1:  "Canceled",
		2:  "Unknown",
		3:  "InvalidArgument",
		4:  "DeadlineExceeded",
		5:  "NotFound",
		6:  "AlreadyExists",
		7:  "PermissionDenied",
		8:  "ResourceExhausted",
		9:  "FailedPrecondition",
		10: "Aborted",
		11: "OutOfRange",
		12: "Unimplemented",
		13: "Internal",
		14: "Unavailable",
		15: "DataLoss",
		16: "Unauthenticated",
	}
	ErrorCode_value = map[string]int32{
		"OK":                 0,
		"Canceled":           1,
		"Unknown":            2,
		"InvalidArgument":    3,
		"DeadlineExceeded":   4,
		"NotFound":           5,
		"AlreadyExists":      6,
		"PermissionDenied":   7,
		"ResourceExhausted":  8,
		"FailedPrecondition": 9,
		"Aborted":            10,
		"OutOfRange":         11,
		"Unimplemented":      12,
		"Internal":           13,
		"Unavailable":        14,
		"DataLoss":           15,
		"Unauthenticated":    16,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_errors_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_errors_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_errors_proto_rawDescGZIP(), []int{0}
}

// Error is the error message that has id, it's code and a detail.
type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code    ErrorCode         `protobuf:"varint,2,opt,name=code,proto3,enum=cgerrors.ErrorCode" json:"code,omitempty"`
	Detail  string            `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	Process string            `protobuf:"bytes,4,opt,name=process,proto3" json:"process,omitempty"`
	Meta    map[string]string `protobuf:"bytes,5,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_errors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_errors_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Error) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_OK
}

func (x *Error) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *Error) GetProcess() string {
	if x != nil {
		return x.Process
	}
	return ""
}

func (x *Error) GetMeta() map[string]string {
	if x != nil {
		return x.Meta
	}
	return nil
}

var File_errors_proto protoreflect.FileDescriptor

var file_errors_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x63, 0x67, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x22, 0xda, 0x01, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x27, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x13, 0x2e, 0x63, 0x67, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x12, 0x2d, 0x0a,
	0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x67,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09,
	0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0xb1, 0x02, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x43,
	0x61, 0x6e, 0x63, 0x65, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b,
	0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x44,
	0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x45, 0x78, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x10,
	0x04, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x05, 0x12,
	0x11, 0x0a, 0x0d, 0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73,
	0x10, 0x06, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x44, 0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x45, 0x78, 0x68, 0x61, 0x75, 0x73, 0x74, 0x65, 0x64, 0x10, 0x08, 0x12,
	0x16, 0x0a, 0x12, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x09, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x62, 0x6f, 0x72, 0x74,
	0x65, 0x64, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x75, 0x74, 0x4f, 0x66, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x10, 0x0b, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x6e, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x65, 0x64, 0x10, 0x0c, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x10, 0x0d, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x10, 0x0e, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x6f,
	0x73, 0x73, 0x10, 0x0f, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x6e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x10, 0x10, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x63, 0x6a, 0x61, 0x63, 0x2f, 0x63,
	0x6c, 0x65, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x63, 0x67, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errors_proto_rawDescOnce sync.Once
	file_errors_proto_rawDescData = file_errors_proto_rawDesc
)

func file_errors_proto_rawDescGZIP() []byte {
	file_errors_proto_rawDescOnce.Do(func() {
		file_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_errors_proto_rawDescData)
	})
	return file_errors_proto_rawDescData
}

var file_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errors_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_errors_proto_goTypes = []interface{}{
	(ErrorCode)(0), // 0: cgerrors.ErrorCode
	(*Error)(nil),  // 1: cgerrors.Error
	nil,            // 2: cgerrors.Error.MetaEntry
}
var file_errors_proto_depIdxs = []int32{
	0, // 0: cgerrors.Error.code:type_name -> cgerrors.ErrorCode
	2, // 1: cgerrors.Error.meta:type_name -> cgerrors.Error.MetaEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_errors_proto_init() }
func file_errors_proto_init() {
	if File_errors_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_errors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errors_proto_goTypes,
		DependencyIndexes: file_errors_proto_depIdxs,
		EnumInfos:         file_errors_proto_enumTypes,
		MessageInfos:      file_errors_proto_msgTypes,
	}.Build()
	File_errors_proto = out.File
	file_errors_proto_rawDesc = nil
	file_errors_proto_goTypes = nil
	file_errors_proto_depIdxs = nil
}