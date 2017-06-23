package auth

import (
	"fmt"
	"time"

	"github.com/zpatrick/go-cache"
	"github.com/zpatrick/rclient"
)

type Auth0Manager struct {
	clientID   string
	connection string
	client     *rclient.RestClient
	cache      *cache.Cache
}

type Auth0Config struct {
	Domain     string
	ClientID   string
	Connection string
}

type oauthReq struct {
	ClientID   string `json:"client_id"`
	Connection string `json:"connection"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	GrantType  string `json:"grant_type"`
	Scope      string `json:"scope"`
}

type authStatus struct {
	isValid bool
	penalty time.Duration
}

const (
	maxPenalty      = 5 * time.Second
	validAuthExpiry = 1 * time.Hour
)

func NewAuth0Manager(config Auth0Config) *Auth0Manager {
	return &Auth0Manager{
		clientID:   config.ClientID,
		connection: config.Connection,
		client:     rclient.NewRestClient(config.Domain),
		cache:      cache.New(),
	}
}

func (a *Auth0Manager) Authenticate(username, password string) (bool, error) {
	key := fmt.Sprintf("%s:%s", username, password)
	var cachedStatus *authStatus
	if result, exists := a.cache.Getf(key); exists {
		cachedStatus = result.(*authStatus)
	} else {
		cachedStatus = &authStatus{}
	}

	if cachedStatus.isValid {
		return true, nil
	}

	// will only sleep if cachedStatus already exists with a penalty
	time.Sleep(cachedStatus.penalty)

	req := oauthReq{
		ClientID:   a.clientID,
		Connection: a.connection,
		Username:   username,
		Password:   password,
		GrantType:  "password",
		Scope:      "openid",
	}

	if err := a.client.Post("/oauth/ro", req, nil); err != nil {
		if err, ok := err.(*rclient.ResponseError); ok && err.Response.StatusCode == 401 {
			cachedStatus.penalty += time.Second
			if cachedStatus.penalty > maxPenalty {
				cachedStatus.penalty = maxPenalty
			}

			a.cache.Add(key, cachedStatus)
			return false, nil
		}

		return false, err
	}

	cachedStatus.isValid = true
	cachedStatus.penalty = 0 * time.Second
	a.cache.Addf(key, cachedStatus, validAuthExpiry)
	return true, nil
}
