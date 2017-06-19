package controllers

import (
	"fmt"
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/zpatrick/fireball"
)

func AuthDecorator(auth auth.Authenticator) fireball.Decorator {
	return func(handler fireball.Handler) fireball.Handler {
		return func(c *fireball.Context) (fireball.Response, error) {
			isAuthenticated, err := auth.Authenticate(c.Request)
			if err != nil {
				return nil, err
			}

			if !isAuthenticated {
				return nil, fireball.NewError(401, fmt.Errorf("Invalid Auth"), nil)
			}

			return handler(c)
		}
	}
}
