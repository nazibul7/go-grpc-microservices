package dto

type SignUpRequest struct {
	Email    string
	Password string
}

type SignUpResponse struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type RefreshTokenResponse struct {
	UserID       int64
	AccessToken  string
	RefreshToken string
}

type SignOutresponse struct {
	Message string
}
