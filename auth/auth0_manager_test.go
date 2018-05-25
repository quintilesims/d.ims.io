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

		assert.Equal(t, req.Username, "invalid username")
		assert.Equal(t, req.Password, "invalid password")

		MarshalAndWrite(t, w, nil, 401)
	}

	auth0Manager, server := newAuth0ManagerAndServer(handler)
	defer server.Close()

	valid, err := auth0Manager.Authenticate("invalid username", "invalid password")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, valid, false)
}
