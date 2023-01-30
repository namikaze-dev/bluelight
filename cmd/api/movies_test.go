package main

import (
	"bytes"
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
		models: models.NewMockModels(),
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	// 10 for id of movie 
	resp, err := ts.Client().Get(ts.URL + "/v1/movies/10")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// movie of id 10 should not exist
	assertEqual(t, resp.StatusCode, http.StatusNotFound)

	// 1 for id of movie, already available by default
	resp, err = ts.Client().Get(ts.URL + "/v1/movies/1")
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
		models: models.NewMockModels(),
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
	assertEqual(t, resp.StatusCode, http.StatusCreated)

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

func TestUpdateMovie(t *testing.T) {
	app := &application{
		config: config{
			env: "development",
		},
		models: models.NewMockModels(),
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	var input struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime string   `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	input.Title = "undead"
	input.Genres = []string{"action", "horror"}
	input.Year = 2020
	input.Runtime = "100 mins"
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	rd := bytes.NewReader(b)
	// movie with id 1 available by default
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/v1/movies/1", rd)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := ts.Client().Do(req)
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

func TestDeleteMovie(t *testing.T) {
	app := &application{
		config: config{
			env: "development",
		},
		models: models.NewMockModels(),
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/v1/movies/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// movie of id 10 should not exist
	assertEqual(t, resp.StatusCode, http.StatusNotFound)

	// movie with id 1 available by default
	req, err = http.NewRequest(http.MethodDelete, ts.URL+"/v1/movies/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assertEqual(t, resp.StatusCode, http.StatusOK)
}
