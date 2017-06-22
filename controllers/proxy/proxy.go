package proxy

import (
	"net/http"
)

type Proxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ProxyFunc func(w http.ResponseWriter, r *http.Request)

func (p ProxyFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p(w, r)
}
