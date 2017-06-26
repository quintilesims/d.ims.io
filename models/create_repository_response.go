package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type CreateRepositoryResponse struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (r CreateRepositoryResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"owner": swagger.NewStringProperty(),
			"name":  swagger.NewStringProperty(),
		},
	}
}
