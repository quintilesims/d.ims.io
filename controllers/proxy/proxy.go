package proxy

import (
	"net/http"
)

type Proxy interface {
	ServeHTTP(*http.Request, http.ResponseWriter)
}
