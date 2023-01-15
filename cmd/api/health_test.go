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

	var got map[string]string
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"status":      "available",
		"environment": "development",
		"version":     "1.0.0",
	}

	assertEqual(t, got["status"], want["status"])
	assertEqual(t, got["environment"], want["environment"])
	assertEqual(t, got["version"], want["version"])
}
