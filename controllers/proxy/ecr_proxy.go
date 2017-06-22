package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewECRProxy(registryEndpoint string) ProxyFunc {
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   registryEndpoint,
		Scheme: "https",
	})

	return ProxyFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = registryEndpoint
		reverseProxy.ServeHTTP(w, r)
	})
}
