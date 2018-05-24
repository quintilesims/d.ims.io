package auth

import "fmt"

func NewCompositeAuthenticator(authenticators ...Authenticator) AuthenticatorFunc {
	return AuthenticatorFunc(func(user, pass string) (bool, error) {
		if user == "" || pass == "" {
			return false, fmt.Errorf("username and/or password is empty")
		}

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
