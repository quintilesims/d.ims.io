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
	bytesize "github.com/zpatrick/go-bytesize"
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
		{
			Path: "/repository/:owner/:name/image",
			Handlers: fireball.Handlers{
				"GET": r.ListRepositoryImages,
			},
		},
		{
			Path: "/repository/:owner/:name/image/:tag",
			Handlers: fireball.Handlers{
				"GET":    r.GetRepositoryImage,
				"DELETE": r.DeleteRepositoryImage,
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

	input := &ecr.DescribeRepositoriesInput{}
	input.SetRepositoryNames([]*string{aws.String(repo)})

	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	output, err := r.ecr.DescribeRepositories(input)
	if err != nil {
		return nil, err
	}

	resp := models.Repository{
		Owner:     owner,
		Name:      name,
		CreatedAt: aws.TimeValue(output.Repositories[0].CreatedAt),
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) ListRepositories(c *fireball.Context) (fireball.Response, error) {
	repositories := []string{}
	fn := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, repository := range output.Repositories {
			repositories = append(repositories, aws.StringValue(repository.RepositoryName))
		}

		return !lastPage
	}

	if err := r.ecr.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{}, fn); err != nil {
		return nil, err
	}

	resp := models.ListRepositoriesResponse{
		Repositories: repositories,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) ListOwnerRepositories(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]

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

	if err := r.ecr.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{}, fn); err != nil {
		return nil, err
	}

	resp := models.ListRepositoriesResponse{
		Repositories: repositories,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) ListRepositoryImages(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]
	name := c.PathVariables["name"]
	repo := fmt.Sprintf("%s/%s", owner, name)

	filter := &ecr.ListImagesFilter{}
	filter.SetTagStatus("TAGGED")

	input := &ecr.ListImagesInput{}
	input.SetRepositoryName(repo)
	input.SetFilter(filter)

	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	images := []string{}
	fn := func(output *ecr.ListImagesOutput, lastPage bool) bool {
		for _, image := range output.ImageIds {
			images = append(images, aws.StringValue(image.ImageTag))
		}

		return !lastPage
	}

	if err := r.ecr.ListImagesPages(input, fn); err != nil {
		return nil, err
	}

	resp := models.ListImagesResponse{
		Images: images,
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) GetRepositoryImage(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]
	name := c.PathVariables["name"]
	tag := c.PathVariables["tag"]
	repo := fmt.Sprintf("%s/%s", owner, name)

	imageID := &ecr.ImageIdentifier{}
	imageID.SetImageTag(tag)

	input := &ecr.DescribeImagesInput{}
	input.SetRepositoryName(repo)
	input.SetImageIds([]*ecr.ImageIdentifier{imageID})

	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	output, err := r.ecr.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	detail := output.ImageDetails[0]
	size := bytesize.Bytesize(aws.Int64Value(detail.ImageSizeInBytes))

	resp := models.Image{
		Digest:   aws.StringValue(detail.ImageDigest),
		Size:     size.Format("MB"),
		PushedAt: aws.TimeValue(detail.ImagePushedAt),
	}

	return fireball.NewJSONResponse(200, resp)
}

func (r *RepositoryController) DeleteRepositoryImage(c *fireball.Context) (fireball.Response, error) {
	owner := c.PathVariables["owner"]
	name := c.PathVariables["name"]
	tag := c.PathVariables["tag"]
	repo := fmt.Sprintf("%s/%s", owner, name)

	imageID := &ecr.ImageIdentifier{}
	imageID.SetImageTag(tag)

	input := &ecr.BatchDeleteImageInput{}
	input.SetRepositoryName(repo)
	input.SetImageIds([]*ecr.ImageIdentifier{imageID})

	if err := input.Validate(); err != nil {
		return fireball.NewJSONError(400, err)
	}

	output, err := r.ecr.BatchDeleteImage(input)
	if err != nil {
		return nil, err
	}

	if failures := output.Failures; len(failures) > 0 {
		return nil, fmt.Errorf(failures[0].String())
	}

	message := fmt.Sprintf("Image '%s:%s' successfully deleted.", repo, tag)
	return fireball.NewJSONResponse(200, message)
}
