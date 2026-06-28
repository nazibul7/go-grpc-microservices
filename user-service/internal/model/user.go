package model

type User struct {
	ID    string
	Name  string
	Email string
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type IDRequest struct {
	ID string
}

type UpdateUserRequest struct {
	ID    string
	Name  string
	Email string
}

type DeleteUserResponse struct {
	Message string
}
