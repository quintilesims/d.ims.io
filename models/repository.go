package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type Repository struct {
	Name      string   `json:"name"`
	URI       string   `json:"uri"`
	ImageTags []string `json:"image_tags"`
}

func (r Repository) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"name":       swagger.NewStringProperty(),
			"uri":        swagger.NewStringProperty(),
			"image_tags": swagger.NewStringSliceProperty(),
		},
	}
}
