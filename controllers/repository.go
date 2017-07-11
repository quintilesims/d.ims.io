package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

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
				"GET": r.ListRepositories,
			},
		},
		{
			Path: "/repository/:owner",
			Handlers: fireball.Handlers{
				"POST": r.CreateRepository,
				"GET":  r.ListOwnerRepositories,
			},
		},
		{
			Path: "/repository/:owner/:name",
			Handlers: fireball.Handlers{
				"GET":    r.GetRepository,
				"DELETE": r.DeleteRepository,
			},
		},
	}
}

func (r *RepositoryController) CreateRepository(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]

	var req models.CreateRepositoryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return fireball.NewJSONError(400, err)
	}

	if err := req.Validate(); err != nil {
		return nil, fireball.NewError(400, err, nil)
	}

	repo := fmt.Sprintf("%s/%s", owner, req.Name)
	input := &ecr.CreateRepositoryInput{}
	input.SetRepositoryName(repo)
	if err := input.Validate(); err != nil {
		return nil, fireball.NewError(400, err, nil)
	}

	if _, err := r.ecr.CreateRepository(input); err != nil {
		return nil, err
	}

	resp := models.CreateRepositoryResponse{
		Owner: owner,
		Name:  req.Name,
	}

	return fireball.NewJSONResponse(202, resp)
}

func (r *RepositoryController) DeleteRepository(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]
	name := c.PathVariables["name"]
	repo := fmt.Sprintf("%s/%s", owner, name)

	input := &ecr.DeleteRepositoryInput{}
	input.SetRepositoryName(repo)
	input.SetForce(true)

	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	if _, err := r.ecr.DeleteRepository(input); err != nil {
		return nil, err
	}

	return fireball.NewResponse(200, []byte("Successfully deleted repository"), nil), nil
}

func (r *RepositoryController) GetRepository(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]
	name := c.PathVariables["name"]
	repo := fmt.Sprintf("%s/%s", owner, name)

	input := &ecr.DescribeImagesInput{}
	input.SetRepositoryName(repo)
	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
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
		Owner:     owner,
		Name:      name,
		ImageTags: tags,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	input := &ecr.DescribeRepositoriesInput{}
	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
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

func (r *RepositoryController) ListOwnerRepositories(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]

	input := &ecr.DescribeRepositoriesInput{}
	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	repositories := []string{}
	fn := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, repository := range output.Repositories {
			prefix := fmt.Sprintf("%s/", owner)
			repositoryName := aws.StringValue(repository.RepositoryName)

			if strings.HasPrefix(repositoryName, prefix) {
				repositories = append(repositories, strings.TrimPrefix(repositoryName, prefix))
			}
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
