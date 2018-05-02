package models

import (
	"github.com/zpatrick/go-plugin-swagger"
)

type ListAccountsResponse struct {
	Accounts []string `json:"accounts"`
}

func (a ListAccountsResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"accounts": swagger.NewStringSliceProperty(),
		},
	}
}
