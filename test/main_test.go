package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/d.ims.io/models"
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

func TestMain(m *testing.M) {
	setup()
	fmt.Println("[INFO] Starting stress test")
	code := m.Run()
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

	fmt.Println("[INFO] Setting docker authentication")
	if err := setDockerToken(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("[INFO] Clearing test repositories")
	if err := clearTestRepos(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// setDockerToken adds the authentication for the registry into ~/.docker/config.json
func setDockerToken() error {
	config := struct {
		Auths map[string]interface{} `json:"auths"`
	}{
		Auths: map[string]interface{}{},
	}

	path := fmt.Sprintf("%s/.docker/config.json", homedir.Get())
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}
	}

	endpoint := Endpoint(false)
	if _, ok := config.Auths[endpoint]; !ok {
		config.Auths[endpoint] = map[string]interface{}{}
	}

	config.Auths[endpoint].(map[string]interface{})["auth"] = Token()

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(path, data, 0600)
}

// clearTestRepos removes all repositories owned by the TEST_REPO_OWNER
// and (re)creates the repositories 'small', 'medium', and 'large' for the same owner
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

func shell(t *testing.T, format string, tokens ...interface{}) {
	args := strings.Split(fmt.Sprintf(format, tokens...), " ")
	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		text := fmt.Sprintf("Error running %v: %v\n", cmd.Args, err)
		for _, line := range strings.Split(string(output), "\n") {
			text += line + "\n"
		}

		t.Fatalf(text)
	}
}
