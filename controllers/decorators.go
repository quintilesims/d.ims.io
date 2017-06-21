package controllers

import (
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/zpatrick/fireball"
)

func AuthDecorator(auth auth.Authenticator) fireball.Decorator {
	return func(handler fireball.Handler) fireball.Handler {
		return func(c *fireball.Context) (fireball.Response, error) {
			headers := map[string]string{"WWW-Authenticate": "Basic realm=\"Restricted\""}
			invalidAuthResponse := fireball.NewResponse(401, []byte("401 Unauthorized\n"), headers)

			user, pass, ok := c.Request.BasicAuth()
			if !ok {
				return invalidAuthResponse, nil
			}

			isAuthenticated, err := auth.Authenticate(user, pass)
			if err != nil {
				return nil, err
			}

			if !isAuthenticated {
				return invalidAuthResponse, nil
			}

			return handler(c)
		}
	}
}
