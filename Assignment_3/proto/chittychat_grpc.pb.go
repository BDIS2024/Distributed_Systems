// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.1
// source: proto/chittychat.proto

package proto

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
	ChittyChatService_SendTimestamp_FullMethodName = "/ChittyChatService/SendTimestamp"
	ChittyChatService_SendMessage_FullMethodName   = "/ChittyChatService/SendMessage"
)

// ChittyChatServiceClient is the client API for ChittyChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatServiceClient interface {
	SendTimestamp(ctx context.Context, in *Timestamp, opts ...grpc.CallOption) (*Empty, error)
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Empty, error)
}

type chittyChatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatServiceClient(cc grpc.ClientConnInterface) ChittyChatServiceClient {
	return &chittyChatServiceClient{cc}
}

func (c *chittyChatServiceClient) SendTimestamp(ctx context.Context, in *Timestamp, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ChittyChatService_SendTimestamp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatServiceClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ChittyChatService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServiceServer is the server API for ChittyChatService service.
// All implementations must embed UnimplementedChittyChatServiceServer
// for forward compatibility.
type ChittyChatServiceServer interface {
	SendTimestamp(context.Context, *Timestamp) (*Empty, error)
	SendMessage(context.Context, *Message) (*Empty, error)
	mustEmbedUnimplementedChittyChatServiceServer()
}

// UnimplementedChittyChatServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChittyChatServiceServer struct{}

func (UnimplementedChittyChatServiceServer) SendTimestamp(context.Context, *Timestamp) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTimestamp not implemented")
}
func (UnimplementedChittyChatServiceServer) SendMessage(context.Context, *Message) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChittyChatServiceServer) mustEmbedUnimplementedChittyChatServiceServer() {}
func (UnimplementedChittyChatServiceServer) testEmbeddedByValue()                           {}

// UnsafeChittyChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServiceServer will
// result in compilation errors.
type UnsafeChittyChatServiceServer interface {
	mustEmbedUnimplementedChittyChatServiceServer()
}

func RegisterChittyChatServiceServer(s grpc.ServiceRegistrar, srv ChittyChatServiceServer) {
	// If the following call pancis, it indicates UnimplementedChittyChatServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChittyChatService_ServiceDesc, srv)
}

func _ChittyChatService_SendTimestamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Timestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).SendTimestamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChittyChatService_SendTimestamp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).SendTimestamp(ctx, req.(*Timestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChatService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChittyChatService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChatService_ServiceDesc is the grpc.ServiceDesc for ChittyChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChittyChatService",
	HandlerType: (*ChittyChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTimestamp",
			Handler:    _ChittyChatService_SendTimestamp_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _ChittyChatService_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/chittychat.proto",
}
