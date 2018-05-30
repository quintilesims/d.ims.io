package controllers

import (
	"log"

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
				log.Printf("[DEBUG] Request %s %s did not contain basic auth", c.Request.Method, c.Request.URL.String())
				return invalidAuthResponse, nil
			}

			isAuthenticated, err := auth.Authenticate(user, pass)
			if err != nil {
				log.Printf("[ERROR] Authenticator encountered an unexpected error: %v", err)
				return nil, err
			}

			if !isAuthenticated {
				log.Printf("[DEBUG] User '%s' failed to authenticate: %v", user, err)
				return invalidAuthResponse, nil
			}

			return handler(c)
		}
	}
}
