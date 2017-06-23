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

		return &fireball.RouteMatch{Handler: doProxy}, nil
	})
}
