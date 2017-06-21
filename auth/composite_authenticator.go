package auth

type CompositeAuthenticator struct {
	authenticators []Authenticator
}

func NewCompositeAuthenticator(authenticators ...Authenticator) *CompositeAuthenticator {
	return &CompositeAuthenticator{
		authenticators: authenticators,
	}
}

func (c *CompositeAuthenticator) Authenticate(user, pass string) (bool, error) {
	for _, authenticator := range c.authenticators {
		isValid, err := authenticator.Authenticate(user, pass)
		if err != nil {
			return false, err
		}

		if isValid {
			return true, nil
		}
	}

	return false, nil
}
