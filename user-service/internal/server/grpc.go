package server

import (
	"database/sql"

	"net"

	"github.com/nazibul7/go-grpc-microservices/user-service/internal/handler"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/service"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/store"
	pb "github.com/nazibul7/go-grpc-microservices/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(db *sql.DB) *grpc.Server {
	grpcServer := grpc.NewServer()

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore)
	pb.RegisterUserServiceServer(grpcServer, handler.NewUserHandler(userService))
	reflection.Register(grpcServer)

	return grpcServer
}

func Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
