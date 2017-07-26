package models

import (
	"time"

	"github.com/zpatrick/go-plugin-swagger"
)

type ListImagesResponse struct {
	Images []Image `json:"images"`
}

type Image struct {
	Repository    string    `json:"repository"`
	ImageTags     []string  `json:"image_tags,omitempty"`
	ImageDigest   string    `json:"image_digest"`
	ImagePushedAt time.Time `json:"image_pushed_at"`
}

func (r ListImagesResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"images": swagger.NewObjectSliceProperty("Image"),
		},
	}
}

func (r Image) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"repository":      swagger.NewStringProperty(),
			"image_tags":      swagger.NewStringSliceProperty(),
			"image_digest":    swagger.NewStringProperty(),
			"image_pushed_at": swagger.NewStringProperty(),
		},
	}
}
