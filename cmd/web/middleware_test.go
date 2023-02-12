package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	// .(type) syntax returns type
	switch v := h.(type) {
	case http.Handler:
		// do nothing

	default:
		t.Errorf("type is not http.Handler, but is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	// .(type) syntax returns type
	switch v := h.(type) {
	case http.Handler:
		// do nothing

	default:
		t.Errorf("type is not http.Handler, but is %T", v)
	}
}
