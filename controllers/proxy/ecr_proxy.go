package proxy

import (
	"net/http"
)

type ECRProxy struct {
	endpoint string
}

func NewECRProxy(endpoint string) *ECRProxy {
	return &ECRProxy{
		endpoint: endpoint,
	}
}

func (e *ECRProxy) ServeHTTP(r *http.Request, w http.ResponseWriter) {
	http.Error(w, "ecr proxy not implemented", 500)
}
