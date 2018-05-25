package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/fireball"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newContextWithBasicAuth(t *testing.T, user, pass string) *fireball.Context {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth(user, pass)
	return &fireball.Context{
		Request: req,
	}
}

func assertResponseCode(t *testing.T, resp fireball.Response, code int) {
	recorder := httptest.NewRecorder()
	resp.Write(recorder, nil)
	assert.Equal(t, code, recorder.Code)
}
