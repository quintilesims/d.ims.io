package controllers

import (
	"fmt"
	//"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/zpatrick/fireball"
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

func (r *RepositoryController) CreateRepository(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *RepositoryController) GetRepository(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *RepositoryController) DeleteRepository(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (r *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("Not implemented")
}
