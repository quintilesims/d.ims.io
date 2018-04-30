package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/zpatrick/fireball"
)

func unmarshalBody(t *testing.T, resp fireball.Response, v interface{}) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	resp.Write(recorder, nil)

	if v != nil {
		if err := json.Unmarshal(recorder.Body.Bytes(), v); err != nil {
			t.Fatal(err)
		}
	}

	return recorder
}
