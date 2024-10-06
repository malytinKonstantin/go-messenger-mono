// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: friendship-service/friendship.proto

package friendship_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FriendshipService_SendFriendRequest_FullMethodName   = "/friendship.FriendshipService/SendFriendRequest"
	FriendshipService_AcceptFriendRequest_FullMethodName = "/friendship.FriendshipService/AcceptFriendRequest"
)

// FriendshipServiceClient is the client API for FriendshipService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FriendshipServiceClient interface {
	SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*AcceptFriendRequestResponse, error)
}

type friendshipServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFriendshipServiceClient(cc grpc.ClientConnInterface) FriendshipServiceClient {
	return &friendshipServiceClient{cc}
}

func (c *friendshipServiceClient) SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*SendFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendshipService_SendFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendshipServiceClient) AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*AcceptFriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptFriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendshipService_AcceptFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FriendshipServiceServer is the server API for FriendshipService service.
// All implementations must embed UnimplementedFriendshipServiceServer
// for forward compatibility.
type FriendshipServiceServer interface {
	SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error)
	AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*AcceptFriendRequestResponse, error)
	mustEmbedUnimplementedFriendshipServiceServer()
}

// UnimplementedFriendshipServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFriendshipServiceServer struct{}

func (UnimplementedFriendshipServiceServer) SendFriendRequest(context.Context, *SendFriendRequestRequest) (*SendFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequest not implemented")
}
func (UnimplementedFriendshipServiceServer) AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*AcceptFriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptFriendRequest not implemented")
}
func (UnimplementedFriendshipServiceServer) mustEmbedUnimplementedFriendshipServiceServer() {}
func (UnimplementedFriendshipServiceServer) testEmbeddedByValue()                           {}

// UnsafeFriendshipServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FriendshipServiceServer will
// result in compilation errors.
type UnsafeFriendshipServiceServer interface {
	mustEmbedUnimplementedFriendshipServiceServer()
}

func RegisterFriendshipServiceServer(s grpc.ServiceRegistrar, srv FriendshipServiceServer) {
	// If the following call pancis, it indicates UnimplementedFriendshipServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FriendshipService_ServiceDesc, srv)
}

func _FriendshipService_SendFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendshipServiceServer).SendFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendshipService_SendFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendshipServiceServer).SendFriendRequest(ctx, req.(*SendFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendshipService_AcceptFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendshipServiceServer).AcceptFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendshipService_AcceptFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendshipServiceServer).AcceptFriendRequest(ctx, req.(*AcceptFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FriendshipService_ServiceDesc is the grpc.ServiceDesc for FriendshipService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FriendshipService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "friendship.FriendshipService",
	HandlerType: (*FriendshipServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendFriendRequest",
			Handler:    _FriendshipService_SendFriendRequest_Handler,
		},
		{
			MethodName: "AcceptFriendRequest",
			Handler:    _FriendshipService_AcceptFriendRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "friendship-service/friendship.proto",
}
