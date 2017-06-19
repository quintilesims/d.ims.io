package token

type TokenManager interface {
	AddToken(user, token string) error
}
