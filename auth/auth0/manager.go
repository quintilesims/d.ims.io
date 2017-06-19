package auth0

type ADManager interface {
	Authenticate(user, pass string) (bool, error)
}
