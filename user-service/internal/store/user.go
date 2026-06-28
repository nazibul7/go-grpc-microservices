package store

import (
	"context"
	"database/sql"

	"github.com/nazibul7/go-grpc-microservices/user-service/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.User, error) {
	query := `INSERT INTO users(name, email) VALUES($1, $2) RETURNING id, name, email`

	var user model.User
	err := s.db.QueryRowContext(ctx, query, req.Name, req.Email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserStore) GetUser(ctx context.Context, id string) (model.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user model.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, err
		}
		return model.User{}, err
	}
	return user, nil
}

func (s *UserStore) UpdateUser(ctx context.Context, req model.UpdateUserRequest) (model.User, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email`

	var user model.User
	err := s.db.QueryRowContext(ctx, query, req.Name, req.Email, req.ID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, err
		}
		return model.User{}, err
	}
	return user, nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id string) (model.DeleteUserResponse, error) {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return model.DeleteUserResponse{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.DeleteUserResponse{}, err
	}
	if rowsAffected == 0 {
		return model.DeleteUserResponse{}, sql.ErrNoRows
	}

	return model.DeleteUserResponse{Message: "user deleted successfully"}, nil
}
