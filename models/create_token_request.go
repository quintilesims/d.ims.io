package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type CreateTokenRequest struct{}

func (r CreateTokenRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type:       "object",
		Properties: map[string]swagger.Property{
		//"name": swagger.NewStringProperty(),
		},
	}
}
