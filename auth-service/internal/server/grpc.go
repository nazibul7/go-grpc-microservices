package server

import (
	"database/sql"
	"net"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/handler"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/service"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/store"
	pb "github.com/nazibul7/go-grpc-microservices/proto/auth"
	"google.golang.org/grpc"
)

func NewGRPCServer(db *sql.DB) *grpc.Server {
	grpcServer := grpc.NewServer()

	userStore := store.NewUserStore()
	refreshTokenStore := store.NewRefreshTokenStore()

	authService := service.NewAuthService(
		db,
		userStore,
		refreshTokenStore,
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
