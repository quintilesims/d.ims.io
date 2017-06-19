package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/quintilesims/d.ims.io/auth/token"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
	"math/rand"
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
	// todo: validate user:pass in auth0
	// todo: get username in AddToken

	raw := fmt.Sprintf("%s:%s", randomString(26), randomString(26))
	token := base64.StdEncoding.EncodeToString([]byte(raw))

	if err := t.tokenManager.AddToken("test", token); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.CreateTokenResponse{Token: token})
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	runes := make([]rune, length)
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return string(runes)
}
