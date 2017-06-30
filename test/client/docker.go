package client

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

type TestDockerClient struct {
	T *testing.T
}

func NewTestDockerClient(t *testing.T) *TestDockerClient {
	return &TestDockerClient{
		T: t,
	}
}

func (d *TestDockerClient) Build(tag string, buildArgs map[string]string) {
	command := fmt.Sprintf("docker build --tag %s ", tag)
	if len(buildArgs) > 0 {
		for k, v := range buildArgs {
			command += fmt.Sprintf("--build-arg %s=%s ", k, v)
		}
	}

	shell(d.T, command)
}

func (d *TestDockerClient) Push(tag string) {
	command := fmt.Sprintf("docker push %s", tag)
	shell(d.T, command)
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
