package handler

import (
	"context"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/dto"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/service"
	pb "github.com/nazibul7/go-grpc-microservices/proto/auth"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer

	authService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(
	ctx context.Context,
	req *pb.SignUpRequest,
) (*pb.AuthResponse, error) {
	dtoReq := dto.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := h.authService.SignUp(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	

	return &pb.AuthResponse{
		User: &pb.User{
			Id:    resp.UserID,
			Email: req.Email,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, err
}

func (h *AuthHandler) SignIn(
	ctx context.Context,
	req *pb.SignInRequest,
) (*pb.AuthResponse, error) {

	dtoReq := dto.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := h.authService.SignIn(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		User: &pb.User{
			Id:    resp.UserID,
			Email: req.Email,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, err
}

func (h *AuthHandler) RefreshToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (*pb.AuthResponse, error) {
	dtoReq := &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
	resp, err := h.authService.RefreshToken(ctx, *dtoReq)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		User: &pb.User{
			Id: resp.UserID,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, err
}

func (h *AuthHandler) SignOut(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (*pb.SignOutResponse, error) {
	dtoReq := dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
	resp, err := h.authService.SignOut(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &pb.SignOutResponse{
		Message: resp.Message,
	}, err
}
