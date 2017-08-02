package models

import (
	"time"

	"github.com/zpatrick/go-plugin-swagger"
)

type Image struct {
	Digest   string    `json:"digest"`
	Size     string    `json:"size"`
	PushedAt time.Time `json:"pushed_at"`
}

func (r Image) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"digest":    swagger.NewStringProperty(),
			"size":      swagger.NewStringProperty(),
			"pushed_at": swagger.NewStringProperty(),
		},
	}
}
