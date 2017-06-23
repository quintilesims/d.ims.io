package controllers

import (
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/quintilesims/d.ims.io/models"
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
				"POST": t.CreateToken,
			},
		},
		{
			Path: "/token/:token",
			Handlers: fireball.Handlers{
				"DELETE": t.DeleteToken,
			},
		},
	}
}

func (t *TokenController) CreateToken(c *fireball.Context) (fireball.Response, error) {
	user, _, _ := c.Request.BasicAuth()
	token, err := t.tokenManager.CreateToken(user)
	if err != nil {
		return nil, err
	}

	resp := models.CreateTokenResponse{
		Token: token,
	}

	return fireball.NewJSONResponse(202, resp)
}

func (t *TokenController) DeleteToken(c *fireball.Context) (fireball.Response, error) {
	token := c.PathVariables["token"]
	if err := t.tokenManager.DeleteToken(token); err != nil {
		return nil, err
	}

	return fireball.NewResponse(200, []byte("Successfully deleted token"), nil), nil
}
