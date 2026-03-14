package domain

type TokenService interface {
	GenerateAccessToken(user *User) (string, error)
	GenerateRefreshToken() (raw string, hash string, err error) // raw → cookie, hash → DB
	ParseAccessToken(tokenStr string) (*AccessTokenClaims, error)
	HashToken(raw string) string
}
