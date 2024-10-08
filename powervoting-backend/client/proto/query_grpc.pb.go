// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.7
// source: query.proto

package snapshot

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Snapshot_GetAddressPower_FullMethodName  = "/snapshot.Snapshot/GetAddressPower"
	Snapshot_SyncDateHeight_FullMethodName   = "/snapshot.Snapshot/SyncDateHeight"
	Snapshot_SyncAddrPower_FullMethodName    = "/snapshot.Snapshot/SyncAddrPower"
	Snapshot_SyncAllAddrPower_FullMethodName = "/snapshot.Snapshot/SyncAllAddrPower"
)

// SnapshotClient is the client API for Snapshot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SnapshotClient interface {
	GetAddressPower(ctx context.Context, in *AddressPowerRequest, opts ...grpc.CallOption) (*AddressPowerResponse, error)
	SyncDateHeight(ctx context.Context, in *SyncDateHeightRequest, opts ...grpc.CallOption) (*SyncDateHeightResponse, error)
	SyncAddrPower(ctx context.Context, in *SyncAddrPowerRequest, opts ...grpc.CallOption) (*SyncAddrPowerResponse, error)
	SyncAllAddrPower(ctx context.Context, in *SyncAllAddrPowerRequest, opts ...grpc.CallOption) (*SyncAllAddrPowerResponse, error)
}

type snapshotClient struct {
	cc grpc.ClientConnInterface
}

func NewSnapshotClient(cc grpc.ClientConnInterface) SnapshotClient {
	return &snapshotClient{cc}
}

func (c *snapshotClient) GetAddressPower(ctx context.Context, in *AddressPowerRequest, opts ...grpc.CallOption) (*AddressPowerResponse, error) {
	out := new(AddressPowerResponse)
	err := c.cc.Invoke(ctx, Snapshot_GetAddressPower_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snapshotClient) SyncDateHeight(ctx context.Context, in *SyncDateHeightRequest, opts ...grpc.CallOption) (*SyncDateHeightResponse, error) {
	out := new(SyncDateHeightResponse)
	err := c.cc.Invoke(ctx, Snapshot_SyncDateHeight_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snapshotClient) SyncAddrPower(ctx context.Context, in *SyncAddrPowerRequest, opts ...grpc.CallOption) (*SyncAddrPowerResponse, error) {
	out := new(SyncAddrPowerResponse)
	err := c.cc.Invoke(ctx, Snapshot_SyncAddrPower_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snapshotClient) SyncAllAddrPower(ctx context.Context, in *SyncAllAddrPowerRequest, opts ...grpc.CallOption) (*SyncAllAddrPowerResponse, error) {
	out := new(SyncAllAddrPowerResponse)
	err := c.cc.Invoke(ctx, Snapshot_SyncAllAddrPower_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SnapshotServer is the server API for Snapshot service.
// All implementations must embed UnimplementedSnapshotServer
// for forward compatibility
type SnapshotServer interface {
	GetAddressPower(context.Context, *AddressPowerRequest) (*AddressPowerResponse, error)
	SyncDateHeight(context.Context, *SyncDateHeightRequest) (*SyncDateHeightResponse, error)
	SyncAddrPower(context.Context, *SyncAddrPowerRequest) (*SyncAddrPowerResponse, error)
	SyncAllAddrPower(context.Context, *SyncAllAddrPowerRequest) (*SyncAllAddrPowerResponse, error)
	mustEmbedUnimplementedSnapshotServer()
}

// UnimplementedSnapshotServer must be embedded to have forward compatible implementations.
type UnimplementedSnapshotServer struct {
}

func (UnimplementedSnapshotServer) GetAddressPower(context.Context, *AddressPowerRequest) (*AddressPowerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressPower not implemented")
}
func (UnimplementedSnapshotServer) SyncDateHeight(context.Context, *SyncDateHeightRequest) (*SyncDateHeightResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncDateHeight not implemented")
}
func (UnimplementedSnapshotServer) SyncAddrPower(context.Context, *SyncAddrPowerRequest) (*SyncAddrPowerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncAddrPower not implemented")
}
func (UnimplementedSnapshotServer) SyncAllAddrPower(context.Context, *SyncAllAddrPowerRequest) (*SyncAllAddrPowerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncAllAddrPower not implemented")
}
func (UnimplementedSnapshotServer) mustEmbedUnimplementedSnapshotServer() {}

// UnsafeSnapshotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SnapshotServer will
// result in compilation errors.
type UnsafeSnapshotServer interface {
	mustEmbedUnimplementedSnapshotServer()
}

func RegisterSnapshotServer(s grpc.ServiceRegistrar, srv SnapshotServer) {
	s.RegisterService(&Snapshot_ServiceDesc, srv)
}

func _Snapshot_GetAddressPower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressPowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServer).GetAddressPower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Snapshot_GetAddressPower_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServer).GetAddressPower(ctx, req.(*AddressPowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Snapshot_SyncDateHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncDateHeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServer).SyncDateHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Snapshot_SyncDateHeight_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServer).SyncDateHeight(ctx, req.(*SyncDateHeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Snapshot_SyncAddrPower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncAddrPowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServer).SyncAddrPower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Snapshot_SyncAddrPower_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServer).SyncAddrPower(ctx, req.(*SyncAddrPowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Snapshot_SyncAllAddrPower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncAllAddrPowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnapshotServer).SyncAllAddrPower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Snapshot_SyncAllAddrPower_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnapshotServer).SyncAllAddrPower(ctx, req.(*SyncAllAddrPowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Snapshot_ServiceDesc is the grpc.ServiceDesc for Snapshot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Snapshot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "snapshot.Snapshot",
	HandlerType: (*SnapshotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAddressPower",
			Handler:    _Snapshot_GetAddressPower_Handler,
		},
		{
			MethodName: "SyncDateHeight",
			Handler:    _Snapshot_SyncDateHeight_Handler,
		},
		{
			MethodName: "SyncAddrPower",
			Handler:    _Snapshot_SyncAddrPower_Handler,
		},
		{
			MethodName: "SyncAllAddrPower",
			Handler:    _Snapshot_SyncAllAddrPower_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "query.proto",
}
