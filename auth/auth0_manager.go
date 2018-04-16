package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/zpatrick/go-cache"
	"github.com/zpatrick/rclient"
)

var timeMultiplier = 1

type Auth0Manager struct {
	clientID   string
	connection string
	client     *rclient.RestClient
	cache      *cache.Cache
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

type authStatus struct {
	isValid bool
	penalty time.Duration
}

const (
	maxPenalty      = 5 * time.Second
	validAuthExpiry = 1 * time.Hour
)

func NewAuth0Manager(domain, clientID, connection string, rateLimit time.Duration) *Auth0Manager {
	return &Auth0Manager{
		clientID:   clientID,
		connection: connection,
		client:     rclient.NewRestClient(domain),
		cache:      cache.New(),
		throttle:   time.Tick(rateLimit),
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
	time.Sleep(cachedStatus.penalty * time.Duration(timeMultiplier))

	isAuthenticated, err := a.authenticate(username, password)
	if err != nil {
		return false, err
	}

	if !isAuthenticated {
		cachedStatus.penalty += time.Second
		if cachedStatus.penalty > maxPenalty {
			cachedStatus.penalty = maxPenalty
		}

		a.cache.Add(key, cachedStatus)
		return false, nil
	}

	cachedStatus.isValid = true
	cachedStatus.penalty = 0 * time.Second
	a.cache.Addf(key, cachedStatus, validAuthExpiry)

	return true, nil
}

func (a *Auth0Manager) authenticate(username, password string) (bool, error) {
	log.Printf("[DEBUG] Attempting to authenticate user %s through Auth0", username)
	req := oauthReq{
		ClientID:   a.clientID,
		Connection: a.connection,
		Username:   username,
		Password:   password,
		GrantType:  "password",
		Scope:      "openid",
	}

	for backoff := time.Duration(0); true; backoff += (time.Millisecond * 500) {
		time.Sleep(backoff * time.Duration(timeMultiplier))
		<-a.throttle

		if err := a.client.Post("/oauth/ro", req, nil); err != nil {
			if err, ok := err.(*rclient.ResponseError); ok {
				switch err.Response.StatusCode {
				case 401:
					log.Printf("[DEBUG] User '%s' sent invalid Auth0 credentials.", username)
					return false, nil
				case 429:
					log.Printf("[DEBUG] Auth0 returned 429 response for user '%s'; retrying.", username)
					continue
				default:
					return false, err
				}
			}

			return false, err
		}

		log.Printf("[DEBUG] User '%s' sent valid Auth0 credentials.", username)
		break
	}

	return true, nil
}
