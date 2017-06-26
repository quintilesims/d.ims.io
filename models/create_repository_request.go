package models

import (
	"fmt"

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

func (r CreateRepositoryRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("Field 'name' is required")
	}

	return nil
}
