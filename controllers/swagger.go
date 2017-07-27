package controllers

import (
	"encoding/json"

	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
	"github.com/zpatrick/go-plugin-swagger"
)

type SwaggerController struct {
	host string
}

func NewSwaggerController(host string) *SwaggerController {
	return &SwaggerController{
		host: host,
	}
}

func (s *SwaggerController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/swagger.json",
			Handlers: fireball.Handlers{
				"GET": s.serveSwaggerJSON,
			},
		},
	}

	return routes
}

func (s *SwaggerController) serveSwaggerJSON(c *fireball.Context) (fireball.Response, error) {
	spec := swagger.Spec{
		SwaggerVersion: "2.0",
		Host:           s.host,
		Schemes:        []string{"https"},
		Info: &swagger.Info{
			Title:   "D.IMS.IO",
			Version: "1.0.0",
		},
		Tags: []swagger.Tag{
			{
				Name:        "Token",
				Description: "Methods for Tokens",
			},
			{
				Name:        "Repository",
				Description: "Methods for Repositories",
			},
			{
				Name:        "Image",
				Description: "Methods for Images",
			},
		},
		Paths: map[string]swagger.Path{
			"/token": map[string]swagger.Method{
				"post": {
					Tags:       []string{"Token"},
					Summary:    "Create a new Token",
					Security:   swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("CreateTokenResponse"),
						},
					},
				},
			},
			"/token/{token}": map[string]swagger.Method{
				"delete": {
					Tags:     []string{"Token"},
					Summary:  "Delete a Token",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("token", "The token to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
						},
					},
				},
			},
			"/repository": map[string]swagger.Method{
				"get": {
					Tags:     []string{"Repository"},
					Summary:  "List all Repositories",
					Security: swagger.BasicAuthSecurity("login"),
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("ListRepositoriesResponse"),
						},
					},
				},
			},
			"/repository/{owner}": map[string]swagger.Method{
				"get": {
					Tags:     []string{"Repository"},
					Summary:  "List all Repositories for an Owner",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner of the Repositories", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("ListRepositoriesResponse"),
						},
					},
				},
				"post": {
					Tags:     []string{"Repository"},
					Summary:  "Create a new Repository",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner for the Repository", true),
						swagger.NewBodyParam("CreateRepositoryRequest", "none", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("CreateRepositoryResponse"),
						},
					},
				},
			},
			"/repository/{owner}/{name}": map[string]swagger.Method{
				"get": {
					Tags:     []string{"Repository"},
					Summary:  "Describe a Repository",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner of the Repository", true),
						swagger.NewStringPathParam("name", "Name of the Repository", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("Repository"),
						},
					},
				},
				"delete": {
					Tags:     []string{"Repository"},
					Summary:  "Delete a Repository",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner of the Repository", true),
						swagger.NewStringPathParam("name", "Name of the Repository", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
						},
					},
				},
			},
			"/images": map[string]swagger.Method{
				"get": {
					Tags:     []string{"Image"},
					Summary:  "List all Images",
					Security: swagger.BasicAuthSecurity("login"),
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("ListImagesResponse"),
						},
					},
				},
			},
			"/repository/{owner}/{name}/image/{tag}": map[string]swagger.Method{
				"get": {
					Tags:    []string{"Image"},
					Summary: "Describe an Image",
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner of the Repository", true),
						swagger.NewStringPathParam("name", "Name of the Repository", true),
						swagger.NewStringPathParam("tag", "Tag of the Image to delete", true),
					},
					Security: swagger.BasicAuthSecurity("login"),
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("Image"),
						},
					},
				},
				"delete": {
					Tags:    []string{"Image"},
					Summary: "Delete an Image",
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("owner", "Owner of the Repository", true),
						swagger.NewStringPathParam("name", "Name of the Repository", true),
						swagger.NewStringPathParam("tag", "Tag of the Image to delete", true),
					},
					Security: swagger.BasicAuthSecurity("login"),
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
						},
					},
				},
			},
		},
		Definitions: map[string]swagger.Definition{
			"CreateRepositoryRequest":  models.CreateRepositoryRequest{}.Definition(),
			"CreateRepositoryResponse": models.CreateRepositoryResponse{}.Definition(),
			"CreateTokenResponse":      models.CreateTokenResponse{}.Definition(),
			"ListRepositoriesResponse": models.ListRepositoriesResponse{}.Definition(),
			"Repository":               models.Repository{}.Definition(),
			"ListImagesResponse":       models.ListImagesResponse{}.Definition(),
			"Image":                    models.Image{}.Definition(),
		},
		SecurityDefinitions: map[string]swagger.SecurityDefinition{
			"login": {
				Type:        "basic",
				Description: "Basic authentication",
			},
		},
	}

	bytes, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		return nil, err
	}

	return fireball.NewResponse(200, bytes, fireball.JSONHeaders), nil
}
