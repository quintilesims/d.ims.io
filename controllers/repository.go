package controllers

import (
	"encoding/json"
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

	if _, err := r.ecr.CreateRepository(input); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(202, models.CreateRepositoryResponse{Name: req.Name})
}

func (r *RepositoryController) GetRepository(c *fireball.Context) (fireball.Response, error) {
	// todo: auth

	name := c.PathVariables["name"]
	describeReposInput := &ecr.DescribeRepositoriesInput{}
	describeReposInput.SetRepositoryNames([]*string{&name})
	if err := describeReposInput.Validate(); err != nil {
		return nil, err
	}

	describeReposOutput, err := r.ecr.DescribeRepositories(describeReposInput)
	if err != nil {
		return nil, err
	}

	repository := describeReposOutput.Repositories[0]

	describeImagesInput := &ecr.DescribeImagesInput{}
	describeImagesInput.SetRepositoryName(name)
	if err := describeImagesInput.Validate(); err != nil {
		return nil, err
	}

	tags := []string{}
	fn := func(output *ecr.DescribeImagesOutput, lastPage bool) bool {
		for _, image := range output.ImageDetails {
			for _, tag := range image.ImageTags {
				tags = append(tags, *tag)
			}
		}

		return !lastPage
	}

	if err := r.ecr.DescribeImagesPages(describeImagesInput, fn); err != nil {
		return nil, err
	}

	resp := models.Repository{
		Name:      *repository.RepositoryName,
		URI:       *repository.RepositoryUri,
		ImageTags: tags,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) DeleteRepository(c *fireball.Context) (fireball.Response, error) {
	// todo: auth
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

func (r *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	input := &ecr.DescribeRepositoriesInput{}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	repositories := []string{}
	fn := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, repository := range output.Repositories {
			repositories = append(repositories, *repository.RepositoryName)
		}

		return !lastPage
	}

	if err := r.ecr.DescribeRepositoriesPages(input, fn); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, models.ListRepositoriesResponse{Repositories: repositories})
}
