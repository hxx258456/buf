// Copyright 2020-2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: buf/alpha/breaking/v1/config.proto

package breakingv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Config represents the breaking change configuration for a module. The rule and category IDs are defined
// by the version and apply across the config. The version is independent of the version of
// the package. The package version refers to the config shape, the version encoded in the Config message
// indicates which rule and category IDs should be used.
//
// The rule and category IDs are not encoded as enums in this package because we may want to support custom rule
// and category IDs in the future. Callers will need to resolve the rule and category ID strings.
type Config struct {
	state                             protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Version                string                 `protobuf:"bytes,1,opt,name=version,proto3"`
	xxx_hidden_UseIds                 []string               `protobuf:"bytes,2,rep,name=use_ids,json=useIds,proto3"`
	xxx_hidden_ExceptIds              []string               `protobuf:"bytes,3,rep,name=except_ids,json=exceptIds,proto3"`
	xxx_hidden_IgnorePaths            []string               `protobuf:"bytes,4,rep,name=ignore_paths,json=ignorePaths,proto3"`
	xxx_hidden_IgnoreIdPaths          *[]*IDPaths            `protobuf:"bytes,5,rep,name=ignore_id_paths,json=ignoreIdPaths,proto3"`
	xxx_hidden_IgnoreUnstablePackages bool                   `protobuf:"varint,6,opt,name=ignore_unstable_packages,json=ignoreUnstablePackages,proto3"`
	unknownFields                     protoimpl.UnknownFields
	sizeCache                         protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_buf_alpha_breaking_v1_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_breaking_v1_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Config) GetVersion() string {
	if x != nil {
		return x.xxx_hidden_Version
	}
	return ""
}

func (x *Config) GetUseIds() []string {
	if x != nil {
		return x.xxx_hidden_UseIds
	}
	return nil
}

func (x *Config) GetExceptIds() []string {
	if x != nil {
		return x.xxx_hidden_ExceptIds
	}
	return nil
}

func (x *Config) GetIgnorePaths() []string {
	if x != nil {
		return x.xxx_hidden_IgnorePaths
	}
	return nil
}

func (x *Config) GetIgnoreIdPaths() []*IDPaths {
	if x != nil {
		if x.xxx_hidden_IgnoreIdPaths != nil {
			return *x.xxx_hidden_IgnoreIdPaths
		}
	}
	return nil
}

func (x *Config) GetIgnoreUnstablePackages() bool {
	if x != nil {
		return x.xxx_hidden_IgnoreUnstablePackages
	}
	return false
}

func (x *Config) SetVersion(v string) {
	x.xxx_hidden_Version = v
}

func (x *Config) SetUseIds(v []string) {
	x.xxx_hidden_UseIds = v
}

func (x *Config) SetExceptIds(v []string) {
	x.xxx_hidden_ExceptIds = v
}

func (x *Config) SetIgnorePaths(v []string) {
	x.xxx_hidden_IgnorePaths = v
}

func (x *Config) SetIgnoreIdPaths(v []*IDPaths) {
	x.xxx_hidden_IgnoreIdPaths = &v
}

func (x *Config) SetIgnoreUnstablePackages(v bool) {
	x.xxx_hidden_IgnoreUnstablePackages = v
}

type Config_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// version represents the version of the breaking change rule and category IDs that should be used with this config.
	Version string
	// use_ids lists the rule and/or category IDs that are included in the breaking change check.
	UseIds []string
	// except_ids lists the rule and/or category IDs that are excluded from the breaking change check.
	ExceptIds []string
	// ignore_paths lists the paths of directories and/or files that should be ignored by the breaking change check.
	// All paths are relative to the root of the module.
	IgnorePaths []string
	// ignore_id_paths is a map of rule and/or category IDs to directory and/or file paths to exclude from the
	// breaking change check. This corresponds with the ignore_only configuration key.
	IgnoreIdPaths []*IDPaths
	// ignore_unstable_packages ignores packages with a last component that is one of the unstable forms recognised
	// by the PACKAGE_VERSION_SUFFIX:
	//
	//	v\d+test.*
	//	v\d+(alpha|beta)\d+
	//	v\d+p\d+(alpha|beta)\d+
	IgnoreUnstablePackages bool
}

func (b0 Config_builder) Build() *Config {
	m0 := &Config{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Version = b.Version
	x.xxx_hidden_UseIds = b.UseIds
	x.xxx_hidden_ExceptIds = b.ExceptIds
	x.xxx_hidden_IgnorePaths = b.IgnorePaths
	x.xxx_hidden_IgnoreIdPaths = &b.IgnoreIdPaths
	x.xxx_hidden_IgnoreUnstablePackages = b.IgnoreUnstablePackages
	return m0
}

// IDPaths represents a rule or category ID and the file and/or directory paths that are ignored for the rule.
type IDPaths struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Id    string                 `protobuf:"bytes,1,opt,name=id,proto3"`
	xxx_hidden_Paths []string               `protobuf:"bytes,2,rep,name=paths,proto3"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *IDPaths) Reset() {
	*x = IDPaths{}
	mi := &file_buf_alpha_breaking_v1_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IDPaths) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDPaths) ProtoMessage() {}

func (x *IDPaths) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_breaking_v1_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *IDPaths) GetId() string {
	if x != nil {
		return x.xxx_hidden_Id
	}
	return ""
}

func (x *IDPaths) GetPaths() []string {
	if x != nil {
		return x.xxx_hidden_Paths
	}
	return nil
}

func (x *IDPaths) SetId(v string) {
	x.xxx_hidden_Id = v
}

func (x *IDPaths) SetPaths(v []string) {
	x.xxx_hidden_Paths = v
}

type IDPaths_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Id    string
	Paths []string
}

func (b0 IDPaths_builder) Build() *IDPaths {
	m0 := &IDPaths{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Id = b.Id
	x.xxx_hidden_Paths = b.Paths
	return m0
}

var File_buf_alpha_breaking_v1_config_proto protoreflect.FileDescriptor

const file_buf_alpha_breaking_v1_config_proto_rawDesc = "" +
	"\n" +
	"\"buf/alpha/breaking/v1/config.proto\x12\x15buf.alpha.breaking.v1\"\xff\x01\n" +
	"\x06Config\x12\x18\n" +
	"\aversion\x18\x01 \x01(\tR\aversion\x12\x17\n" +
	"\ause_ids\x18\x02 \x03(\tR\x06useIds\x12\x1d\n" +
	"\n" +
	"except_ids\x18\x03 \x03(\tR\texceptIds\x12!\n" +
	"\fignore_paths\x18\x04 \x03(\tR\vignorePaths\x12F\n" +
	"\x0fignore_id_paths\x18\x05 \x03(\v2\x1e.buf.alpha.breaking.v1.IDPathsR\rignoreIdPaths\x128\n" +
	"\x18ignore_unstable_packages\x18\x06 \x01(\bR\x16ignoreUnstablePackages\"/\n" +
	"\aIDPaths\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05paths\x18\x02 \x03(\tR\x05pathsB\xee\x01\n" +
	"\x19com.buf.alpha.breaking.v1B\vConfigProtoP\x01ZMgithub.com/bufbuild/buf/private/gen/proto/go/buf/alpha/breaking/v1;breakingv1\xa2\x02\x03BAB\xaa\x02\x15Buf.Alpha.Breaking.V1\xca\x02\x15Buf\\Alpha\\Breaking\\V1\xe2\x02!Buf\\Alpha\\Breaking\\V1\\GPBMetadata\xea\x02\x18Buf::Alpha::Breaking::V1b\x06proto3"

var file_buf_alpha_breaking_v1_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_alpha_breaking_v1_config_proto_goTypes = []any{
	(*Config)(nil),  // 0: buf.alpha.breaking.v1.Config
	(*IDPaths)(nil), // 1: buf.alpha.breaking.v1.IDPaths
}
var file_buf_alpha_breaking_v1_config_proto_depIdxs = []int32{
	1, // 0: buf.alpha.breaking.v1.Config.ignore_id_paths:type_name -> buf.alpha.breaking.v1.IDPaths
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_buf_alpha_breaking_v1_config_proto_init() }
func file_buf_alpha_breaking_v1_config_proto_init() {
	if File_buf_alpha_breaking_v1_config_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_buf_alpha_breaking_v1_config_proto_rawDesc), len(file_buf_alpha_breaking_v1_config_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_buf_alpha_breaking_v1_config_proto_goTypes,
		DependencyIndexes: file_buf_alpha_breaking_v1_config_proto_depIdxs,
		MessageInfos:      file_buf_alpha_breaking_v1_config_proto_msgTypes,
	}.Build()
	File_buf_alpha_breaking_v1_config_proto = out.File
	file_buf_alpha_breaking_v1_config_proto_goTypes = nil
	file_buf_alpha_breaking_v1_config_proto_depIdxs = nil
}
