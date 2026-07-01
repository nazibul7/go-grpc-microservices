package client

import (
	interceptor "github.com/nazibul7/go-grpc-microservices/api-gateway/internal/grpc"
	pb "github.com/nazibul7/go-grpc-microservices/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient() (*grpc.ClientConn, pb.UserServiceClient, error) {
	// Create a gRPC ClientConn to the User Service.
	//
	// A ClientConn is responsible for managing communication
	// with the gRPC server. Internally gRPC uses HTTP/2,
	// but the application interacts with the higher-level
	// gRPC abstraction rather than raw HTTP/2.

	// grpc.WithTransportCredentials(insecure.NewCredentials()) is required
	// when not using TLS. Without this option, gRPC defaults to expecting
	// a secure connection and will refuse to connect to a plaintext server.
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.AuthInterceptor),
	)
	if err != nil {
		return nil, nil, err
	}

	// It takes grpc client connectuon & returns all the methods which we created while protobuf file generate times.
	// So that we could call correct rpc methods.
	// The client and server can communicate because both are generated from the same .proto contract,
	// which defines the RPC method names, request types, response types, and serialization format.
	// The matching method names are part of the contract, but they are not the whole reason communication works.
	// The shared .proto definition is the real reason.
	client := pb.NewUserServiceClient(conn)
	return conn, client, nil
}
