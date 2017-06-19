package auth

type Authenticator interface {
	Authenticate(username, password string) (bool, error)
}
