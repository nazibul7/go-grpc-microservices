package server

import (
	"database/sql"
	"net"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/handler"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/service"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/store"
	pb "github.com/nazibul7/go-grpc-microservices/proto/auth"
	userpb "github.com/nazibul7/go-grpc-microservices/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCServer(db *sql.DB) *grpc.Server {
	grpcServer := grpc.NewServer()

	userStore := store.NewUserStore()
	refreshTokenStore := store.NewRefreshTokenStore()

	conn, err := grpc.NewClient(":50051", 
	grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	userClient := userpb.NewUserServiceClient(conn)

	authService := service.NewAuthService(
		db,
		userStore,
		refreshTokenStore,
		userClient,
	)

	authHandler := handler.NewAuthHandler(
		authService,
	)

	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	return grpcServer
}

func Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
