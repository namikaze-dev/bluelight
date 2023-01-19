package main

import (
	"bytes"
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

	var got struct {
		Movie struct {
			ID int64 `json:"id"`
		} `json:"movie"`
	}
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 1
	assertEqual(t, got.Movie.ID, want)
}

func TestCreateMovie(t *testing.T) {
	app := &application{
		config: config{
			env: "development",
		},
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	var input struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime string   `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// test for bad request
	input.Title = "undead"
	input.Genres = []string{"action"}
	input.Year = 2020
	// missing runtime field should yield a bad request response

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	rd := bytes.NewReader(b)
	resp, err := ts.Client().Post(ts.URL+"/v1/movies", "application/json", rd)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assertEqual(t, resp.StatusCode, http.StatusBadRequest)

	// test for unprocessible entity
	input.Title = "undead"
	input.Genres = []string{"action", "action"}
	input.Year = 2020
	input.Runtime = "100 mins"
	// duplicates genres field should yield a unprocessible entity response
	b, err = json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	rd = bytes.NewReader(b)
	resp, err = ts.Client().Post(ts.URL+"/v1/movies", "application/json", rd)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assertEqual(t, resp.StatusCode, http.StatusUnprocessableEntity)

	// test for ok
	input.Title = "undead"
	input.Genres = []string{"action", "horror"}
	input.Year = 2020
	input.Runtime = "100 mins"
	b, err = json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	rd = bytes.NewReader(b)
	resp, err = ts.Client().Post(ts.URL+"/v1/movies", "application/json", rd)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assertEqual(t, resp.StatusCode, http.StatusOK)

	// check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var got struct {
		Movie struct {
			Runtime string `json:"runtime"`
			Title   string `json:"title"`
			Year    int    `json:"year"`
		} `json:"movie"`
	}
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, got.Movie.Runtime, input.Runtime)
	assertEqual(t, got.Movie.Year, input.Year)
	assertEqual(t, got.Movie.Title, input.Title)
}
