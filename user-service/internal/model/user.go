package model

type User struct {
	ID    int64
	Name  string
	Email string
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type IDRequest struct {
	ID int64
}

type UpdateUserRequest struct {
	ID    int64
	Name  string
	Email string
}

type DeleteUserResponse struct {
	Message string
}