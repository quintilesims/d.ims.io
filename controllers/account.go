package controllers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
)

type AccountController struct {
	ecr    ecriface.ECRAPI
	access auth.AccountManager
}

func NewAccountController(e ecriface.ECRAPI, a auth.AccountManager) *AccountController {
	return &AccountController{
		ecr:    e,
		access: a,
	}
}

func (a *AccountController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/account",
			Handlers: fireball.Handlers{
				"GET":  a.ListAccounts,
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

func (a *AccountController) ListAccounts(c *fireball.Context) (fireball.Response, error) {
	response, err := a.access.Accounts()
	if err != nil {
		return fireball.NewJSONError(500, err)
	}

	return fireball.NewJSONResponse(200, models.ListAccountsResponse{Accounts: response})
}

func (a *AccountController) GrantAccess(c *fireball.Context) (fireball.Response, error) {
	var request models.GrantAccessRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		return fireball.NewJSONError(400, err)
	}

	repositories, err := listRepositories(a.ecr)
	if err != nil {
		return fireball.NewJSONError(500, err)
	}

	accounts, err := a.access.Accounts()
	if err != nil {
		return fireball.NewJSONError(500, err)
	}

	for _, r := range repositories {
		if err := addToRepositoryPolicy(a.ecr, r, append(accounts, request.Account)); err != nil {
			return fireball.NewJSONError(500, err)
		}
	}

	if err := a.access.GrantAccess(request.Account); err != nil {
		return fireball.NewJSONError(500, err)
	}

	return fireball.NewJSONResponse(204, nil)
}

func (a *AccountController) RevokeAccess(c *fireball.Context) (fireball.Response, error) {
	accountID := c.PathVariables["id"]
	repositories, err := listRepositories(a.ecr)
	if err != nil {
		return fireball.NewJSONError(500, err)
	}

	for _, r := range repositories {
		if err := removeFromRepositoryPolicy(a.ecr, r, accountID); err != nil {
			return fireball.NewJSONError(500, err)
		}
	}

	if err := a.access.RevokeAccess(accountID); err != nil {
		return fireball.NewJSONError(500, err)
	}

	return fireball.NewJSONResponse(204, nil)
}
