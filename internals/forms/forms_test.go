package forms

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/any", nil)

	_ = r.ParseForm()

	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}

	form.Errors.Add("error", "error")

	if form.Valid() {
		t.Error("should be invalid")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/any", nil)

	_ = r.ParseForm()

	form := New(r.PostForm)

	if form.Has("abc") {
		t.Error("shouldn't have a value")
	}

	postedData := url.Values{}
	postedData.Add("abc", "12345")

	// here I create a new request and write data to it
	r = httptest.NewRequest("POST", "/any", strings.NewReader(postedData.Encode()))
	// appropriate header must be set in order to get data from the request then
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_ = r.ParseForm()

	form = New(r.PostForm)

	if !form.Has("abc") {
		t.Error("should have a value for the given key")
	}

}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/any", nil)

	_ = r.ParseForm()

	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	// here I don't write data to a request but pass it inside later
	r = httptest.NewRequest("POST", "/any", nil)

	// pass data to the request before parsing
	r.PostForm = postedData

	// parse ONLY when data is inside yet
	_ = r.ParseForm()

	// OR it's possible to pass postedData directly to New() then it doesn't require to PareForm()
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows doesn't have required fields when it does")
	}

}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("name", "Eugene")
	postedData.Add("email", "123")

	r := httptest.NewRequest("POST", "/any", nil)

	r.PostForm = postedData

	_ = r.ParseForm()

	form := New(r.PostForm)

	if !form.MinLength("name", 5) {
		t.Error("should return true")
	}

	if form.MinLength("email", 5) {
		t.Error("should return false")
	}

	if form.MinLength("abc", 5) {
		t.Error("should return false")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email1", "123@mail")
	postedData.Add("email2", "eugene.eugene@mail.com")

	r := httptest.NewRequest("POST", "/any", nil)

	r.PostForm = postedData

	_ = r.ParseForm()

	form := New(r.PostForm)

	form.IsEmail("email2")
	if !form.Valid() {
		t.Error("should be valid")
	}

	form.IsEmail("email1")
	if form.Valid() {
		t.Error("should be invalid")
	}
}
