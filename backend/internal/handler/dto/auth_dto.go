package dto

type RegisterRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}
