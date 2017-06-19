package auth

import (
	"fmt"
	"github.com/quintilesims/d.ims.io/auth/auth0"
	"github.com/quintilesims/d.ims.io/auth/token"
	"github.com/zpatrick/go-cache"
	"net/http"
	"time"
)

const (
	CACHE_EXPIRY         = time.Minute * 5
	INVALID_AUTH_PENALTY = time.Second * 5
)

type CompositeAuthenticator struct {
	Auth0 auth0.ADManager
	Token token.TokenManager
	cache *cache.Cache
}

func NewCompositeAuthenticator(a auth0.ADManager, t token.TokenManager) *CompositeAuthenticator {
	return &CompositeAuthenticator{
		Auth0: a,
		Token: t,
		cache: cache.New(),
	}
}

func (c *CompositeAuthenticator) Authenticate(r *http.Request) (bool, error) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false, nil
	}

	key := fmt.Sprintf("%s:%s", user, pass)
	isValid, exists := c.cache.Getf(key)
	if exists && isValid.(bool) {
		return true, nil
	}

	// punish the user for sending bad auth in subsequent requests
	if exists && !isValid.(bool) {
		time.Sleep(INVALID_AUTH_PENALTY)
	}

	isValidToken, err := c.Token.Authenticate(user, pass)
	if err != nil {
		return false, err
	}

	if isValidToken {
		c.cache.Addf(key, true, CACHE_EXPIRY)
		return true, nil
	}

	isValidAD, err := c.Auth0.Authenticate(user, pass)
	if err != nil {
		return false, err
	}

	if isValidAD {
		c.cache.Addf(key, true, CACHE_EXPIRY)
		return true, nil
	}

	c.cache.Addf(key, false, CACHE_EXPIRY)
	return false, nil
}
