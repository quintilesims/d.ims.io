package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewECRProxy(registryEndpoint string) ProxyFunc {
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   registryEndpoint,
		Scheme: "https",
	})

	return ProxyFunc(func(token string, w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
		r.Host = registryEndpoint

		log.Printf("[DEBUG] Performing reverse proxy for %s %s", r.Method, r.URL.String())
		reverseProxy.ServeHTTP(w, r)
	})
}
