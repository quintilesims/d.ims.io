package helpers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/rclient"
)

type TestAPIClient struct {
	T      *testing.T
	Client *rclient.RestClient
}

func NewTestAPIClient(t *testing.T, endpoint, token string) *TestAPIClient {
	doer := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	addAuth := rclient.Header("Authorization", fmt.Sprintf("Basic %s", token))
	client := rclient.NewRestClient(endpoint, rclient.Doer(doer),
		rclient.RequestOptions(addAuth))

	return &TestAPIClient{
		T:      t,
		Client: client,
	}
}

func (a *TestAPIClient) CreateRepository(owner, name string) {
	req := models.CreateRepositoryRequest{Name: name}
	path := fmt.Sprintf("/repository/%s", owner)
	if err := a.Client.Post(path, req, nil); err != nil {
		a.T.Fatal(err)
	}
}
