package store

import (
	"context"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/model"
)

type CredentialStore struct{}

func NewCredentialStore() *CredentialStore {
    return &CredentialStore{}
}

func (s *CredentialStore) Create(
	ctx context.Context,
	db DBTX,
	userID string,
	email string,
	passwordHash string,
) (*model.Credential, error) {

	query := `
		INSERT INTO credentials(id, email, password_hash)
		VALUES($1, $2, $3)
		RETURNING id, email
	`

	Credential := &model.Credential{}

	err := db.QueryRowContext(
		ctx,
		query,
		userID,
		email,
		passwordHash,
	).Scan(
		&Credential.ID,
		&Credential.Email,
	)

	if err != nil {
		return nil, err
	}

	return Credential, nil
}

func (s *CredentialStore) GetByEmail(
	ctx context.Context,
	db DBTX,
	email string,
) (*model.Credential, error) {

	query := `
		SELECT id, email, password_hash, role
		FROM credentials
		WHERE email = $1
	`

	Credential := &model.Credential{}

	err := db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&Credential.ID,
		&Credential.Email,
		&Credential.PasswordHash,
		&Credential.Role,
	)

	if err != nil {
		return nil, err
	}

	return Credential, nil
}

func (s *CredentialStore) GetByID(
	ctx context.Context,
	db DBTX,
	id string,
) (*model.Credential, error) {

	query := `
		SELECT id, email, role
		FROM credentials
		WHERE id = $1
	`

	Credential := &model.Credential{}

	err := db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&Credential.ID,
		&Credential.Email,
		&Credential.Role,
	)

	if err != nil {
		return nil, err
	}

	return Credential, nil
}
