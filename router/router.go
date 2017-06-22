package router

import (
	"net/http"

	"github.com/quintilesims/d.ims.io/controllers"
	"github.com/zpatrick/fireball"
)

type Router struct {
	routes          []*fireball.Route
	proxyController *controllers.ProxyController
	router          *fireball.BasicRouter
}

func NewRouter(routes []*fireball.Route, p *controllers.ProxyController) *Router {
	return &Router{
		routes:          routes,
		proxyController: p,
		router:          fireball.NewBasicRouter(routes),
	}
}

func (r *Router) Match(req *http.Request) (*fireball.RouteMatch, error) {
	return r.router.Match(req)
}
