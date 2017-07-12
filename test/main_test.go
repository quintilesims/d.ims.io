package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/quintilesims/d.ims.io/models"
	"github.com/quintilesims/d.ims.io/test/helpers"
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

func Count() int {
	p := flag.Lookup("test.parallel")
	count, err := strconv.Atoi(p.Value.String())
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse -parallel flag: %v", err)
		os.Exit(1)
	}

	return count
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())

	setup()
	fmt.Println("[INFO] Starting stress test")
	code := m.Run()
	os.Exit(code)
}

func setup() {
	flag.Parse()

	if Endpoint(true) == "" {
		fmt.Printf("Required environment variable %s not set. Skipping test.\n", ENVVAR_ENDPOINT)
		os.Exit(0)
	}

	if Token() == "" {
		fmt.Printf("Required environment variable %s not set. Skipping test.\n", ENVVAR_TOKEN)
		os.Exit(0)
	}

	fmt.Println("[INFO] Clearing test repositories")
	if err := clearTestRepos(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func clearTestRepos() error {
	api := helpers.NewTestAPIClient(nil, Endpoint(true), Token())

	var resp models.ListRepositoriesResponse
	path := fmt.Sprintf("/repository/%s", TEST_REPO_OWNER)
	if err := api.Client.Get(path, &resp); err != nil {
		return err
	}

	for _, name := range resp.Repositories {
		path := fmt.Sprintf("/repository/%s/%s", TEST_REPO_OWNER, name)
		if err := api.Client.Delete(path, nil, nil); err != nil {
			return err
		}
	}

	return nil
}
