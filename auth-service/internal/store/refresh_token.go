package store

import (
	"context"
	"time"
)

type RefreshTokenStore struct{}

func NewRefreshTokenStore() *RefreshTokenStore {
	return &RefreshTokenStore{}
}

func (s *RefreshTokenStore) Create(
	ctx context.Context,
	db DBTX,
	userID string,
	tokenHash string,
	expiresAt time.Time,
) error {

	query := `
		INSERT INTO refresh_tokens(
			user_id,
			token_hash,
			expires_at
		)
		VALUES($1, $2, $3)
	`

	_, err := db.ExecContext(
		ctx,
		query,
		userID,
		tokenHash,
		expiresAt,
	)

	return err
}

func (s *RefreshTokenStore) GetByHash(
	ctx context.Context,
	db DBTX,
	tokenHash string,
) (string, error) {

	query := `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		AND revoked = FALSE
		AND expires_at > NOW()
	`

	var userID string

	err := db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *RefreshTokenStore) Revoke(
	ctx context.Context,
	db DBTX,
	tokenHash string,
) error {

	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE token_hash = $1
	`

	_, err := db.ExecContext(
		ctx,
		query,
		tokenHash,
	)

	return err
}
