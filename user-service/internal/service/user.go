package service

import (
	"context"
	"errors"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/model"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/store"
)

type UserService struct {
	store *store.UserStore
}

func NewUserService(store *store.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.User, error) {
	if req.Name == "" {
		return model.User{}, errors.New("name is required")
	}

	if req.Email == "" {
		return model.User{}, errors.New("email is required")
	}

	return s.store.CreateUser(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req model.IDRequest) (model.User, error) {
	if req.ID <= 0 {
		return model.User{}, errors.New("invalid user id")
	}

	return s.store.GetUser(ctx, req.ID)
}

func (s *UserService) UpdateUser(ctx context.Context, req model.UpdateUserRequest) (model.User, error) {
	if req.ID <= 0 {
		return model.User{}, errors.New("invalid user id")
	}

	if req.Name == "" {
		return model.User{}, errors.New("name is required")
	}

	if req.Email == "" {
		return model.User{}, errors.New("email is required")
	}

	return s.store.UpdateUser(ctx, req)
}

func (s *UserService) DeleteUser(ctx context.Context, req model.IDRequest) (model.DeleteUserResponse, error) {
	if req.ID <= 0 {
		return model.DeleteUserResponse{},
			errors.New("invalid user id")
	}

	msg, err := s.store.DeleteUser(ctx, req.ID)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
