package router

import (
	"net/http"

	"github.com/zpatrick/fireball"
)

func NewRouter(routes []*fireball.Route, doProxy fireball.Handler) fireball.RouterFunc {
	router := fireball.NewBasicRouter(routes)
	return fireball.RouterFunc(func(req *http.Request) (*fireball.RouteMatch, error) {
		match, err := router.Match(req)
		if match != nil || err != nil {
			return match, err
		}

		// perform reverse proxy on any request that doesn't match the d.ims.io api
		// this could be improved by verifying the request is part of the registry api
		// see: https://docs.docker.com/registry/spec/api
		return &fireball.RouteMatch{Handler: doProxy}, nil
	})
}
