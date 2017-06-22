package auth

func NewCompositeAuthenticator(authenticators ...Authenticator) AuthenticatorFunc {
	return AuthenticatorFunc(func(user, pass string) (bool, error) {
		for _, authenticator := range authenticators {
			isValid, err := authenticator.Authenticate(user, pass)
			if err != nil {
				return false, err
			}

			if isValid {
				return true, nil
			}
		}

		return false, nil
	})
}
