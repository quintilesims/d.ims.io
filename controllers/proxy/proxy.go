package proxy

import (
	"net/http"
)

type Proxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
