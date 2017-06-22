package controllers

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/quintilesims/d.ims.io/models"
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
	var req models.CreateRepositoryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, err
	}

	input := &ecr.CreateRepositoryInput{}
	input.SetRepositoryName(req.Name)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	if _, err := r.ecr.CreateRepository(input); err != nil {
		return nil, err
	}

	resp := models.CreateRepositoryResponse{
		Name: req.Name,
	}

	return fireball.NewJSONResponse(202, resp)
}

func (r *RepositoryController) DeleteRepository(c *fireball.Context) (fireball.Response, error) {
	name := c.PathVariables["name"]
	input := &ecr.DeleteRepositoryInput{}
	input.SetRepositoryName(name)
	input.SetForce(true)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if _, err := r.ecr.DeleteRepository(input); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (r *RepositoryController) GetRepository(c *fireball.Context) (fireball.Response, error) {
	name := c.PathVariables["name"]
	input := &ecr.DescribeImagesInput{}
	input.SetRepositoryName(name)
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tags := []string{}
	fn := func(output *ecr.DescribeImagesOutput, lastPage bool) bool {
		for _, image := range output.ImageDetails {
			for _, tag := range image.ImageTags {
				tags = append(tags, aws.StringValue(tag))
			}
		}

		return !lastPage
	}

	if err := r.ecr.DescribeImagesPages(input, fn); err != nil {
		return nil, err
	}

	resp := models.Repository{
		Name:      name,
		ImageTags: tags,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	input := &ecr.DescribeRepositoriesInput{}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	repositories := []string{}
	fn := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, repository := range output.Repositories {
			repositories = append(repositories, aws.StringValue(repository.RepositoryName))
		}

		return !lastPage
	}

	if err := r.ecr.DescribeRepositoriesPages(input, fn); err != nil {
		return nil, err
	}

	resp := models.ListRepositoriesResponse{
		Repositories: repositories,
	}

	return fireball.NewJSONResponse(200, resp)
}
