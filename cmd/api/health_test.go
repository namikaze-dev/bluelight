package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
	if err != nil {
		t.Fatal("unexpected error", err)
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
	want := `{"status": "available", "environment": "development", "version": "1.0.0"}`
	assertEqual(t, string(body), want)
}
