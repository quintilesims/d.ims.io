package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ECRProxy struct {
	endpoint     string
	reverseProxy *httputil.ReverseProxy
}

func NewECRProxy(endpoint string) *ECRProxy {
	return &ECRProxy{
		endpoint: endpoint,
		reverseProxy: httputil.NewSingleHostReverseProxy(&url.URL{
			Host:   endpoint,
			Scheme: "https",
		}),
	}
}

func (e *ECRProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Host = e.endpoint
	e.reverseProxy.ServeHTTP(w, r)
}
