package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/quintilesims/d.ims.io/auth"
	"github.com/zpatrick/fireball"
	cache "github.com/zpatrick/go-cache"
)

const (
	validAuthExpiry = time.Hour
)

func hash(user, pass string) string {
	sum := sha256.Sum256([]byte(user + pass))
	return fmt.Sprintf("%x", sum)
}

func AuthDecorator(auth auth.Authenticator) fireball.Decorator {
	che := cache.New()
	return func(handler fireball.Handler) fireball.Handler {
		return func(c *fireball.Context) (fireball.Response, error) {
			headers := map[string]string{"WWW-Authenticate": "Basic realm=\"Restricted\""}
			invalidAuthResponse := fireball.NewResponse(401, []byte("401 Unauthorized\n"), headers)

			user, pass, ok := c.Request.BasicAuth()
			if !ok {
				log.Printf("[DEBUG] Request %s %s did not contain basic auth", c.Request.Method, c.Request.URL.String())
				return invalidAuthResponse, nil
			}

			log.Printf("[DEBUG] Attempting to authenticate user '%s'", user)

			key := hash(user, pass)
			isValidCreds, ok := che.GetOK(key)
			if ok && isValidCreds.(bool) {
				log.Printf("[DEBUG] Allowing valid cached creds for user '%s'", user)
				return handler(c)
			}

			if ok && !isValidCreds.(bool) {
				log.Printf("[DEBUG] Denying invalid cached creds for user '%s'", user)
				return invalidAuthResponse, nil
			}

			isAuthenticated, err := auth.Authenticate(user, pass)
			if err != nil {
				log.Printf("[ERROR] Authenticator encountered an unexpected error: %v", err)
				return nil, err
			}

			if !isAuthenticated {
				log.Printf("[DEBUG] User '%s' failed to authenticate", user)
				che.Set(key, false)
				return invalidAuthResponse, nil
			}

			log.Printf("[DEBUG] User '%s' successfully authenticated", user)
			che.Set(key, true, cache.Expire(validAuthExpiry))
			return handler(c)
		}
	}
}
