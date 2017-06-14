package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type RepositoryController struct {
	ecr ecriface.ECRAPI
}

func NewRepositoryController(e ecriface.ECRAPI) *RepositoryController {
	return &RepositoryController{
		ecr: e,
	}
}

func (r *RepositoryController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/repository",
			Handlers: fireball.Handlers{
				"GET":  r.ListRepositories,
				"POST": r.CreateRepository,
			},
		},
		{
			Path: "/repository/:name",
			Handlers: fireball.Handlers{
				"GET":    r.GetRepository,
				"DELETE": r.DeleteRepository,
			},
		},
	}
}

func (t *RepositoryController) CreateRepository(c *fireball.Context) (fireball.Response, error) {
	// todo: auth

	var req models.CreateRepositoryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, err
	}

	input := &ecr.CreateRepositoryInput{}
	input.SetRepositoryName(req.Name)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if _, err := t.ecr.CreateRepository(input); err != nil {
		return nil, err
	}
	
	return fireball.NewJSONResponse(202, models.CreateRepositoryResponse{Name: req.Name})
}

func (t *RepositoryController) GetRepository(c *fireball.Context) (fireball.Response, error) {
	// todo: auth
	return nil, fmt.Errorf("Not implemented")
}

func (t *RepositoryController) DeleteRepository(c *fireball.Context) (fireball.Response, error) {
	// todo: auth
	return nil, fmt.Errorf("Not implemented")
}

func (t *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	// todo: auth
	return nil, fmt.Errorf("Not implemented")
}
