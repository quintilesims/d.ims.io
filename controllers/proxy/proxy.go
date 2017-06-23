package proxy

import (
	"net/http"
)

type Proxy interface {
	ServeHTTP(token string, w http.ResponseWriter, r *http.Request)
}

type ProxyFunc func(token string, w http.ResponseWriter, r *http.Request)

func (p ProxyFunc) ServeHTTP(token string, w http.ResponseWriter, r *http.Request) {
	p(token, w, r)
}
