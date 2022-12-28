// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/proto/gokeeper.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	SignIn(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*JWT, error)
	SignUp(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*JWT, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) SignIn(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*JWT, error) {
	out := new(JWT)
	err := c.cc.Invoke(ctx, "/gokeeper.Auth/SignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) SignUp(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*JWT, error) {
	out := new(JWT)
	err := c.cc.Invoke(ctx, "/gokeeper.Auth/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	SignIn(context.Context, *Credentials) (*JWT, error)
	SignUp(context.Context, *Credentials) (*JWT, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) SignIn(context.Context, *Credentials) (*JWT, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedAuthServer) SignUp(context.Context, *Credentials) (*JWT, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Auth/SignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).SignIn(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Auth/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).SignUp(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gokeeper.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignIn",
			Handler:    _Auth_SignIn_Handler,
		},
		{
			MethodName: "SignUp",
			Handler:    _Auth_SignUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/gokeeper.proto",
}

// KeeperClient is the client API for Keeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperClient interface {
	StoreSecret(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*ClientSecret, error)
	UpdateSecret(ctx context.Context, in *ClientSecret, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteSecret(ctx context.Context, in *ClientSecret, opts ...grpc.CallOption) (*empty.Empty, error)
	StoreMeta(ctx context.Context, in *StoreMetaRequest, opts ...grpc.CallOption) (*ClientMeta, error)
	UpdateMeta(ctx context.Context, in *ClientMeta, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteMeta(ctx context.Context, in *ClientMeta, opts ...grpc.CallOption) (*empty.Empty, error)
	GetSecrets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClientSecrets, error)
	GetEncryptedKey(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClientSecret, error)
}

type keeperClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperClient(cc grpc.ClientConnInterface) KeeperClient {
	return &keeperClient{cc}
}

func (c *keeperClient) StoreSecret(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*ClientSecret, error) {
	out := new(ClientSecret)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/StoreSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) UpdateSecret(ctx context.Context, in *ClientSecret, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/UpdateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) DeleteSecret(ctx context.Context, in *ClientSecret, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/DeleteSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) StoreMeta(ctx context.Context, in *StoreMetaRequest, opts ...grpc.CallOption) (*ClientMeta, error) {
	out := new(ClientMeta)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/StoreMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) UpdateMeta(ctx context.Context, in *ClientMeta, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/UpdateMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) DeleteMeta(ctx context.Context, in *ClientMeta, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/DeleteMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) GetSecrets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClientSecrets, error) {
	out := new(ClientSecrets)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/GetSecrets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) GetEncryptedKey(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClientSecret, error) {
	out := new(ClientSecret)
	err := c.cc.Invoke(ctx, "/gokeeper.Keeper/GetEncryptedKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServer is the server API for Keeper service.
// All implementations must embed UnimplementedKeeperServer
// for forward compatibility
type KeeperServer interface {
	StoreSecret(context.Context, *Secret) (*ClientSecret, error)
	UpdateSecret(context.Context, *ClientSecret) (*empty.Empty, error)
	DeleteSecret(context.Context, *ClientSecret) (*empty.Empty, error)
	StoreMeta(context.Context, *StoreMetaRequest) (*ClientMeta, error)
	UpdateMeta(context.Context, *ClientMeta) (*empty.Empty, error)
	DeleteMeta(context.Context, *ClientMeta) (*empty.Empty, error)
	GetSecrets(context.Context, *empty.Empty) (*ClientSecrets, error)
	GetEncryptedKey(context.Context, *empty.Empty) (*ClientSecret, error)
	mustEmbedUnimplementedKeeperServer()
}

// UnimplementedKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServer struct {
}

func (UnimplementedKeeperServer) StoreSecret(context.Context, *Secret) (*ClientSecret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreSecret not implemented")
}
func (UnimplementedKeeperServer) UpdateSecret(context.Context, *ClientSecret) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSecret not implemented")
}
func (UnimplementedKeeperServer) DeleteSecret(context.Context, *ClientSecret) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (UnimplementedKeeperServer) StoreMeta(context.Context, *StoreMetaRequest) (*ClientMeta, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreMeta not implemented")
}
func (UnimplementedKeeperServer) UpdateMeta(context.Context, *ClientMeta) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMeta not implemented")
}
func (UnimplementedKeeperServer) DeleteMeta(context.Context, *ClientMeta) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMeta not implemented")
}
func (UnimplementedKeeperServer) GetSecrets(context.Context, *empty.Empty) (*ClientSecrets, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecrets not implemented")
}
func (UnimplementedKeeperServer) GetEncryptedKey(context.Context, *empty.Empty) (*ClientSecret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEncryptedKey not implemented")
}
func (UnimplementedKeeperServer) mustEmbedUnimplementedKeeperServer() {}

// UnsafeKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperServer will
// result in compilation errors.
type UnsafeKeeperServer interface {
	mustEmbedUnimplementedKeeperServer()
}

func RegisterKeeperServer(s grpc.ServiceRegistrar, srv KeeperServer) {
	s.RegisterService(&Keeper_ServiceDesc, srv)
}

func _Keeper_StoreSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Secret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).StoreSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/StoreSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).StoreSecret(ctx, req.(*Secret))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_UpdateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientSecret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).UpdateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/UpdateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).UpdateSecret(ctx, req.(*ClientSecret))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_DeleteSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientSecret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).DeleteSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/DeleteSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).DeleteSecret(ctx, req.(*ClientSecret))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_StoreMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreMetaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).StoreMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/StoreMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).StoreMeta(ctx, req.(*StoreMetaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_UpdateMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientMeta)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).UpdateMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/UpdateMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).UpdateMeta(ctx, req.(*ClientMeta))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_DeleteMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientMeta)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).DeleteMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/DeleteMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).DeleteMeta(ctx, req.(*ClientMeta))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_GetSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/GetSecrets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetSecrets(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_GetEncryptedKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).GetEncryptedKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gokeeper.Keeper/GetEncryptedKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).GetEncryptedKey(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Keeper_ServiceDesc is the grpc.ServiceDesc for Keeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Keeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gokeeper.Keeper",
	HandlerType: (*KeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreSecret",
			Handler:    _Keeper_StoreSecret_Handler,
		},
		{
			MethodName: "UpdateSecret",
			Handler:    _Keeper_UpdateSecret_Handler,
		},
		{
			MethodName: "DeleteSecret",
			Handler:    _Keeper_DeleteSecret_Handler,
		},
		{
			MethodName: "StoreMeta",
			Handler:    _Keeper_StoreMeta_Handler,
		},
		{
			MethodName: "UpdateMeta",
			Handler:    _Keeper_UpdateMeta_Handler,
		},
		{
			MethodName: "DeleteMeta",
			Handler:    _Keeper_DeleteMeta_Handler,
		},
		{
			MethodName: "GetSecrets",
			Handler:    _Keeper_GetSecrets_Handler,
		},
		{
			MethodName: "GetEncryptedKey",
			Handler:    _Keeper_GetEncryptedKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/gokeeper.proto",
}
