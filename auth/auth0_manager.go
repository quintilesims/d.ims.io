package auth

import (
	"github.com/zpatrick/rclient"
	"log"
)

type Auth0Manager struct {
	token  string
	client *rclient.RestClient
}

func NewAuth0Manager(endpoint, token string) *Auth0Manager {
	return &Auth0Manager{
		client: rclient.NewRestClient(endpoint),
	}
}

func (a *Auth0Manager) Authenticate(user, pass string) (bool, error) {
	log.Println("[ERROR] - Auth0Manager.Authenticate not implemented")
	return true, nil
}
