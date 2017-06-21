package controllers

import (
	"fmt"
	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"github.com/zpatrick/fireball"
	//"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
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

func (r *ProxyController) DoProxy(c *fireball.Context) (fireball.Response, error) {
	return nil, fmt.Errorf("proxy controller implemented")
}
