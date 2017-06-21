package auth

type Authenticator interface {
	Authenticate(user, pass string) (bool, error)
}
