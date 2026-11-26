package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/nazibul7/go-grpc-microservices/proto/user"
	"google.golang.org/grpc/metadata"
)

func (s *AuthService) compensateProfile(ctx context.Context, userID string, cause error) error {
	compensationCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()

	serviceToken := os.Getenv("AUTH_SERVICE_TOKEN")
	if serviceToken == "" {
		return errors.Join(
			cause,
			errors.New("AUTH_SERVICE_TOKEN is not configured"),
		)
	}

	compensationCtx = metadata.AppendToOutgoingContext(
		compensationCtx,
		"authorization",
		"Bearer "+serviceToken,
	)

	_, err := s.userClient.DeleteUser(compensationCtx, &user.DeleteUserRequest{Id: userID})

	if err != nil {
		return errors.Join(cause, fmt.Errorf("compensation failed for user %s: %w", userID, err))
	}

	return cause
}
