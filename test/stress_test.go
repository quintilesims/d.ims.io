package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/quintilesims/d.ims.io/test/helpers"
)

func newRepoGenerator() func() (string, string, string) {
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

	var i int
	return func() (string, string, string) {
		name := nato[i%len(nato)]
		tag := fmt.Sprintf("%s/%s/%s", Endpoint(false), TEST_REPO_OWNER, name)
		size := fmt.Sprintf("%dMB", rand.Intn(1000))

		i++
		return name, tag, size
	}
}

func TestStress(t *testing.T) {
	repoGenerator := newRepoGenerator()

	for i := 0; i < Count(); i++ {
		name, tag, size := repoGenerator()

		t.Run(name, func(tt *testing.T) {
			tt.Parallel()

			api := helpers.NewTestAPIClient(tt, Endpoint(true), Token())
			docker := helpers.NewTestDockerClient(tt)

			tt.Logf("Creating repository %s", name)
			api.CreateRepository(TEST_REPO_OWNER, name)

			tt.Logf("Building image %s", tag)
			docker.Build(tag, map[string]string{"size": size})

			tt.Logf("Pushing %s image to %s", size, tag)
			docker.Push(tag)

			tt.Logf("Pulling image %s", tag)
			docker.RMI(tag)
			docker.Pull(tag)

			tt.Logf("Deleting repository %s", name)
			api.DeleteRepository(TEST_REPO_OWNER, name)
			docker.RMI(tag)
		})
	}
}
