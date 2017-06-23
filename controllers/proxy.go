package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"

	"log"
	"net/http"
	"time"

	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"github.com/zpatrick/fireball"
	"github.com/zpatrick/go-cache"
)

// ecr tokens last for 12 hours: https://github.com/aws/aws-sdk-go/blob/master/service/ecr/api.go#L1022
const TOKEN_EXPIRY = time.Hour * 12

type ProxyController struct {
	ecr   ecriface.ECRAPI
	proxy proxy.Proxy
	cache *cache.Cache
}

func NewProxyController(ecr ecriface.ECRAPI, p proxy.Proxy) *ProxyController {
	return &ProxyController{
		ecr:   ecr,
		proxy: p,
		cache: cache.New(),
	}
}

func (p *ProxyController) DoProxy(c *fireball.Context) (fireball.Response, error) {
	token, err := p.getRegistryAuthToken()
	if err != nil {
		log.Printf("[ERROR] Failed to get auth token for registry: %v", err)
		return nil, err
	}

	response := fireball.ResponseFunc(func(w http.ResponseWriter, r *http.Request) {
		p.proxy.ServeHTTP(token, w, r)
	})

	return response, nil
}

func (p *ProxyController) getRegistryAuthToken() (string, error) {
	if token, ok := p.cache.Getf("token"); ok {
		return token.(string), nil
	}

	input := &ecr.GetAuthorizationTokenInput{}
	output, err := p.ecr.GetAuthorizationToken(input)
	if err != nil {
		return "", err
	}

	token := aws.StringValue(output.AuthorizationData[0].AuthorizationToken)
	p.cache.Addf("token", token, TOKEN_EXPIRY)

	return token, nil
}
