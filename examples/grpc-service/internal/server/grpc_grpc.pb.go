// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: grpc.proto

package server

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProfileControllerClient is the client API for ProfileController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileControllerClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Stream(ctx context.Context, opts ...grpc.CallOption) (ProfileController_StreamClient, error)
	Stream2(ctx context.Context, opts ...grpc.CallOption) (ProfileController_Stream2Client, error)
	Stream3(ctx context.Context, in *Stream3Request, opts ...grpc.CallOption) (ProfileController_Stream3Client, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type profileControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileControllerClient(cc grpc.ClientConnInterface) ProfileControllerClient {
	return &profileControllerClient{cc}
}

func (c *profileControllerClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/server.ProfileController/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileControllerClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/server.ProfileController/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileControllerClient) Stream(ctx context.Context, opts ...grpc.CallOption) (ProfileController_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProfileController_ServiceDesc.Streams[0], "/server.ProfileController/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &profileControllerStreamClient{stream}
	return x, nil
}

type ProfileController_StreamClient interface {
	Send(*Profile) error
	Recv() (*Statistic, error)
	grpc.ClientStream
}

type profileControllerStreamClient struct {
	grpc.ClientStream
}

func (x *profileControllerStreamClient) Send(m *Profile) error {
	return x.ClientStream.SendMsg(m)
}

func (x *profileControllerStreamClient) Recv() (*Statistic, error) {
	m := new(Statistic)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *profileControllerClient) Stream2(ctx context.Context, opts ...grpc.CallOption) (ProfileController_Stream2Client, error) {
	stream, err := c.cc.NewStream(ctx, &ProfileController_ServiceDesc.Streams[1], "/server.ProfileController/Stream2", opts...)
	if err != nil {
		return nil, err
	}
	x := &profileControllerStream2Client{stream}
	return x, nil
}

type ProfileController_Stream2Client interface {
	Send(*Profile) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type profileControllerStream2Client struct {
	grpc.ClientStream
}

func (x *profileControllerStream2Client) Send(m *Profile) error {
	return x.ClientStream.SendMsg(m)
}

func (x *profileControllerStream2Client) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *profileControllerClient) Stream3(ctx context.Context, in *Stream3Request, opts ...grpc.CallOption) (ProfileController_Stream3Client, error) {
	stream, err := c.cc.NewStream(ctx, &ProfileController_ServiceDesc.Streams[2], "/server.ProfileController/Stream3", opts...)
	if err != nil {
		return nil, err
	}
	x := &profileControllerStream3Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ProfileController_Stream3Client interface {
	Recv() (*Statistic, error)
	grpc.ClientStream
}

type profileControllerStream3Client struct {
	grpc.ClientStream
}

func (x *profileControllerStream3Client) Recv() (*Statistic, error) {
	m := new(Statistic)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *profileControllerClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/server.ProfileController/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileControllerServer is the server API for ProfileController service.
// All implementations must embed UnimplementedProfileControllerServer
// for forward compatibility
type ProfileControllerServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Remove(context.Context, *RemoveRequest) (*emptypb.Empty, error)
	Stream(ProfileController_StreamServer) error
	Stream2(ProfileController_Stream2Server) error
	Stream3(*Stream3Request, ProfileController_Stream3Server) error
	Update(context.Context, *UpdateRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedProfileControllerServer()
}

// UnimplementedProfileControllerServer must be embedded to have forward compatible implementations.
type UnimplementedProfileControllerServer struct {
}

func (UnimplementedProfileControllerServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedProfileControllerServer) Remove(context.Context, *RemoveRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedProfileControllerServer) Stream(ProfileController_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedProfileControllerServer) Stream2(ProfileController_Stream2Server) error {
	return status.Errorf(codes.Unimplemented, "method Stream2 not implemented")
}
func (UnimplementedProfileControllerServer) Stream3(*Stream3Request, ProfileController_Stream3Server) error {
	return status.Errorf(codes.Unimplemented, "method Stream3 not implemented")
}
func (UnimplementedProfileControllerServer) Update(context.Context, *UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedProfileControllerServer) mustEmbedUnimplementedProfileControllerServer() {}

// UnsafeProfileControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileControllerServer will
// result in compilation errors.
type UnsafeProfileControllerServer interface {
	mustEmbedUnimplementedProfileControllerServer()
}

func RegisterProfileControllerServer(s grpc.ServiceRegistrar, srv ProfileControllerServer) {
	s.RegisterService(&ProfileController_ServiceDesc, srv)
}

func _ProfileController_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileControllerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.ProfileController/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileControllerServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileController_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileControllerServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.ProfileController/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileControllerServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileController_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProfileControllerServer).Stream(&profileControllerStreamServer{stream})
}

type ProfileController_StreamServer interface {
	Send(*Statistic) error
	Recv() (*Profile, error)
	grpc.ServerStream
}

type profileControllerStreamServer struct {
	grpc.ServerStream
}

func (x *profileControllerStreamServer) Send(m *Statistic) error {
	return x.ServerStream.SendMsg(m)
}

func (x *profileControllerStreamServer) Recv() (*Profile, error) {
	m := new(Profile)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProfileController_Stream2_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProfileControllerServer).Stream2(&profileControllerStream2Server{stream})
}

type ProfileController_Stream2Server interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*Profile, error)
	grpc.ServerStream
}

type profileControllerStream2Server struct {
	grpc.ServerStream
}

func (x *profileControllerStream2Server) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *profileControllerStream2Server) Recv() (*Profile, error) {
	m := new(Profile)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProfileController_Stream3_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Stream3Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProfileControllerServer).Stream3(m, &profileControllerStream3Server{stream})
}

type ProfileController_Stream3Server interface {
	Send(*Statistic) error
	grpc.ServerStream
}

type profileControllerStream3Server struct {
	grpc.ServerStream
}

func (x *profileControllerStream3Server) Send(m *Statistic) error {
	return x.ServerStream.SendMsg(m)
}

func _ProfileController_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileControllerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.ProfileController/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileControllerServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileController_ServiceDesc is the grpc.ServiceDesc for ProfileController service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileController_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.ProfileController",
	HandlerType: (*ProfileControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ProfileController_Create_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _ProfileController_Remove_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ProfileController_Update_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _ProfileController_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Stream2",
			Handler:       _ProfileController_Stream2_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Stream3",
			Handler:       _ProfileController_Stream3_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc.proto",
}
