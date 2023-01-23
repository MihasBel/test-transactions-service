// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: delivery/grpc/proto/v1/transaction/api.proto

package transaction

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

// TransactionAPIClient is the client API for TransactionAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionAPIClient interface {
	// Getting transaction info
	ByID(ctx context.Context, in *ByIDRequest, opts ...grpc.CallOption) (*Transaction, error)
}

type transactionAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionAPIClient(cc grpc.ClientConnInterface) TransactionAPIClient {
	return &transactionAPIClient{cc}
}

func (c *transactionAPIClient) ByID(ctx context.Context, in *ByIDRequest, opts ...grpc.CallOption) (*Transaction, error) {
	out := new(Transaction)
	err := c.cc.Invoke(ctx, "/v1.transaction.TransactionAPI/ByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionAPIServer is the server API for TransactionAPI service.
// All implementations should embed UnimplementedTransactionAPIServer
// for forward compatibility
type TransactionAPIServer interface {
	// Getting transaction info
	ByID(context.Context, *ByIDRequest) (*Transaction, error)
}

// UnimplementedTransactionAPIServer should be embedded to have forward compatible implementations.
type UnimplementedTransactionAPIServer struct {
}

func (UnimplementedTransactionAPIServer) ByID(context.Context, *ByIDRequest) (*Transaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ByID not implemented")
}

// UnsafeTransactionAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionAPIServer will
// result in compilation errors.
type UnsafeTransactionAPIServer interface {
	mustEmbedUnimplementedTransactionAPIServer()
}

func RegisterTransactionAPIServer(s grpc.ServiceRegistrar, srv TransactionAPIServer) {
	s.RegisterService(&TransactionAPI_ServiceDesc, srv)
}

func _TransactionAPI_ByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionAPIServer).ByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.transaction.TransactionAPI/ByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionAPIServer).ByID(ctx, req.(*ByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionAPI_ServiceDesc is the grpc.ServiceDesc for TransactionAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.transaction.TransactionAPI",
	HandlerType: (*TransactionAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ByID",
			Handler:    _TransactionAPI_ByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "delivery/grpc/proto/v1/transaction/api.proto",
}
