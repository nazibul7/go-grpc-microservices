package service

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/dto"
	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/store"
	userpb "github.com/nazibul7/go-grpc-microservices/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockUserService struct {
	userpb.UnimplementedUserServiceServer
	called bool
}

func (s *mockUserService) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.UserResponse, error) {
	s.called=true

	return &userpb.UserResponse{
		Id: "test-user-id",
	}, nil
}

func TestSignUp_UserServiceUnavailable(t *testing.T) {
	conn, err := grpc.NewClient(
		"localhost:59999",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)

	authService := NewAuthService(
		&sql.DB{},
		store.NewCredentialStore(),
		store.NewRefreshTokenStore(),
		userClient,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	resp, err := authService.SignUp(ctx, dto.SignUpRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "TestPassword123",
	})

	if err == nil {
		t.Fatal("expected signup to fail when User Service is unavailable")
	}

	if resp != nil {
		t.Fatal("expected nil response when signup fails")
	}

	st, ok := status.FromError(err)

	if !ok {
		t.Fatalf("expected a gRPC status error, got: %v", err)
	}

	if st.Code() != codes.Unavailable {
		t.Fatalf("expected Unavailable, got %s: %v", st.Code(), err)
	}
}

func TestSignUp_AuthDBUnavailableAfterProfileCreation(t *testing.T) {
	// Create mock Auth DB.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock database: %v", err)
	}
	defer db.Close()

	// Start a local TCP listener on an available port.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Start mock User Service as a real gRPC server.
	mockUser := &mockUserService{}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, mockUser)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Logf("gRPC server stopped: %v", err)
		}
	}()
	defer grpcServer.Stop()

	// Connect Auth Service to mock User Service.
	conn, err := grpc.NewClient(
		listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)

	// Make Auth DB transaction creation fail.
	authDBErr := errors.New("auth database unavailable")
	mock.ExpectBegin().WillReturnError(authDBErr)

	authService := NewAuthService(
		db,
		nil,
		nil,
		userClient,
	)

	resp, err := authService.SignUp(
		context.Background(),
		dto.SignUpRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		},
	)

	if err == nil {
		t.Fatal("expected signup error, got nil")
	}

	if resp != nil {
		t.Fatal("expected nil signup response")
	}

	// Verify User Service was called.
	if !mockUser.called {
		t.Fatal("expected User Service CreateUser to be called")
	}

	// Verify Auth DB transaction was attempted.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled database expectations: %v", err)
	}
}
