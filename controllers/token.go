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

func (r *TokenController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/token",
			Handlers: fireball.Handlers{
				"DELETE": r.DeleteToken,
				"POST":   r.CreateToken,
			},
		},
	}
}

func (r *TokenController) CreateToken(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *TokenController) DeleteToken(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}
