package main

import (
	"fmt"

	"github.com/quintilesims/d.ims.io/test/client"

	"math/rand"
	"testing"
)

func getRandomRepo() (string, string) {
	nato := []string{
		"alpha",
		"bravo",
		"charlie",
		"delta",
		"echo",
		"foxtrop",
		"golf",
		"hotel",
		"india",
		"juliet",
		"kilo",
		"lima",
		"mike",
		"november",
		"oscar",
		"papa",
		"quebec",
		"romeo",
		"sierra",
		"tango",
		"uniform",
		"victor",
		"whiskey",
		"x-ray",
		"yankee",
		"zulu",
	}

	name := nato[rand.Intn(len(nato)-1)]
	tag := fmt.Sprintf("%s/%s/%s", Endpoint(false), TEST_REPO_OWNER, name)
	return name, tag
}

// todo: stres test  pull
func TestStressPush(t *testing.T) {
	api := client.NewTestAPIClient(t, Endpoint(true), Token())
	docker := client.NewTestDockerClient(t)

	for i := 0; i < 5; i++ {
	name, tag := getRandomRepo()
	api.CreateRepository(TEST_REPO_OWNER, name)

	size := fmt.Sprintf("%dMB", rand.Intn(1000))
	docker.Build(tag, map[string]string{"size": size})

	t.Logf("Pushing %s image to %s", size, tag)
	docker.Push(tag)
	}

}

