// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: account_center.proto

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

// AccountCenterClient is the client API for AccountCenter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountCenterClient interface {
	SetAdmin(ctx context.Context, in *SetAdminRequest, opts ...grpc.CallOption) (*SetAdminReply, error)
	SendEmailCode(ctx context.Context, in *SendEmailCodeRequest, opts ...grpc.CallOption) (*SendEmailCodeReply, error)
	Login(ctx context.Context, in *CommonLoginRequest, opts ...grpc.CallOption) (*CommonLoginReply, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutReply, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error)
	WXLogin(ctx context.Context, in *WXLoginRequest, opts ...grpc.CallOption) (*WXLoginReply, error)
	GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*GetAccountInfoReply, error)
	PasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetReply, error)
	ListAccount(ctx context.Context, in *ListAccountRequest, opts ...grpc.CallOption) (*ListAccountReply, error)
	UpdateAccountInfo(ctx context.Context, in *UpdateAccountInfoRequest, opts ...grpc.CallOption) (*UpdateAccountInfoReply, error)
	GetPorn(ctx context.Context, in *GetPornRequest, opts ...grpc.CallOption) (*GetPornReply, error)
	GetGuest(ctx context.Context, in *GetGuestRequest, opts ...grpc.CallOption) (*GetGuestReply, error)
}

type accountCenterClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountCenterClient(cc grpc.ClientConnInterface) AccountCenterClient {
	return &accountCenterClient{cc}
}

func (c *accountCenterClient) SetAdmin(ctx context.Context, in *SetAdminRequest, opts ...grpc.CallOption) (*SetAdminReply, error) {
	out := new(SetAdminReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/SetAdmin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) SendEmailCode(ctx context.Context, in *SendEmailCodeRequest, opts ...grpc.CallOption) (*SendEmailCodeReply, error) {
	out := new(SendEmailCodeReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/SendEmailCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) Login(ctx context.Context, in *CommonLoginRequest, opts ...grpc.CallOption) (*CommonLoginReply, error) {
	out := new(CommonLoginReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutReply, error) {
	out := new(LogoutReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error) {
	out := new(RegisterReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) WXLogin(ctx context.Context, in *WXLoginRequest, opts ...grpc.CallOption) (*WXLoginReply, error) {
	out := new(WXLoginReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/WXLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*GetAccountInfoReply, error) {
	out := new(GetAccountInfoReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/GetAccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) PasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetReply, error) {
	out := new(PasswordResetReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/PasswordReset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) ListAccount(ctx context.Context, in *ListAccountRequest, opts ...grpc.CallOption) (*ListAccountReply, error) {
	out := new(ListAccountReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/ListAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) UpdateAccountInfo(ctx context.Context, in *UpdateAccountInfoRequest, opts ...grpc.CallOption) (*UpdateAccountInfoReply, error) {
	out := new(UpdateAccountInfoReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/UpdateAccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) GetPorn(ctx context.Context, in *GetPornRequest, opts ...grpc.CallOption) (*GetPornReply, error) {
	out := new(GetPornReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/GetPorn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountCenterClient) GetGuest(ctx context.Context, in *GetGuestRequest, opts ...grpc.CallOption) (*GetGuestReply, error) {
	out := new(GetGuestReply)
	err := c.cc.Invoke(ctx, "/ac.service.v1.AccountCenter/GetGuest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountCenterServer is the server API for AccountCenter service.
// All implementations must embed UnimplementedAccountCenterServer
// for forward compatibility
type AccountCenterServer interface {
	SetAdmin(context.Context, *SetAdminRequest) (*SetAdminReply, error)
	SendEmailCode(context.Context, *SendEmailCodeRequest) (*SendEmailCodeReply, error)
	Login(context.Context, *CommonLoginRequest) (*CommonLoginReply, error)
	Logout(context.Context, *LogoutRequest) (*LogoutReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	WXLogin(context.Context, *WXLoginRequest) (*WXLoginReply, error)
	GetAccountInfo(context.Context, *GetAccountInfoRequest) (*GetAccountInfoReply, error)
	PasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetReply, error)
	ListAccount(context.Context, *ListAccountRequest) (*ListAccountReply, error)
	UpdateAccountInfo(context.Context, *UpdateAccountInfoRequest) (*UpdateAccountInfoReply, error)
	GetPorn(context.Context, *GetPornRequest) (*GetPornReply, error)
	GetGuest(context.Context, *GetGuestRequest) (*GetGuestReply, error)
	mustEmbedUnimplementedAccountCenterServer()
}

// UnimplementedAccountCenterServer must be embedded to have forward compatible implementations.
type UnimplementedAccountCenterServer struct {
}

func (UnimplementedAccountCenterServer) SetAdmin(context.Context, *SetAdminRequest) (*SetAdminReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAdmin not implemented")
}
func (UnimplementedAccountCenterServer) SendEmailCode(context.Context, *SendEmailCodeRequest) (*SendEmailCodeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmailCode not implemented")
}
func (UnimplementedAccountCenterServer) Login(context.Context, *CommonLoginRequest) (*CommonLoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountCenterServer) Logout(context.Context, *LogoutRequest) (*LogoutReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAccountCenterServer) Register(context.Context, *RegisterRequest) (*RegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAccountCenterServer) WXLogin(context.Context, *WXLoginRequest) (*WXLoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WXLogin not implemented")
}
func (UnimplementedAccountCenterServer) GetAccountInfo(context.Context, *GetAccountInfoRequest) (*GetAccountInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountInfo not implemented")
}
func (UnimplementedAccountCenterServer) PasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PasswordReset not implemented")
}
func (UnimplementedAccountCenterServer) ListAccount(context.Context, *ListAccountRequest) (*ListAccountReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccount not implemented")
}
func (UnimplementedAccountCenterServer) UpdateAccountInfo(context.Context, *UpdateAccountInfoRequest) (*UpdateAccountInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccountInfo not implemented")
}
func (UnimplementedAccountCenterServer) GetPorn(context.Context, *GetPornRequest) (*GetPornReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPorn not implemented")
}
func (UnimplementedAccountCenterServer) GetGuest(context.Context, *GetGuestRequest) (*GetGuestReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGuest not implemented")
}
func (UnimplementedAccountCenterServer) mustEmbedUnimplementedAccountCenterServer() {}

// UnsafeAccountCenterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountCenterServer will
// result in compilation errors.
type UnsafeAccountCenterServer interface {
	mustEmbedUnimplementedAccountCenterServer()
}

func RegisterAccountCenterServer(s grpc.ServiceRegistrar, srv AccountCenterServer) {
	s.RegisterService(&AccountCenter_ServiceDesc, srv)
}

func _AccountCenter_SetAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).SetAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/SetAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).SetAdmin(ctx, req.(*SetAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_SendEmailCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEmailCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).SendEmailCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/SendEmailCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).SendEmailCode(ctx, req.(*SendEmailCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).Login(ctx, req.(*CommonLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_WXLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WXLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).WXLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/WXLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).WXLogin(ctx, req.(*WXLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_GetAccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).GetAccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/GetAccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).GetAccountInfo(ctx, req.(*GetAccountInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_PasswordReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).PasswordReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/PasswordReset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).PasswordReset(ctx, req.(*PasswordResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_ListAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).ListAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/ListAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).ListAccount(ctx, req.(*ListAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_UpdateAccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).UpdateAccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/UpdateAccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).UpdateAccountInfo(ctx, req.(*UpdateAccountInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_GetPorn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPornRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).GetPorn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/GetPorn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).GetPorn(ctx, req.(*GetPornRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountCenter_GetGuest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGuestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountCenterServer).GetGuest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ac.service.v1.AccountCenter/GetGuest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountCenterServer).GetGuest(ctx, req.(*GetGuestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountCenter_ServiceDesc is the grpc.ServiceDesc for AccountCenter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountCenter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ac.service.v1.AccountCenter",
	HandlerType: (*AccountCenterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetAdmin",
			Handler:    _AccountCenter_SetAdmin_Handler,
		},
		{
			MethodName: "SendEmailCode",
			Handler:    _AccountCenter_SendEmailCode_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AccountCenter_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AccountCenter_Logout_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _AccountCenter_Register_Handler,
		},
		{
			MethodName: "WXLogin",
			Handler:    _AccountCenter_WXLogin_Handler,
		},
		{
			MethodName: "GetAccountInfo",
			Handler:    _AccountCenter_GetAccountInfo_Handler,
		},
		{
			MethodName: "PasswordReset",
			Handler:    _AccountCenter_PasswordReset_Handler,
		},
		{
			MethodName: "ListAccount",
			Handler:    _AccountCenter_ListAccount_Handler,
		},
		{
			MethodName: "UpdateAccountInfo",
			Handler:    _AccountCenter_UpdateAccountInfo_Handler,
		},
		{
			MethodName: "GetPorn",
			Handler:    _AccountCenter_GetPorn_Handler,
		},
		{
			MethodName: "GetGuest",
			Handler:    _AccountCenter_GetGuest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account_center.proto",
}
