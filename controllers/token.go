package controllers

import (
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
		{
			Path: "/token/:token",
			Handlers: fireball.Handlers{
				"DELETE": r.DeleteToken,
			},
		},
	}
}

func (t *TokenController) CreateToken(c *fireball.Context) (fireball.Response, error) {
	user, _, _ := c.Request.BasicAuth()
	token, err := t.tokenManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateTokenResponse{Token: token})
}

func (t *TokenController) DeleteToken(c *fireball.Context) (fireball.Response, error) {
	tk := c.PathVariables["token"]

	if err := t.tokenManager.DeleteToken(tk); err != nil {
		return nil, err
	}

	return fireball.NewResponse(200, []byte("Deleted token."), nil), nil
}
