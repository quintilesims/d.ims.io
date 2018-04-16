package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type GrantAccessRequest struct {
	Account string `json:"account"`
}

func (r GrantAccessRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"account": swagger.NewStringProperty(),
		},
	}
}
