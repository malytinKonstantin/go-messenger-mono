// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/messaging-service/messaging.proto

package messaging_service

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
	MessagingService_SendMessage_FullMethodName         = "/messaging.MessagingService/SendMessage"
	MessagingService_GetMessages_FullMethodName         = "/messaging.MessagingService/GetMessages"
	MessagingService_UpdateMessageStatus_FullMethodName = "/messaging.MessagingService/UpdateMessageStatus"
)

// MessagingServiceClient is the client API for MessagingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessagingServiceClient interface {
	// Отправка сообщения
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	// Получение истории переписки с пользователем
	GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error)
	// Обновление статуса сообщения
	UpdateMessageStatus(ctx context.Context, in *UpdateMessageStatusRequest, opts ...grpc.CallOption) (*UpdateMessageStatusResponse, error)
}

type messagingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessagingServiceClient(cc grpc.ClientConnInterface) MessagingServiceClient {
	return &messagingServiceClient{cc}
}

func (c *messagingServiceClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, MessagingService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagingServiceClient) GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMessagesResponse)
	err := c.cc.Invoke(ctx, MessagingService_GetMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagingServiceClient) UpdateMessageStatus(ctx context.Context, in *UpdateMessageStatusRequest, opts ...grpc.CallOption) (*UpdateMessageStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMessageStatusResponse)
	err := c.cc.Invoke(ctx, MessagingService_UpdateMessageStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessagingServiceServer is the server API for MessagingService service.
// All implementations must embed UnimplementedMessagingServiceServer
// for forward compatibility.
type MessagingServiceServer interface {
	// Отправка сообщения
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	// Получение истории переписки с пользователем
	GetMessages(context.Context, *GetMessagesRequest) (*GetMessagesResponse, error)
	// Обновление статуса сообщения
	UpdateMessageStatus(context.Context, *UpdateMessageStatusRequest) (*UpdateMessageStatusResponse, error)
	mustEmbedUnimplementedMessagingServiceServer()
}

// UnimplementedMessagingServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMessagingServiceServer struct{}

func (UnimplementedMessagingServiceServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessagingServiceServer) GetMessages(context.Context, *GetMessagesRequest) (*GetMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedMessagingServiceServer) UpdateMessageStatus(context.Context, *UpdateMessageStatusRequest) (*UpdateMessageStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMessageStatus not implemented")
}
func (UnimplementedMessagingServiceServer) mustEmbedUnimplementedMessagingServiceServer() {}
func (UnimplementedMessagingServiceServer) testEmbeddedByValue()                          {}

// UnsafeMessagingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessagingServiceServer will
// result in compilation errors.
type UnsafeMessagingServiceServer interface {
	mustEmbedUnimplementedMessagingServiceServer()
}

func RegisterMessagingServiceServer(s grpc.ServiceRegistrar, srv MessagingServiceServer) {
	// If the following call pancis, it indicates UnimplementedMessagingServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MessagingService_ServiceDesc, srv)
}

func _MessagingService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagingServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessagingService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagingServiceServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessagingService_GetMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagingServiceServer).GetMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessagingService_GetMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagingServiceServer).GetMessages(ctx, req.(*GetMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessagingService_UpdateMessageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMessageStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagingServiceServer).UpdateMessageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessagingService_UpdateMessageStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagingServiceServer).UpdateMessageStatus(ctx, req.(*UpdateMessageStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MessagingService_ServiceDesc is the grpc.ServiceDesc for MessagingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessagingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messaging.MessagingService",
	HandlerType: (*MessagingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _MessagingService_SendMessage_Handler,
		},
		{
			MethodName: "GetMessages",
			Handler:    _MessagingService_GetMessages_Handler,
		},
		{
			MethodName: "UpdateMessageStatus",
			Handler:    _MessagingService_UpdateMessageStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/messaging-service/messaging.proto",
}
