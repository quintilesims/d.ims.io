package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/quintilesims/layer0/setup/docker"
	"github.com/zpatrick/rclient"
)

const (
	ENVVAR_ENDPOINT = "DIMSIO_TEST_ENDPOINT"
	ENVVAR_TOKEN    = "DIMSIO_TEST_TOKEN"
	TEST_REPO_OWNER = "test"
)

func Endpoint(showProtocol bool) string {
	endpoint := os.Getenv(ENVVAR_ENDPOINT)
	if !showProtocol {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	return endpoint
}

func Token() string {
	return os.Getenv(ENVVAR_TOKEN)
}

func RepositoryNames() []string {
	return []string{"small", "medium", "large"}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	if Endpoint(true) == "" {
		fmt.Printf("Required environment variable %s not set\n", ENVVAR_ENDPOINT)
		os.Exit(1)
	}

	if Token() == "" {
		fmt.Printf("Required environment variable %s not set\n", ENVVAR_TOKEN)
		os.Exit(1)
	}

	if err := setDockerToken(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := clearTestRepos(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func setDockerToken() error {
	path := fmt.Sprintf("%s/.docker/config.json", homedir.Get())
	config, err := docker.LoadConfig(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if config == nil {
		config = &docker.Config{Auths: map[string]docker.Auth{}}
	}

	endpoint := Endpoint(false)
	if _, ok := config.Auths[endpoint]; !ok {
		config.Auths[endpoint] = map[string]interface{}{}
	}

	config.Auths[endpoint]["auth"] = Token()
	return docker.WriteConfig(path, config)
}

func clearTestRepos() error {
	doer := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	addAuth := rclient.Header("Authorization", fmt.Sprintf("Basic %s", Token()))

	client := rclient.NewRestClient(Endpoint(true),
		rclient.Doer(doer),
		rclient.RequestOptions(addAuth))

	var resp models.ListRepositoriesResponse
	path := fmt.Sprintf("/repository/%s", TEST_REPO_OWNER)
	if err := client.Get(path, &resp); err != nil {
		return err
	}

	for _, name := range resp.Repositories {
		path := fmt.Sprintf("/repository/%s/%s", TEST_REPO_OWNER, name)
		if err := client.Delete(path, nil, nil); err != nil {
			return err
		}
	}

	for _, name := range []string{"small", "medium", "large"} {
		req := models.CreateRepositoryRequest{Name: name}
		path := fmt.Sprintf("/repository/%s", TEST_REPO_OWNER)
		if err := client.Post(path, req, nil); err != nil {
			return err
		}
	}

	return nil
}

func teardown() {
}
