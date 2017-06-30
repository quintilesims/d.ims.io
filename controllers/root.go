package controllers

import (
	"github.com/zpatrick/fireball"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

func (r *RootController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": r.GetRoot,
			},
		},
	}
}

func (r *RootController) GetRoot(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewResponse(200, nil, nil), nil
	//return fireball.Redirect(301, "/api/?url=/swagger.json"), nil
}
