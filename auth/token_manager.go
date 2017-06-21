package auth

type TokenManager interface {
	CreateToken(user string) (string, error)
	DeleteToken(token string) error
}
