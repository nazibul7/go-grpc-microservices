package handler

import (
	"context"

	"log"

	pb "github.com/nazibul7/go-grpc-microservices/proto/user"

	"github.com/nazibul7/go-grpc-microservices/user-service/internal/model"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.service.CreateUser(ctx, model.CreateUserRequest{Name: req.Name, Email: req.Email})
	if err != nil {
		log.Printf("returning grpc error: %v", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
func (s *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetUser(ctx, model.IDRequest{ID: req.Id})
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
func (s *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.UpdateUser(ctx, model.UpdateUserRequest{
		ID:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
func (s *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	result, err := s.service.DeleteUser(ctx, model.IDRequest{ID: req.Id})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{
		Message: result.Message,
	}, nil
}
