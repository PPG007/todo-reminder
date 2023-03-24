// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: service.proto

package proto

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

// ChatGPTServiceClient is the client API for ChatGPTService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatGPTServiceClient interface {
	GetTextResponse(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	GetImageResponse(ctx context.Context, in *String, opts ...grpc.CallOption) (*Image, error)
}

type chatGPTServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatGPTServiceClient(cc grpc.ClientConnInterface) ChatGPTServiceClient {
	return &chatGPTServiceClient{cc}
}

func (c *chatGPTServiceClient) GetTextResponse(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/ChatGPTService/GetTextResponse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatGPTServiceClient) GetImageResponse(ctx context.Context, in *String, opts ...grpc.CallOption) (*Image, error) {
	out := new(Image)
	err := c.cc.Invoke(ctx, "/ChatGPTService/GetImageResponse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatGPTServiceServer is the server API for ChatGPTService service.
// All implementations must embed UnimplementedChatGPTServiceServer
// for forward compatibility
type ChatGPTServiceServer interface {
	GetTextResponse(context.Context, *String) (*String, error)
	GetImageResponse(context.Context, *String) (*Image, error)
	mustEmbedUnimplementedChatGPTServiceServer()
}

// UnimplementedChatGPTServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatGPTServiceServer struct {
}

func (UnimplementedChatGPTServiceServer) GetTextResponse(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTextResponse not implemented")
}
func (UnimplementedChatGPTServiceServer) GetImageResponse(context.Context, *String) (*Image, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetImageResponse not implemented")
}
func (UnimplementedChatGPTServiceServer) mustEmbedUnimplementedChatGPTServiceServer() {}

// UnsafeChatGPTServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatGPTServiceServer will
// result in compilation errors.
type UnsafeChatGPTServiceServer interface {
	mustEmbedUnimplementedChatGPTServiceServer()
}

func RegisterChatGPTServiceServer(s grpc.ServiceRegistrar, srv ChatGPTServiceServer) {
	s.RegisterService(&ChatGPTService_ServiceDesc, srv)
}

func _ChatGPTService_GetTextResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatGPTServiceServer).GetTextResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ChatGPTService/GetTextResponse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatGPTServiceServer).GetTextResponse(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatGPTService_GetImageResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatGPTServiceServer).GetImageResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ChatGPTService/GetImageResponse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatGPTServiceServer).GetImageResponse(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatGPTService_ServiceDesc is the grpc.ServiceDesc for ChatGPTService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatGPTService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChatGPTService",
	HandlerType: (*ChatGPTServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTextResponse",
			Handler:    _ChatGPTService_GetTextResponse_Handler,
		},
		{
			MethodName: "GetImageResponse",
			Handler:    _ChatGPTService_GetImageResponse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}