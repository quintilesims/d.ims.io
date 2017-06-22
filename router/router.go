package router

import (
	"github.com/zpatrick/fireball"
	"net/http"
)

type Router struct {
	router  *fireball.BasicRouter
	doProxy fireball.Handler
}

func NewRouter(routes []*fireball.Route, doProxy fireball.Handler) *Router {
	return &Router{
		router:  fireball.NewBasicRouter(routes),
		doProxy: doProxy,
	}
}

func (r *Router) Match(req *http.Request) (*fireball.RouteMatch, error) {
	match, err := r.router.Match(req)
	if match != nil || err != nil {
		return match, err
	}

	return &fireball.RouteMatch{Handler: r.doProxy}, nil
}
