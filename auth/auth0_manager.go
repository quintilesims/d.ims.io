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
	clientID           string
	connection         string
	client             *rclient.RestClient
	cache              *cache.Cache
	throttle           <-chan time.Time
	throttleMin        int
	throttleMultiplier int
	lastThrottled      time.Time
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

func NewAuth0Manager(domain, clientID, connection string, throttleMin int) *Auth0Manager {
	return &Auth0Manager{
		clientID:           clientID,
		connection:         connection,
		client:             rclient.NewRestClient(domain),
		cache:              cache.New(),
		throttle:           time.Tick(time.Millisecond * time.Duration(throttleMin)),
		throttleMin:        throttleMin,
		throttleMultiplier: throttleMin,
		lastThrottled:      time.Now(),
	}
}

func (a *Auth0Manager) Authenticate(username, password string) (bool, error) {
	if a.throttleMultiplier > a.throttleMin && time.Since(a.lastThrottled) > 1*time.Minute {
		a.throttleMultiplier = a.throttleMultiplier / 2
		if a.throttleMultiplier < a.throttleMin {
			a.throttleMultiplier = a.throttleMin
		}

		log.Printf("[INFO] 1 min since last throttle update, decreasing throttle to %v", time.Millisecond*time.Duration(a.throttleMultiplier))

		a.throttle = time.Tick(time.Millisecond / time.Duration(a.throttleMultiplier))
		a.lastThrottled = time.Now()
	}

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

	// throttle to deal with rate limiting
	<-a.throttle

	// will only sleep if cachedStatus already exists with a penalty
	time.Sleep(cachedStatus.penalty * time.Duration(timeMultiplier))

	req := oauthReq{
		ClientID:   a.clientID,
		Connection: a.connection,
		Username:   username,
		Password:   password,
		GrantType:  "password",
		Scope:      "openid",
	}

	if err := a.client.Post("/oauth/ro", req, nil); err != nil {
		if err, ok := err.(*rclient.ResponseError); ok && (err.Response.StatusCode == 401 || err.Response.StatusCode == 429) {
			cachedStatus.penalty += time.Second
			if cachedStatus.penalty > maxPenalty {
				cachedStatus.penalty = maxPenalty
			}

			if err.Response.StatusCode == 429 {
				if a.throttleMultiplier > 0 {
					a.throttleMultiplier = a.throttleMultiplier * 2
				} else {
					a.throttleMultiplier = 1
				}

				log.Printf("[INFO] Too many requests, increasing throttle to %v", time.Millisecond*time.Duration(a.throttleMultiplier))

				a.throttle = time.Tick(time.Millisecond * time.Duration(a.throttleMultiplier))
				a.lastThrottled = time.Now()
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
