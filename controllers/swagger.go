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
		//BasePath:       "/v1",
		Schemes: []string{"https"},
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
		},
		Paths: map[string]swagger.Path{
			"/token": map[string]swagger.Method{
				"post": {
					Tags:     []string{"Token"},
					Summary:  "Create a new Token",
					Security: swagger.BasicAuthSecurity("login"),
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateTokenRequest", "none", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "success",
							Schema:      swagger.NewObjectSchema("CreateTokenResponse"),
						},
					},
				},
			},
			 "/repository": map[string]swagger.Method{
                                "post": {
                                        Tags:     []string{"Repository"},
                                        Summary:  "Create a new Repository",
                                        Security: swagger.BasicAuthSecurity("login"),
                                        Parameters: []swagger.Parameter{
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
		},
		Definitions: map[string]swagger.Definition{
			"CreateRepositoryRequest": models.CreateRepositoryRequest{}.Definition(),
			"CreateRepositoryResponse": models.CreateRepositoryResponse{}.Definition(),
			"CreateTokenRequest":      models.CreateTokenRequest{}.Definition(),
			"CreateTokenResponse":     models.CreateTokenResponse{}.Definition(),
			"Repository":              models.Repository{}.Definition(),
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
