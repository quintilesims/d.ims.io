package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth0ManagerAuthenticate_ValidCreds(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/oauth/ro")

		var req oauthReq
		Unmarshal(t, r, &req)

		assert.Equal(t, req.Username, "valid username")
		assert.Equal(t, req.Password, "valid password")

		MarshalAndWrite(t, w, nil, 200)
	}

	auth0Manager, server := newAuth0ManagerAndServer(handler)
	defer server.Close()

	valid, err := auth0Manager.Authenticate("valid username", "valid password")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, valid, true)
}

func TestAuth0ManagerAuthenticate_InvalidCreds(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/oauth/ro")

		var req oauthReq
		Unmarshal(t, r, &req)

		assert.Equal(t, req.Username, "valid username")
		assert.Equal(t, req.Password, "invalid password")

		MarshalAndWrite(t, w, nil, 401)
	}

	auth0Manager, server := newAuth0ManagerAndServer(handler)
	defer server.Close()

	valid, err := auth0Manager.Authenticate("valid username", "invalid password")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, valid, false)
}

func TestAuth0ManagerAuthenticate_ValidCredsAreCached(t *testing.T) {
	count := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		MarshalAndWrite(t, w, nil, 200)
		count++
	}

	auth0Manager, server := newAuth0ManagerAndServer(handler)
	defer server.Close()

	for i := 0; i < 2; i++ {
		valid, err := auth0Manager.Authenticate("valid username", "valid password")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, valid, true)
	}

	assert.Equal(t, count, 1)
}

func TestAuth0ManagerAuthenticate_NotValidCredsAreChecked(t *testing.T) {
	timeMultiplier = 0
	count := 0

	handler := func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			MarshalAndWrite(t, w, nil, 401)
		} else {
			MarshalAndWrite(t, w, nil, 200)
		}

		count++
	}

	auth0Manager, server := newAuth0ManagerAndServer(handler)
	defer server.Close()

	for i, expected := range []bool{false, true} {
		valid, err := auth0Manager.Authenticate("username", "password")
		if err != nil {
			t.Fatal(err)
		}

		if valid != expected {
			t.Errorf("Error on iteration %d: result was %t, expected %t", i, valid, expected)
		}
	}

	assert.Equal(t, count, 2)
}
