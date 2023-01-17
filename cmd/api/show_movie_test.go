package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowMovie(t *testing.T) {
	app := &application{
		config: config{
			env: "development",
		},
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()
	resp, err := ts.Client().Get(ts.URL + "/v1/movies/1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertEqual(t, resp.StatusCode, http.StatusOK)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var js map[string]interface{}
	err = json.Unmarshal(body, &js)
	if err != nil {
		t.Fatal(err)
	}

	js2, ok := js["movie"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected movie, got %v", js["movie"])
	}

	got, ok := js2["id"].(float64)
	if !ok {
		t.Errorf("%T", js2["id"])
		t.Fatalf("expected id of float6464, got %v", js2["id"])
	}

	var want float64 = 1
	assertEqual(t, got, want)
}
