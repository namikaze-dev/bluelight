package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/namikaze-dev/bluelight/internal/models"
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
		return
	}

	var got models.Movie
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}

	want := models.Movie{
		ID: 1,
	}

	assertEqual(t, got.ID, want.ID)
}
