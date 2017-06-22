package auth

type Authenticator interface {
	Authenticate(user, pass string) (bool, error)
}

type AuthenticatorFunc func(string, string) (bool, error)

func (a AuthenticatorFunc) Authenticate(user, pass string) (bool, error) {
	return a(user, pass)
}
