package auth

import (
	"testing"
)

func newTestAuthenticator(isValid bool, err error) AuthenticatorFunc {
	return AuthenticatorFunc(func(string, string) (bool, error) {
		return isValid, err
	})
}

func TestEmptyUserPass(t *testing.T) {
	target := NewCompositeAuthenticator()
	if _, err := target.Authenticate("", ""); err == nil {
		t.Fatalf("Error expected when authenticating with no user and pass")
	}
}

func TestCompositeAuthenticator(t *testing.T) {
	cases := []struct {
		Name           string
		ExpectedResult bool
		Authenticators []Authenticator
	}{
		{
			Name:           "empty",
			ExpectedResult: false,
		},
		{
			Name:           "all true",
			ExpectedResult: true,
			Authenticators: []Authenticator{
				newTestAuthenticator(true, nil),
				newTestAuthenticator(true, nil),
				newTestAuthenticator(true, nil),
			},
		},
		{
			Name:           "all false",
			ExpectedResult: false,
			Authenticators: []Authenticator{
				newTestAuthenticator(false, nil),
				newTestAuthenticator(false, nil),
				newTestAuthenticator(false, nil),
			},
		},
		{
			Name:           "mixed",
			ExpectedResult: true,
			Authenticators: []Authenticator{
				newTestAuthenticator(false, nil),
				newTestAuthenticator(true, nil),
				newTestAuthenticator(false, nil),
				newTestAuthenticator(true, nil),
			},
		},
	}

	for _, c := range cases {
		target := NewCompositeAuthenticator(c.Authenticators...)
		result, err := target.Authenticate("user", "pass")
		if err != nil {
			t.Fatalf("Error on case %s: %v", c.Name, err)
		}

		if v, want := result, c.ExpectedResult; v != want {
			t.Errorf("case %s: result was %v, expected %v", c.Name, v, want)
		}
	}
}
