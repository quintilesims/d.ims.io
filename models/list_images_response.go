package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type ListImagesResponse struct {
	Images []string `json:"images"`
}

func (r ListImagesResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"images": swagger.NewStringSliceProperty(),
		},
	}
}
