package controllers

import (
	"fmt"
	"github.com/quintilesims/d.ims.io/auth/token"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
)

type TokenController struct {
	tokenManager token.TokenManager
}

func NewTokenController(t token.TokenManager) *TokenController {
	return &TokenController{
		tokenManager: t,
	}
}

func (r *TokenController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/token",
			Handlers: fireball.Handlers{
				"POST": r.CreateToken,
			},
		},
	}
}

func (t *TokenController) CreateToken(c *fireball.Context) (fireball.Response, error) {
	// todo: auth0 auth
	user, _, ok := c.Request.BasicAuth()
	if !ok {
		// todo: use stnadard unauthorized response
		return nil, fmt.Errorf("Must pass auth")
	}

	token, err := t.tokenManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateTokenResponse{Token: token})
}
