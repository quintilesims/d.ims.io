package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type Handler func(w http.ResponseWriter, r *http.Request)

func newAuth0AuthenticatorAndServer(handler Handler) (*Auth0Authenticator, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	auth0Manager := NewAuth0Authenticator(server.URL, "", "", time.Second/2)

	return auth0Manager, server
}

func MarshalAndWrite(t *testing.T, w http.ResponseWriter, body interface{}, status int) {
	MarshalAndWriteHeader(t, w, body, nil, status)
}

func MarshalAndWriteHeader(t *testing.T, w http.ResponseWriter, body interface{}, headers map[string]string, status int) {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	for key, val := range headers {
		w.Header().Set(key, val)
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, string(b))
}

func Unmarshal(t *testing.T, r *http.Request, content interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(body, &content); err != nil {
		t.Fatal(err)
	}
}
