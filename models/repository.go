package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type Repository struct {
	Owner     string   `json:"owner"`
	Name      string   `json:"name"`
	ImageTags []string `json:"image_tags"`
}

func (r Repository) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"owner":      swagger.NewStringProperty(),
			"name":       swagger.NewStringProperty(),
			"image_tags": swagger.NewStringSliceProperty(),
		},
	}
}
