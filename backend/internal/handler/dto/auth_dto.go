package dto

// GoogleAuthRequest — body for POST /auth/google
type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

// TokenResponse — returned after successful authentication
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
