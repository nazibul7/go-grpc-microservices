package service

import (
	"context"
	"database/sql"
	"net"
	"testing"
	"time"

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

	createCalled bool
	deleteCalled bool
}

func (s *mockUserService) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.UserResponse, error) {
	s.createCalled = true

	return &userpb.UserResponse{
		Id: "test-user-id",
	}, nil
}

func (s *mockUserService) DeleteUser(
	ctx context.Context,
	req *userpb.DeleteUserRequest,
) (*userpb.DeleteUserResponse, error) {
	s.deleteCalled = true

	return &userpb.DeleteUserResponse{
		Message: "user deleted",
	}, nil
}

// Starts a real local gRPC server for integration tests.
func startMockUserServer(
	t *testing.T,
	mockUser userpb.UserServiceServer,
) userpb.UserServiceClient {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, mockUser)

	go func() {
		_ = grpcServer.Serve(listener)
	}()

	conn, err := grpc.NewClient(
		listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		grpcServer.Stop()
		listener.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		conn.Close()
		grpcServer.Stop()
		listener.Close()
	})

	return userpb.NewUserServiceClient(conn)
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
		t.Fatal("expected signup to fail")
	}

	if resp != nil {
		t.Fatal("expected nil response")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got: %v", err)
	}

	if st.Code() != codes.Unavailable {
		t.Fatalf("expected Unavailable, got: %v", st.Code())
	}
}

func TestSignUp_AuthDBUnavailableAfterProfileCreation(t *testing.T) {
	mockUser := &mockUserService{}

	userClient := startMockUserServer(t, mockUser)

	// Auth DB is deliberately unavailable.
	// The actual SQL transaction is not created.
	db := &sql.DB{}

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
		t.Fatal("expected signup error")
	}

	if resp != nil {
		t.Fatal("expected nil response")
	}

	if !mockUser.createCalled {
		t.Fatal("expected CreateUser to be called")
	}
}
