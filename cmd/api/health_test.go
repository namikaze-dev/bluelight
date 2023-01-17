package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := &application{
		config: config{
			env: "development",
		},
	}
	app.healthHandler(rr, r)
	resp := rr.Result()
	defer resp.Body.Close()

	assertEqual(t, resp.StatusCode, http.StatusOK)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var env map[string]interface{}
	err = json.Unmarshal(body, &env)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := env["status"].(string)
	if !ok {
		t.Errorf("value %v is not of type string", env["status"])
	}

	assertEqual(t, got, "available")
}
