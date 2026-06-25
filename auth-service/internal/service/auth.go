package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/dto"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/store"
)

type AuthService struct {
	db                *sql.DB
	userStore         *store.UserStore
	refreshTokenStore *store.RefreshTokenStore
}

func NewAuthService(
	db *sql.DB,
	userStore *store.UserStore,
	refreshTokenStore *store.RefreshTokenStore,
) *AuthService {
	return &AuthService{
		db:                db,
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
	}
}

func (s *AuthService) SignUp(
	ctx context.Context,
	req dto.SignUpRequest,
) (*dto.SignUpResponse, error) {

	// hash password
	passwordHash := req.Password // replace with bcrypt

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	user, err := s.userStore.Create(
		ctx,
		tx,
		req.Email,
		passwordHash,
	)
	if err != nil {
		return nil, err
	}

	// generate tokens
	accessToken := "access-token"
	refreshToken := "refresh-token"

	// store hashed refresh token
	err = s.refreshTokenStore.Create(
		ctx,
		tx,
		user.ID,
		refreshToken,
		time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
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

	user, err := s.userStore.GetByEmail(
		ctx,
		s.db,
		req.Email,
	)
	if err != nil {
		return nil, err
	}

	// verify password

	accessToken := "access-token"
	refreshToken := "refresh-token"

	err = s.refreshTokenStore.Create(
		ctx,
		s.db,
		user.ID,
		refreshToken,
		time.Now().Add(7*24*time.Hour),
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

	user, err := s.userStore.GetByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		return nil, err
	}

	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"

	return &dto.RefreshTokenResponse{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) SignOut(
	ctx context.Context,
	req dto.RefreshTokenRequest,
) (*dto.SignOutresponse, error) {

	err := s.refreshTokenStore.Revoke(
		ctx,
		s.db,
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	return &dto.SignOutresponse{
		Message: "logged out successfully",
	}, nil
}
