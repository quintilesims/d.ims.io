package token

type TokenManager interface {
	GenerateToken(user string) (string, error)
	DeleteToken(token string) error
	Authenticate(user, pass string) (bool, error)
}
