package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/dto"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/model"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/store"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/utils"
	"google.golang.org/grpc/metadata"

	userpb "github.com/nazibul7/go-grpc-microservices/proto/user"
)

type AuthService struct {
	db                *sql.DB
	credentialStore   *store.CredentialStore
	refreshTokenStore *store.RefreshTokenStore
	userClient        userpb.UserServiceClient
}

func NewAuthService(
	db *sql.DB,
	credentialStore *store.CredentialStore,
	refreshTokenStore *store.RefreshTokenStore,
	userClient userpb.UserServiceClient,
) *AuthService {
	return &AuthService{
		db:                db,
		credentialStore:   credentialStore,
		refreshTokenStore: refreshTokenStore,
		userClient:        userClient,
	}
}

func (s *AuthService) SignUp(
	ctx context.Context,
	req dto.SignUpRequest,
) (*dto.SignUpResponse, error) {
	serviceToken := os.Getenv("AUTH_SERVICE_TOKEN")
	if serviceToken == "" {
		return nil, errors.New("AUTH_SERVICE_TOKEN is not configured")
	}

	userCtx := metadata.AppendToOutgoingContext(
		ctx,
		"authorization",
		"Bearer "+serviceToken,
	)
	userResp, err := s.userClient.CreateUser(userCtx, &userpb.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})

	if err != nil {
		return nil, err
	}

	// hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}
	log.Println("out begin", userResp, err)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("begin", err)
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}

	defer tx.Rollback()

	user, err := s.credentialStore.Create(
		ctx,
		tx,
		userResp.Id,
		req.Email,
		passwordHash,
	)
	if err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}

	// generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, model.Role(user.Role), 15*time.Minute, "")
	if err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}
	refreshToken, expires_at, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}
	// store hashed refresh token
	err = s.refreshTokenStore.Create(
		ctx,
		tx,
		user.ID,
		refreshToken,
		expires_at,
	)
	if err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, s.compensateProfile(ctx, userResp.Id, err)
	}

	return &dto.SignUpResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) SignIn(
	ctx context.Context,
	req dto.SignInRequest,
) (*dto.SignInResponse, error) {

	user, err := s.credentialStore.GetByEmail(
		ctx,
		s.db,
		req.Email,
	)
	if err != nil {
		return nil, err
	}

	if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, model.Role(user.Role), 15*time.Minute, "")
	if err != nil {
		return nil, err
	}
	refreshToken, expires_at, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	err = s.refreshTokenStore.Create(
		ctx,
		s.db,
		user.ID,
		refreshToken,
		expires_at,
	)
	if err != nil {
		return nil, err
	}

	return &dto.SignInResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	req dto.RefreshTokenRequest,
) (*dto.RefreshTokenResponse, error) {

	userID, err := s.refreshTokenStore.GetByHash(
		ctx,
		s.db,
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	user, err := s.credentialStore.GetByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Email, model.Role(user.Role), 15*time.Minute, "")
	if err != nil {
		return nil, err
	}
	newRefreshToken, expires_at, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err := s.refreshTokenStore.Revoke(ctx, tx, req.RefreshToken); err != nil {
		return nil, err
	}
	// store hashed refresh token
	err = s.refreshTokenStore.Create(
		ctx,
		tx,
		user.ID,
		newRefreshToken,
		expires_at,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) SignOut(
	ctx context.Context,
	req dto.RefreshTokenRequest,
) (*dto.SignOutResponse, error) {

	err := s.refreshTokenStore.Revoke(
		ctx,
		s.db,
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	return &dto.SignOutResponse{
		Message: "logged out successfully",
	}, nil
}
