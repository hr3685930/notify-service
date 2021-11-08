package errs

import (
	"google.golang.org/grpc/codes"
	"notify-service/pkg/rpc"
	"runtime/debug"
)

func RpcRequestTimeout(msg string) error {
	return rpc.Err(codes.Canceled, msg, debug.Stack())
}

func RpcInternalServer(msg string) error {
	return rpc.Err(codes.Unknown, msg, debug.Stack())
}

func RpcValidationFailed(msg string) error {
	return rpc.Err(codes.InvalidArgument, msg, debug.Stack())
}

func RpcBadRequest(msg string) error {
	return rpc.Err(codes.OutOfRange, msg, debug.Stack())
}

func RpcResourceNotFound(msg string) error {
	return rpc.Err(codes.NotFound, msg, debug.Stack())
}

func RpcConflict(msg string) error {
	return rpc.Err(codes.AlreadyExists, msg, debug.Stack())
}

func RpcAuthorizationFailed(msg string) error {
	return rpc.Err(codes.PermissionDenied, msg, debug.Stack())
}

func RpcAuthenticationFailed(msg string) error {
	return rpc.Err(codes.Unauthenticated, msg, debug.Stack())
}

func RpcTooManyRequests(msg string) error {
	return rpc.Err(codes.ResourceExhausted, msg, debug.Stack())
}

func RpcServiceUnavailable(msg string) error {
	return rpc.Err(codes.Unavailable, msg, debug.Stack())
}
