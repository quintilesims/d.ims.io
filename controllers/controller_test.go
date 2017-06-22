package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/zpatrick/fireball"
	"net/http"
	"testing"
)

func generateContext(t *testing.T, v interface{}, pathVariables map[string]string) *fireball.Context {
	context := &fireball.Context{
		PathVariables: pathVariables,
	}

	if v != nil {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(v); err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest("", "", b)
		if err != nil {
			t.Fatal(err)
		}

		context.Request = request
	}

	return context
}
