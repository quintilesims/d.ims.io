package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type ListRepositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

func (r ListRepositoriesResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"repositories": swagger.NewStringSliceProperty(),
		},
	}
}
