package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/dto"
	userpb "github.com/nazibul7/go-grpc-microservices/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Separate mock for compensation failure.
type failingDeleteUserService struct {
	mockUserService
}

func (s *failingDeleteUserService) DeleteUser(
	ctx context.Context,
	req *userpb.DeleteUserRequest,
) (*userpb.DeleteUserResponse, error) {
	s.deleteCalled = true

	return nil, status.Error(
		codes.Unavailable,
		"user service unavailable",
	)
}

func TestSignUp_CompensatesWhenAuthDBTransactionFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockUser := &mockUserService{}
	userClient := startMockUserServer(t, mockUser)

	dbErr := errors.New("auth database unavailable")
	mock.ExpectBegin().WillReturnError(dbErr)

	authService := NewAuthService(db, nil, nil, userClient)

	resp, err := authService.SignUp(
		context.Background(),
		dto.SignUpRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		},
	)

	if err == nil {
		t.Fatal("expected signup error")
	}

	if resp != nil {
		t.Fatal("expected nil response")
	}

	if !mockUser.createCalled {
		t.Fatal("expected CreateUser to be called")
	}

	if !mockUser.deleteCalled {
		t.Fatal("expected DeleteUser compensation to be called")
	}

	if !errors.Is(err, dbErr) {
		t.Fatalf("expected database error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSignUp_CompensationFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockUser := &failingDeleteUserService{}
	userClient := startMockUserServer(t, mockUser)

	dbErr := errors.New("auth database unavailable")
	mock.ExpectBegin().WillReturnError(dbErr)

	authService := NewAuthService(db, nil, nil, userClient)

	resp, err := authService.SignUp(
		context.Background(),
		dto.SignUpRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		},
	)

	if err == nil {
		t.Fatal("expected signup error")
	}

	if resp != nil {
		t.Fatal("expected nil response")
	}

	if !mockUser.deleteCalled {
		t.Fatal("expected DeleteUser compensation to be called")
	}

	// Original DB error must be preserved.
	if !errors.Is(err, dbErr) {
		t.Errorf("expected database error, got: %v", err)
	}

	// Compensation gRPC error must be preserved.
	if status.Code(err) != codes.Unavailable {
		t.Errorf("expected Unavailable, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
