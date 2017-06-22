package controllers

import (
	"fmt"
	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"github.com/zpatrick/fireball"
	//"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"net/http"
)

type ProxyController struct {
	ecr   ecriface.ECRAPI
	proxy proxy.Proxy
}

func NewProxyController(ecr ecriface.ECRAPI, p proxy.Proxy) *ProxyController {
	return &ProxyController{
		ecr:   ecr,
		proxy: p,
	}
}

func (p *ProxyController) DoProxy(c *fireball.Context) (fireball.Response, error) {

	token, err := p.getRegistryAuthToken()
	if err != nil {
		return nil, err
	}

	response := fireball.ResponseFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
		p.proxy.ServeHTTP(w, r)
	})

	return response, nil
}

func (p *ProxyController) getRegistryAuthToken() (string, error) {
	return "", fmt.Errorf("getRegistryAuthToken not implemented")
}
