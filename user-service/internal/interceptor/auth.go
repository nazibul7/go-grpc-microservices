package interceptor

import (
	"context"
	"strings"

	"github.com/nazibul7/go-grpc-microservices/user-service/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	if info.FullMethod == "/user.UserService/CreateUser" {
		return handler(ctx, req)
	}
	// Read incoming metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := values[0]

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
	}

	token := parts[1]

	// Verify JWT.
	claims, err := utils.VerifyToken(token, "")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Store claims in context so handlers/services can use them.
	ctx = context.WithValue(ctx, ClaimsKey, claims)

	// Continue to the actual RPC handler.
	return handler(ctx, req)
}
