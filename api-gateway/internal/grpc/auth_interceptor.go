package grpc

import (
	"context"

	"github.com/nazibul7/go-grpc-microservices/api-gateway/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is a CLIENT interceptor.
//
// It runs BEFORE every outgoing gRPC request made by the API Gateway.
//
// Flow:
//
// HTTP Client
//      │
// Authorization: Bearer <jwt>
//      │
// HTTP Middleware
//      │
// Store token in context
//      │
// HTTP Handler
//      │
// userClient.GetUser(...)
//      │
// AuthInterceptor  <----- HERE
//      │
// Attach JWT as gRPC metadata
//      │
// Auth/User Service
func AuthInterceptor(
	ctx context.Context,

	// Full gRPC method.
	//
	// Example:
	// /user.UserService/GetUser
	// /auth.AuthService/SignUp
	method string,

	// Incoming protobuf request.
	req any,

	// Empty protobuf response.
	//
	// gRPC fills this after the server replies.
	reply any,

	// Underlying gRPC connection.
	cc *grpc.ClientConn,

	// Function that performs the ACTUAL RPC.
	//
	// Similar to next.ServeHTTP() in HTTP middleware.
	invoker grpc.UnaryInvoker,

	// Extra gRPC call options.
	opts ...grpc.CallOption,
) error {

	// Read JWT stored by HTTP middleware.
	token, ok := ctx.Value(middleware.TokenKey).(string)

	if ok && token != "" {

		// Create gRPC metadata.
		//
		// Equivalent to HTTP:
		//
		// Authorization: Bearer <jwt>
		md := metadata.Pairs(
			"authorization",
			"Bearer "+token,
		)

		// Attach metadata to outgoing context.
		//
		// Only metadata, deadlines and cancellation
		// are transmitted over gRPC.
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Continue with the actual RPC.
	//
	// If this is NOT called,
	// the request never reaches the server.
	return invoker(
		ctx,
		method,
		req,
		reply,
		cc,
		opts...,
	)
}