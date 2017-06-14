package auth0

type Authenticator interface {
	Authenticate(username, password string) error
}
