// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: common.proto

package v1

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

// CommonClient is the client API for Common service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommonClient interface {
	CreateNotify(ctx context.Context, in *ReqCreateNotify, opts ...grpc.CallOption) (*RespCreateNotify, error)
	DeleteNotify(ctx context.Context, in *ReqDeleteNotify, opts ...grpc.CallOption) (*RespDeleteNotify, error)
	GetNotifyList(ctx context.Context, in *ReqGetNotifyList, opts ...grpc.CallOption) (*RespGetNotifyList, error)
	GetNotifyObject(ctx context.Context, in *ReqGetNotifyObject, opts ...grpc.CallOption) (*RespGetNotifyObject, error)
	CreateNotifyType(ctx context.Context, in *ReqCreateNotifyType, opts ...grpc.CallOption) (*RespCreateNotifyType, error)
	UpdateNotifyType(ctx context.Context, in *ReqUpdateNotifyType, opts ...grpc.CallOption) (*RespUpdateNotifyType, error)
	DeleteNotifyType(ctx context.Context, in *ReqDeleteNotifyType, opts ...grpc.CallOption) (*RespDeleteNotifyType, error)
	GetNotifyTypeList(ctx context.Context, in *ReqGetNotifyTypeList, opts ...grpc.CallOption) (*RespGetNotifyTypeList, error)
}

type commonClient struct {
	cc grpc.ClientConnInterface
}

func NewCommonClient(cc grpc.ClientConnInterface) CommonClient {
	return &commonClient{cc}
}

func (c *commonClient) CreateNotify(ctx context.Context, in *ReqCreateNotify, opts ...grpc.CallOption) (*RespCreateNotify, error) {
	out := new(RespCreateNotify)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/CreateNotify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) DeleteNotify(ctx context.Context, in *ReqDeleteNotify, opts ...grpc.CallOption) (*RespDeleteNotify, error) {
	out := new(RespDeleteNotify)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/DeleteNotify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) GetNotifyList(ctx context.Context, in *ReqGetNotifyList, opts ...grpc.CallOption) (*RespGetNotifyList, error) {
	out := new(RespGetNotifyList)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/GetNotifyList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) GetNotifyObject(ctx context.Context, in *ReqGetNotifyObject, opts ...grpc.CallOption) (*RespGetNotifyObject, error) {
	out := new(RespGetNotifyObject)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/GetNotifyObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) CreateNotifyType(ctx context.Context, in *ReqCreateNotifyType, opts ...grpc.CallOption) (*RespCreateNotifyType, error) {
	out := new(RespCreateNotifyType)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/CreateNotifyType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) UpdateNotifyType(ctx context.Context, in *ReqUpdateNotifyType, opts ...grpc.CallOption) (*RespUpdateNotifyType, error) {
	out := new(RespUpdateNotifyType)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/UpdateNotifyType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) DeleteNotifyType(ctx context.Context, in *ReqDeleteNotifyType, opts ...grpc.CallOption) (*RespDeleteNotifyType, error) {
	out := new(RespDeleteNotifyType)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/DeleteNotifyType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) GetNotifyTypeList(ctx context.Context, in *ReqGetNotifyTypeList, opts ...grpc.CallOption) (*RespGetNotifyTypeList, error) {
	out := new(RespGetNotifyTypeList)
	err := c.cc.Invoke(ctx, "/common.service.v1.Common/GetNotifyTypeList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommonServer is the server API for Common service.
// All implementations must embed UnimplementedCommonServer
// for forward compatibility
type CommonServer interface {
	CreateNotify(context.Context, *ReqCreateNotify) (*RespCreateNotify, error)
	DeleteNotify(context.Context, *ReqDeleteNotify) (*RespDeleteNotify, error)
	GetNotifyList(context.Context, *ReqGetNotifyList) (*RespGetNotifyList, error)
	GetNotifyObject(context.Context, *ReqGetNotifyObject) (*RespGetNotifyObject, error)
	CreateNotifyType(context.Context, *ReqCreateNotifyType) (*RespCreateNotifyType, error)
	UpdateNotifyType(context.Context, *ReqUpdateNotifyType) (*RespUpdateNotifyType, error)
	DeleteNotifyType(context.Context, *ReqDeleteNotifyType) (*RespDeleteNotifyType, error)
	GetNotifyTypeList(context.Context, *ReqGetNotifyTypeList) (*RespGetNotifyTypeList, error)
	mustEmbedUnimplementedCommonServer()
}

// UnimplementedCommonServer must be embedded to have forward compatible implementations.
type UnimplementedCommonServer struct {
}

func (UnimplementedCommonServer) CreateNotify(context.Context, *ReqCreateNotify) (*RespCreateNotify, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNotify not implemented")
}
func (UnimplementedCommonServer) DeleteNotify(context.Context, *ReqDeleteNotify) (*RespDeleteNotify, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNotify not implemented")
}
func (UnimplementedCommonServer) GetNotifyList(context.Context, *ReqGetNotifyList) (*RespGetNotifyList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotifyList not implemented")
}
func (UnimplementedCommonServer) GetNotifyObject(context.Context, *ReqGetNotifyObject) (*RespGetNotifyObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotifyObject not implemented")
}
func (UnimplementedCommonServer) CreateNotifyType(context.Context, *ReqCreateNotifyType) (*RespCreateNotifyType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNotifyType not implemented")
}
func (UnimplementedCommonServer) UpdateNotifyType(context.Context, *ReqUpdateNotifyType) (*RespUpdateNotifyType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNotifyType not implemented")
}
func (UnimplementedCommonServer) DeleteNotifyType(context.Context, *ReqDeleteNotifyType) (*RespDeleteNotifyType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNotifyType not implemented")
}
func (UnimplementedCommonServer) GetNotifyTypeList(context.Context, *ReqGetNotifyTypeList) (*RespGetNotifyTypeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotifyTypeList not implemented")
}
func (UnimplementedCommonServer) mustEmbedUnimplementedCommonServer() {}

// UnsafeCommonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommonServer will
// result in compilation errors.
type UnsafeCommonServer interface {
	mustEmbedUnimplementedCommonServer()
}

func RegisterCommonServer(s grpc.ServiceRegistrar, srv CommonServer) {
	s.RegisterService(&Common_ServiceDesc, srv)
}

func _Common_CreateNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqCreateNotify)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).CreateNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/CreateNotify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).CreateNotify(ctx, req.(*ReqCreateNotify))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_DeleteNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqDeleteNotify)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).DeleteNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/DeleteNotify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).DeleteNotify(ctx, req.(*ReqDeleteNotify))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_GetNotifyList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNotifyList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).GetNotifyList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/GetNotifyList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).GetNotifyList(ctx, req.(*ReqGetNotifyList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_GetNotifyObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNotifyObject)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).GetNotifyObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/GetNotifyObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).GetNotifyObject(ctx, req.(*ReqGetNotifyObject))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_CreateNotifyType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqCreateNotifyType)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).CreateNotifyType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/CreateNotifyType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).CreateNotifyType(ctx, req.(*ReqCreateNotifyType))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_UpdateNotifyType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqUpdateNotifyType)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).UpdateNotifyType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/UpdateNotifyType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).UpdateNotifyType(ctx, req.(*ReqUpdateNotifyType))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_DeleteNotifyType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqDeleteNotifyType)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).DeleteNotifyType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/DeleteNotifyType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).DeleteNotifyType(ctx, req.(*ReqDeleteNotifyType))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_GetNotifyTypeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNotifyTypeList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).GetNotifyTypeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.service.v1.Common/GetNotifyTypeList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).GetNotifyTypeList(ctx, req.(*ReqGetNotifyTypeList))
	}
	return interceptor(ctx, in, info, handler)
}

// Common_ServiceDesc is the grpc.ServiceDesc for Common service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Common_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "common.service.v1.Common",
	HandlerType: (*CommonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNotify",
			Handler:    _Common_CreateNotify_Handler,
		},
		{
			MethodName: "DeleteNotify",
			Handler:    _Common_DeleteNotify_Handler,
		},
		{
			MethodName: "GetNotifyList",
			Handler:    _Common_GetNotifyList_Handler,
		},
		{
			MethodName: "GetNotifyObject",
			Handler:    _Common_GetNotifyObject_Handler,
		},
		{
			MethodName: "CreateNotifyType",
			Handler:    _Common_CreateNotifyType_Handler,
		},
		{
			MethodName: "UpdateNotifyType",
			Handler:    _Common_UpdateNotifyType_Handler,
		},
		{
			MethodName: "DeleteNotifyType",
			Handler:    _Common_DeleteNotifyType_Handler,
		},
		{
			MethodName: "GetNotifyTypeList",
			Handler:    _Common_GetNotifyTypeList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common.proto",
}
