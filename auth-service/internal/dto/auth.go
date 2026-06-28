package dto

type SignUpRequest struct {
	Name     string
	Email    string
	Password string
}

type SignUpResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type RefreshTokenResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type SignOutresponse struct {
	Message string
}
