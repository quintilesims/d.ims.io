package auth

import (
	"log"
	"time"

	"github.com/zpatrick/rclient"
)

type Auth0Authenticator struct {
	clientID   string
	connection string
	client     *rclient.RestClient
	throttle   <-chan time.Time
}

type oauthReq struct {
	ClientID   string `json:"client_id"`
	Connection string `json:"connection"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	GrantType  string `json:"grant_type"`
	Scope      string `json:"scope"`
}

func NewAuth0Authenticator(domain, clientID, connection string, rateLimit time.Duration) *Auth0Authenticator {
	return &Auth0Authenticator{
		clientID:   clientID,
		connection: connection,
		client:     rclient.NewRestClient(domain),
		throttle:   time.Tick(rateLimit),
	}
}

func (a *Auth0Authenticator) Authenticate(username, password string) (bool, error) {
	log.Printf("[DEBUG] Attempting to authenticate user '%s' through Auth0", username)

	req := oauthReq{
		ClientID:   a.clientID,
		Connection: a.connection,
		Username:   username,
		Password:   password,
		GrantType:  "password",
		Scope:      "openid",
	}

	<-a.throttle
	if err := a.client.Post("/oauth/ro", req, nil); err != nil {
		if re, ok := err.(*rclient.ResponseError); ok && re.Response.StatusCode == 401 {
			log.Printf("[DEBUG] User '%s' sent invalid Auth0 credentials", username)
			return false, nil
		}

		return false, err
	}

	log.Printf("[DEBUG] User '%s' sent valid Auth0 credentials", username)
	return true, nil
}
