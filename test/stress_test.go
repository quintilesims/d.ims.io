package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

type RunInfo struct {
	RepositoryName string
	Dockerfile     string
	Iterations     int
}

func TestStressSmall(t *testing.T) {
	doStressTest(t, RunInfo{"small", "Dockerfile.small", 1})
}

func doStressTest(t *testing.T, r RunInfo) {
	t.Parallel()

	for i := 0; i < r.Iterations; i++ {
		tag := fmt.Sprintf("%s/%s/%s", Endpoint(false), TEST_REPO_OWNER, r.RepositoryName)
		shell(t, "docker build --tag %s --file %s .", tag, r.Dockerfile)
		shell(t, "docker push %s", tag)
	}
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
