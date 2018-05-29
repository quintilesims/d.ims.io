package controllers

import (
	"testing"

	"github.com/quintilesims/d.ims.io/auth"
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/fireball"
)

func TestAuthDecoratorHonorsValidAuth(t *testing.T) {
	handler := func(c *fireball.Context) (fireball.Response, error) {
		return fireball.NewResponse(200, nil, nil), nil
	}

	authenticator := auth.AuthenticatorFunc(func(user, pass string) (bool, error) {
		return true, nil
	})

	c := newContextWithBasicAuth(t, "user", "pass")
	resp, err := AuthDecorator(authenticator)(handler)(c)
	if err != nil {
		t.Fatal(err)
	}

	assertResponseCode(t, resp, 200)
}

func TestAuthDecoratorHonorsInvalidAuth(t *testing.T) {
	handler := func(c *fireball.Context) (fireball.Response, error) {
		t.Fatal("handler was called")
		return nil, nil
	}

	authenticator := auth.AuthenticatorFunc(func(user, pass string) (bool, error) {
		return false, nil
	})

	c := newContextWithBasicAuth(t, "user", "pass")
	resp, err := AuthDecorator(authenticator)(handler)(c)
	if err != nil {
		t.Fatal(err)
	}

	assertResponseCode(t, resp, 401)
}

func TestAuthDecoratorCachesCredStatus(t *testing.T) {
	cases := map[string]bool{
		"ValidCreds":   true,
		"InvalidCreds": false,
	}

	for name, isValidCreds := range cases {
		t.Run(name, func(t *testing.T) {
			handler := func(c *fireball.Context) (fireball.Response, error) {
				return fireball.NewResponse(200, nil, nil), nil
			}

			var authenticatorCalls int
			authenticator := auth.AuthenticatorFunc(func(user, pass string) (bool, error) {
				authenticatorCalls++
				return isValidCreds, nil
			})

			// use the same decorated handler for multiple calls to
			// ensure we only use a single cache
			handler = AuthDecorator(authenticator)(handler)
			for i := 0; i < 5; i++ {
				c := newContextWithBasicAuth(t, "user", "pass")
				resp, err := handler(c)
				if err != nil {
					t.Fatal(err)
				}

				if isValidCreds {
					assertResponseCode(t, resp, 200)
				} else {
					assertResponseCode(t, resp, 401)
				}
			}

			assert.Equal(t, 1, authenticatorCalls)
		})
	}
}
