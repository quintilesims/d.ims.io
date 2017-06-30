package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewECRProxy(registryEndpoint string) ProxyFunc {
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   registryEndpoint,
		Scheme: "https",
	})

	return ProxyFunc(func(token string, w http.ResponseWriter, r *http.Request) {
		originalHost := r.Host
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
		r.Host = registryEndpoint

		reverseProxy.ModifyResponse = func(resp *http.Response) error {
			location := resp.Header.Get("Location")
			location = strings.Replace(location, registryEndpoint, originalHost, 1)
			resp.Header.Set("Location", location)

			return nil
		}

		log.Printf("[DEBUG] Performing reverse proxy for %s %s", r.Method, r.URL.String())
		reverseProxy.ServeHTTP(w, r)
	})
}
