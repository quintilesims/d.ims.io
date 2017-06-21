package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (r CreateTokenResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"token": swagger.NewStringProperty(),
		},
	}
}
