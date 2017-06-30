package main

import (
	"fmt"
	"testing"
)

//func TestStressSmall(t *testing.T)  { doStressTest(t, "small", "Dockerfile.small", 1) }
//func TestStressMedium(t *testing.T) { doStressTest(t, "medium", "Dockerfile.medium", 1) }
func TestStressLarge(t *testing.T) { doStressTest(t, "large", "Dockerfile.large", 1) }

func doStressTest(t *testing.T, repositoryName, dockerfile string, iterations int) {
	t.Parallel()

	for i := 0; i < iterations; i++ {
		tag := fmt.Sprintf("%s/%s/%s", Endpoint(false), TEST_REPO_OWNER, repositoryName)
		shell(t, "docker build --tag %s --file %s .", tag, dockerfile)
		shell(t, "docker push %s", tag)
	}
}
