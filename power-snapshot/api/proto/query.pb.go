// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snapshot

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

type SyncAllAddrPowerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetId int64 `protobuf:"varint,1,opt,name=netId,proto3" json:"netId,omitempty"`
}

func (x *SyncAllAddrPowerRequest) Reset() {
	*x = SyncAllAddrPowerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllAddrPowerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllAddrPowerRequest) ProtoMessage() {}

func (x *SyncAllAddrPowerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllAddrPowerRequest.ProtoReflect.Descriptor instead.
func (*SyncAllAddrPowerRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{0}
}

func (x *SyncAllAddrPowerRequest) GetNetId() int64 {
	if x != nil {
		return x.NetId
	}
	return 0
}

type AddressPowerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetId     int64  `protobuf:"varint,1,opt,name=netId,proto3" json:"netId,omitempty"`
	Address   string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	RandomNum int32  `protobuf:"varint,3,opt,name=random_num,json=randomNum,proto3" json:"random_num,omitempty"`
}

func (x *AddressPowerRequest) Reset() {
	*x = AddressPowerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressPowerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressPowerRequest) ProtoMessage() {}

func (x *AddressPowerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressPowerRequest.ProtoReflect.Descriptor instead.
func (*AddressPowerRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{1}
}

func (x *AddressPowerRequest) GetNetId() int64 {
	if x != nil {
		return x.NetId
	}
	return 0
}

func (x *AddressPowerRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *AddressPowerRequest) GetRandomNum() int32 {
	if x != nil {
		return x.RandomNum
	}
	return 0
}

type SyncDateHeightRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetId int64 `protobuf:"varint,1,opt,name=netId,proto3" json:"netId,omitempty"`
}

func (x *SyncDateHeightRequest) Reset() {
	*x = SyncDateHeightRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncDateHeightRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDateHeightRequest) ProtoMessage() {}

func (x *SyncDateHeightRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDateHeightRequest.ProtoReflect.Descriptor instead.
func (*SyncDateHeightRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{2}
}

func (x *SyncDateHeightRequest) GetNetId() int64 {
	if x != nil {
		return x.NetId
	}
	return 0
}

type SyncAddrPowerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetId   int64  `protobuf:"varint,1,opt,name=netId,proto3" json:"netId,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *SyncAddrPowerRequest) Reset() {
	*x = SyncAddrPowerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAddrPowerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAddrPowerRequest) ProtoMessage() {}

func (x *SyncAddrPowerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAddrPowerRequest.ProtoReflect.Descriptor instead.
func (*SyncAddrPowerRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{3}
}

func (x *SyncAddrPowerRequest) GetNetId() int64 {
	if x != nil {
		return x.NetId
	}
	return 0
}

func (x *SyncAddrPowerRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type AddressPowerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address          string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	SpPower          string `protobuf:"bytes,2,opt,name=sp_power,json=spPower,proto3" json:"sp_power,omitempty"`
	ClientPower      string `protobuf:"bytes,3,opt,name=client_power,json=clientPower,proto3" json:"client_power,omitempty"`
	TokenHolderPower string `protobuf:"bytes,4,opt,name=token_holder_power,json=tokenHolderPower,proto3" json:"token_holder_power,omitempty"`
	DeveloperPower   string `protobuf:"bytes,5,opt,name=developer_power,json=developerPower,proto3" json:"developer_power,omitempty"`
	BlockHeight      int64  `protobuf:"varint,6,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	DateStr          string `protobuf:"bytes,7,opt,name=date_str,json=dateStr,proto3" json:"date_str,omitempty"`
}

func (x *AddressPowerResponse) Reset() {
	*x = AddressPowerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressPowerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressPowerResponse) ProtoMessage() {}

func (x *AddressPowerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressPowerResponse.ProtoReflect.Descriptor instead.
func (*AddressPowerResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{4}
}

func (x *AddressPowerResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *AddressPowerResponse) GetSpPower() string {
	if x != nil {
		return x.SpPower
	}
	return ""
}

func (x *AddressPowerResponse) GetClientPower() string {
	if x != nil {
		return x.ClientPower
	}
	return ""
}

func (x *AddressPowerResponse) GetTokenHolderPower() string {
	if x != nil {
		return x.TokenHolderPower
	}
	return ""
}

func (x *AddressPowerResponse) GetDeveloperPower() string {
	if x != nil {
		return x.DeveloperPower
	}
	return ""
}

func (x *AddressPowerResponse) GetBlockHeight() int64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *AddressPowerResponse) GetDateStr() string {
	if x != nil {
		return x.DateStr
	}
	return ""
}

type SyncAllAddrPowerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAllAddrPowerResponse) Reset() {
	*x = SyncAllAddrPowerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllAddrPowerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllAddrPowerResponse) ProtoMessage() {}

func (x *SyncAllAddrPowerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllAddrPowerResponse.ProtoReflect.Descriptor instead.
func (*SyncAllAddrPowerResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{5}
}

type SyncDateHeightResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncDateHeightResponse) Reset() {
	*x = SyncDateHeightResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncDateHeightResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDateHeightResponse) ProtoMessage() {}

func (x *SyncDateHeightResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDateHeightResponse.ProtoReflect.Descriptor instead.
func (*SyncDateHeightResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{6}
}

type SyncAddrPowerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAddrPowerResponse) Reset() {
	*x = SyncAddrPowerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAddrPowerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAddrPowerResponse) ProtoMessage() {}

func (x *SyncAddrPowerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAddrPowerResponse.ProtoReflect.Descriptor instead.
func (*SyncAddrPowerResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{7}
}

type SyncAllDeveloperWeightRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAllDeveloperWeightRequest) Reset() {
	*x = SyncAllDeveloperWeightRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllDeveloperWeightRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllDeveloperWeightRequest) ProtoMessage() {}

func (x *SyncAllDeveloperWeightRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllDeveloperWeightRequest.ProtoReflect.Descriptor instead.
func (*SyncAllDeveloperWeightRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{8}
}

type SyncAllDeveloperWeightResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncAllDeveloperWeightResponse) Reset() {
	*x = SyncAllDeveloperWeightResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllDeveloperWeightResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllDeveloperWeightResponse) ProtoMessage() {}

func (x *SyncAllDeveloperWeightResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllDeveloperWeightResponse.ProtoReflect.Descriptor instead.
func (*SyncAllDeveloperWeightResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{9}
}

type SyncDeveloperWeightRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DateStr string `protobuf:"bytes,1,opt,name=date_str,json=dateStr,proto3" json:"date_str,omitempty"`
}

func (x *SyncDeveloperWeightRequest) Reset() {
	*x = SyncDeveloperWeightRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncDeveloperWeightRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDeveloperWeightRequest) ProtoMessage() {}

func (x *SyncDeveloperWeightRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDeveloperWeightRequest.ProtoReflect.Descriptor instead.
func (*SyncDeveloperWeightRequest) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{10}
}

func (x *SyncDeveloperWeightRequest) GetDateStr() string {
	if x != nil {
		return x.DateStr
	}
	return ""
}

type SyncDeveloperWeightResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncDeveloperWeightResponse) Reset() {
	*x = SyncDeveloperWeightResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncDeveloperWeightResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncDeveloperWeightResponse) ProtoMessage() {}

func (x *SyncDeveloperWeightResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncDeveloperWeightResponse.ProtoReflect.Descriptor instead.
func (*SyncDeveloperWeightResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_proto_rawDescGZIP(), []int{11}
}

var File_proto_query_proto protoreflect.FileDescriptor

var file_proto_query_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x22, 0x2f, 0x0a,
	0x17, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x65, 0x74, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x22, 0x64,
	0x0a, 0x13, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x5f,
	0x6e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x72, 0x61, 0x6e, 0x64, 0x6f,
	0x6d, 0x4e, 0x75, 0x6d, 0x22, 0x2d, 0x0a, 0x15, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x65,
	0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x65,
	0x74, 0x49, 0x64, 0x22, 0x46, 0x0a, 0x14, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x64, 0x64, 0x72, 0x50,
	0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e,
	0x65, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x65, 0x74, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x83, 0x02, 0x0a, 0x14,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x19,
	0x0a, 0x08, 0x73, 0x70, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x70, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x12,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x77,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x48,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x65,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x50, 0x6f,
	0x77, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73,
	0x74, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74,
	0x72, 0x22, 0x1a, 0x0a, 0x18, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x41, 0x64, 0x64, 0x72,
	0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x0a,
	0x16, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x65, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x0a, 0x15, 0x53, 0x79, 0x6e, 0x63, 0x41,
	0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x1f, 0x0a, 0x1d, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x44, 0x65, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x20, 0x0a, 0x1e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x44, 0x65, 0x76, 0x65,
	0x6c, 0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x37, 0x0a, 0x1a, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x22, 0x1d, 0x0a, 0x1b,
	0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xbb, 0x04, 0x0a, 0x08,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x52, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x73, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f,
	0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x77,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0e,
	0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x65, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1f,
	0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61,
	0x74, 0x65, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x44,
	0x61, 0x74, 0x65, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x0d, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x64, 0x64, 0x72, 0x50,
	0x6f, 0x77, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e,
	0x53, 0x79, 0x6e, 0x63, 0x41, 0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e,
	0x53, 0x79, 0x6e, 0x63, 0x41, 0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5b, 0x0a, 0x10, 0x53, 0x79, 0x6e, 0x63, 0x41,
	0x6c, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x73, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x41, 0x64,
	0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c,
	0x6c, 0x41, 0x64, 0x64, 0x72, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x6d, 0x0a, 0x16, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x44,
	0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x27,
	0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c,
	0x6c, 0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68,
	0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f,
	0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x64, 0x0a, 0x13, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x24, 0x2e, 0x73, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f,
	0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x25, 0x2e, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x44, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_query_proto_rawDescOnce sync.Once
	file_proto_query_proto_rawDescData = file_proto_query_proto_rawDesc
)

func file_proto_query_proto_rawDescGZIP() []byte {
	file_proto_query_proto_rawDescOnce.Do(func() {
		file_proto_query_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_query_proto_rawDescData)
	})
	return file_proto_query_proto_rawDescData
}

var file_proto_query_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_query_proto_goTypes = []any{
	(*SyncAllAddrPowerRequest)(nil),        // 0: snapshot.SyncAllAddrPowerRequest
	(*AddressPowerRequest)(nil),            // 1: snapshot.AddressPowerRequest
	(*SyncDateHeightRequest)(nil),          // 2: snapshot.SyncDateHeightRequest
	(*SyncAddrPowerRequest)(nil),           // 3: snapshot.SyncAddrPowerRequest
	(*AddressPowerResponse)(nil),           // 4: snapshot.AddressPowerResponse
	(*SyncAllAddrPowerResponse)(nil),       // 5: snapshot.SyncAllAddrPowerResponse
	(*SyncDateHeightResponse)(nil),         // 6: snapshot.SyncDateHeightResponse
	(*SyncAddrPowerResponse)(nil),          // 7: snapshot.SyncAddrPowerResponse
	(*SyncAllDeveloperWeightRequest)(nil),  // 8: snapshot.SyncAllDeveloperWeightRequest
	(*SyncAllDeveloperWeightResponse)(nil), // 9: snapshot.SyncAllDeveloperWeightResponse
	(*SyncDeveloperWeightRequest)(nil),     // 10: snapshot.SyncDeveloperWeightRequest
	(*SyncDeveloperWeightResponse)(nil),    // 11: snapshot.SyncDeveloperWeightResponse
}
var file_proto_query_proto_depIdxs = []int32{
	1,  // 0: snapshot.Snapshot.GetAddressPower:input_type -> snapshot.AddressPowerRequest
	2,  // 1: snapshot.Snapshot.SyncDateHeight:input_type -> snapshot.SyncDateHeightRequest
	3,  // 2: snapshot.Snapshot.SyncAddrPower:input_type -> snapshot.SyncAddrPowerRequest
	0,  // 3: snapshot.Snapshot.SyncAllAddrPower:input_type -> snapshot.SyncAllAddrPowerRequest
	8,  // 4: snapshot.Snapshot.SyncAllDeveloperWeight:input_type -> snapshot.SyncAllDeveloperWeightRequest
	10, // 5: snapshot.Snapshot.SyncDeveloperWeight:input_type -> snapshot.SyncDeveloperWeightRequest
	4,  // 6: snapshot.Snapshot.GetAddressPower:output_type -> snapshot.AddressPowerResponse
	6,  // 7: snapshot.Snapshot.SyncDateHeight:output_type -> snapshot.SyncDateHeightResponse
	7,  // 8: snapshot.Snapshot.SyncAddrPower:output_type -> snapshot.SyncAddrPowerResponse
	5,  // 9: snapshot.Snapshot.SyncAllAddrPower:output_type -> snapshot.SyncAllAddrPowerResponse
	9,  // 10: snapshot.Snapshot.SyncAllDeveloperWeight:output_type -> snapshot.SyncAllDeveloperWeightResponse
	11, // 11: snapshot.Snapshot.SyncDeveloperWeight:output_type -> snapshot.SyncDeveloperWeightResponse
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_proto_query_proto_init() }
func file_proto_query_proto_init() {
	if File_proto_query_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_query_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAllAddrPowerRequest); i {
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
		file_proto_query_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AddressPowerRequest); i {
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
		file_proto_query_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SyncDateHeightRequest); i {
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
		file_proto_query_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAddrPowerRequest); i {
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
		file_proto_query_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*AddressPowerResponse); i {
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
		file_proto_query_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAllAddrPowerResponse); i {
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
		file_proto_query_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SyncDateHeightResponse); i {
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
		file_proto_query_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAddrPowerResponse); i {
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
		file_proto_query_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAllDeveloperWeightRequest); i {
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
		file_proto_query_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*SyncAllDeveloperWeightResponse); i {
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
		file_proto_query_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*SyncDeveloperWeightRequest); i {
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
		file_proto_query_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*SyncDeveloperWeightResponse); i {
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
			RawDescriptor: file_proto_query_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_query_proto_goTypes,
		DependencyIndexes: file_proto_query_proto_depIdxs,
		MessageInfos:      file_proto_query_proto_msgTypes,
	}.Build()
	File_proto_query_proto = out.File
	file_proto_query_proto_rawDesc = nil
	file_proto_query_proto_goTypes = nil
	file_proto_query_proto_depIdxs = nil
}
