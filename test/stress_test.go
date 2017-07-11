package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/quintilesims/d.ims.io/test/helpers"
)

func getRandomRepo() (string, string) {
	nato := []string{
		"alpha",
		"bravo",
		"charlie",
		"delta",
		"echo",
		"foxtrot",
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

func TestSimpleWorkflow(t *testing.T) {
	api := helpers.NewTestAPIClient(t, Endpoint(true), Token())
	docker := helpers.NewTestDockerClient(t)

	name, tag := getRandomRepo()
	api.CreateRepository(TEST_REPO_OWNER, name)

	size := fmt.Sprintf("%dMB", rand.Intn(1000))
	docker.Build(tag, map[string]string{"size": size})

	t.Logf("Pushing %s image to %s", size, tag)
	docker.Push(tag)

	docker.RMI(tag)

	t.Logf("Pulling image %s", tag)
	docker.Pull(tag)
}
