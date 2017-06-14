package token

type Authenticator interface {
	AddToken(user, token string) error
	Authenticate(token string) error
}
