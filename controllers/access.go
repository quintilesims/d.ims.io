package controllers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
)

type AccessController struct {
	ecr    ecriface.ECRAPI
	access auth.AccessManager
}

func NewAccessController(e ecriface.ECRAPI, a auth.AccessManager) *AccessController {
	return &AccessController{
		ecr:    e,
		access: a,
	}
}

func (a *AccessController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/account",
			Handlers: fireball.Handlers{
				"GET": a.ListAccounts,
			},
		},
		{
			Path: "/account/",
			Handlers: fireball.Handlers{
				"POST": a.GrantAccess,
			},
		},
		{
			Path: "/account/:id",
			Handlers: fireball.Handlers{
				"DELETE": a.RevokeAccess,
			},
		},
	}
}

func (a *AccessController) ListAccounts(c *fireball.Context) (fireball.Response, error) {
	response, err := a.access.Accounts()
	if err != nil {
		return fireball.NewJSONError(500, err)
	}

	return fireball.NewJSONResponse(200, response)
}

func (a *AccessController) GrantAccess(c *fireball.Context) (fireball.Response, error) {
	var request models.GrantAccessRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		return nil, err
	}

	if err := a.access.GrantAccess(request.Account); err != nil {
		return fireball.NewJSONError(500, err)
	}

	return fireball.NewJSONResponse(200, nil)
}

func (a *AccessController) RevokeAccess(c *fireball.Context) (fireball.Response, error) {
	if err := a.access.RevokeAccess(c.PathVariables["id"]); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}
