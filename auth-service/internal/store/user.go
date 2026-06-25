package store

import (
	"context"

	"github.com/nazibul7/go-grpc-microservices/auth-service/internal/model"
)

type UserStore struct{}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (s *UserStore) Create(
	ctx context.Context,
	db DBTX,
	email string,
	passwordHash string,
) (*model.User, error) {

	query := `
		INSERT INTO users(email, password_hash)
		VALUES($1, $2)
		RETURNING id, email
	`

	user := &model.User{}

	err := db.QueryRowContext(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByEmail(
	ctx context.Context,
	db DBTX,
	email string,
) (*model.User, error) {

	query := `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	user := &model.User{}

	err := db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByID(
	ctx context.Context,
	db DBTX,
	id int64,
) (*model.User, error) {

	query := `
		SELECT id, email, role
		FROM users
		WHERE id = $1
	`

	user := &model.User{}

	err := db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}