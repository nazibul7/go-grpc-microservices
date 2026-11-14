package model

type Credential struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
}
