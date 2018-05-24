package models

import (
	"fmt"

	"github.com/zpatrick/go-plugin-swagger"
)

type GrantAccessRequest struct {
	Account string `json:"account"`
}

func (r GrantAccessRequest) Validate() error {
	if r.Account == "" {
		return fmt.Errorf("account is a required field")
	}

	return nil
}

func (r GrantAccessRequest) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"account": swagger.NewStringProperty(),
		},
	}
}
