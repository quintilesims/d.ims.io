package controllers

import (
	"fmt"

	"github.com/quintilesims/d.ims.io/auth"
	"github.com/zpatrick/fireball"
)

type TokenController struct {
	tokenManager auth.TokenManager
}

func NewTokenController(t auth.TokenManager) *TokenController {
	return &TokenController{
		tokenManager: t,
	}
}

func (t *TokenController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/token",
			Handlers: fireball.Handlers{
				"DELETE": t.DeleteToken,
				"POST":   t.CreateToken,
			},
		},
	}
}

func (t *TokenController) CreateToken(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (t *TokenController) DeleteToken(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}
