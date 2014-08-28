package testhelper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// TestFormValues ensures that all the URL parameters given to the http.Request are the same as values.
func TestFormValues(t *testing.T, r *http.Request, values map[string]string) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if !reflect.DeepEqual(want, r.Form) {
		t.Errorf("Request parameters = %v, want %v", r.Form, want)
	}
}

// TestMethod checks that the Request has the expected method (e.g. GET, POST).
func TestMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

// TestHeader checks that the header on the http.Request matches the expected value.
func TestHeader(t *testing.T, r *http.Request, header string, expected string) {
	if actual := r.Header.Get(header); expected != actual {
		t.Errorf("Header %s = %s, expected %s", header, actual, expected)
	}
}

// TestBody verifies that the request body matches an expected body.
func TestBody(t *testing.T, r *http.Request, expected string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read body: %v", err)
	}
	str := string(b)
	if expected != str {
		t.Errorf("Body = %s, expected %s", str, expected)
	}
}

// TestJSONRequest verifies that the JSON payload of a request matches an expected structure, without asserting things about
// whitespace or ordering.
func TestJSONRequest(t *testing.T, r *http.Request, expected string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read request body: %v", err)
	}

	var expectedJSON interface{}
	err = json.Unmarshal([]byte(expected), &expectedJSON)
	if err != nil {
		t.Errorf("Unable to parse expected value as JSON: %v", err)
	}

	var actualJSON interface{}
	err = json.Unmarshal(b, &actualJSON)
	if err != nil {
		t.Errorf("Unable to parse request body as JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedJSON, actualJSON) {
		t.Errorf("Response body did not contain the correct JSON.")
	}
}
