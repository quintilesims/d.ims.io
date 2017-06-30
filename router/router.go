package router

import (
	"net/http"
	"log"

	"github.com/zpatrick/fireball"
)

func NewRouter(routes []*fireball.Route, doProxy fireball.Handler) fireball.RouterFunc {
	router := fireball.NewBasicRouter(routes)

	return fireball.RouterFunc(func(req *http.Request) (*fireball.RouteMatch, error) {
		 log.Printf("HGOt requrest %s %s", req.Method, req.URL.String())
		

		match, err := router.Match(req)
		if match != nil || err != nil {
			return match, err
		}

		/*
		if req.URL.String() == "/v2/" && (req.Method == "GET" || req.Method == "HEAD") {
			handler := fireball.Handler(func(*fireball.Context) (fireball.Response, error) {
				headers := map[string]string{
					"Docker-Distribution-Api-Version": "registry/2.0",
					"WWW-Authenticate":                "Basic realm=\"https://docker.ims.io/\"service=\"ecr.amazonaws.com\"",
				}


				fmt.Println("!!!!!!!!!!v2 hit!")
				return fireball.NewResponse(401, nil, headers), nil
			})

			return &fireball.RouteMatch{Handler: handler}, nil
		}

		fmt.Println(req.URL.String())
		*/

		// perform reverse proxy on any request that doesn't match the d.ims.io api
		// this could be improved by verifying the request is part of the registry api
		// see: https://docs.docker.com/registry/spec/api
		return &fireball.RouteMatch{Handler: doProxy}, nil
	})
}
