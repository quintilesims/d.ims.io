package controllers

import (
	"github.com/zpatrick/fireball"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

func (a *RootController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": a.GetRoot,
			},
		},
	}
}

func (a *RootController) GetRoot(c *fireball.Context) (fireball.Response, error) {
	return fireball.Redirect(301, "/api/?url=/swagger.json"), nil
}
