package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type CreateRepositoryRequest struct {
	Name string `json:"name"`
}

func (r CreateRepositoryRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"name": swagger.NewStringProperty(),
		},
	}
}
