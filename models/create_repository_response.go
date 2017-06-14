package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type CreateRepositoryResponse struct {
	Name string `json:"name"`
}

func (r CreateRepositoryResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"name": swagger.NewStringProperty(),
		},
	}
}
