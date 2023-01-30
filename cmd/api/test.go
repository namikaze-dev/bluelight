package main

import "testing"

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Fatalf("got: %v; want: %v", actual, expected)
	}
}
