package models

import (
	"time"

	"github.com/zpatrick/go-plugin-swagger"
)

type Repository struct {
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Repository) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"owner":      swagger.NewStringProperty(),
			"name":       swagger.NewStringProperty(),
			"created_at": swagger.NewStringProperty(),
		},
	}
}
