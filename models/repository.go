package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type Repository struct {
	Name string `json:"name"`
}

func (r Repository) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"name": swagger.NewStringProperty(),
		},
	}
}
