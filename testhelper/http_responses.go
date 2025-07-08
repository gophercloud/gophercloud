package testhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type FakeServer struct {
	// Mux is a multiplexer that can be used to register handlers.
	Mux *http.ServeMux

	// Server is an in-memory HTTP server for testing.
	Server *httptest.Server
}

func (fakeServer FakeServer) Teardown() {
	fakeServer.Server.Close()
}

func (fakeServer FakeServer) Endpoint() string {
	return fakeServer.Server.URL + "/"
}

// Serves a static content at baseURL/relPath
func (fakeServer FakeServer) ServeFile(t *testing.T, baseURL, relPath, contentType, content string) string {
	rawURL := strings.Join([]string{baseURL, relPath}, "/")
	parsedURL, err := url.Parse(rawURL)
	AssertNoErr(t, err)
	fakeServer.Mux.HandleFunc(parsedURL.Path, func(w http.ResponseWriter, r *http.Request) {
		TestMethod(t, r, "GET")
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, content)
	})

	return rawURL
}

// SetupPersistentPortHTTP prepares the Mux and Server listening specific port.
func SetupPersistentPortHTTP(t *testing.T, port int) FakeServer {
	mux := http.NewServeMux()
	server := httptest.NewUnstartedServer(mux)
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Errorf("Failed to listen to 127.0.0.1:%d: %s", port, err)
	}
	server.Listener = l
	server.Start()

	return FakeServer{
		Mux:    mux,
		Server: server,
	}
}

// SetupHTTP prepares the Mux and Server.
func SetupHTTP() FakeServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return FakeServer{
		Mux:    mux,
		Server: server,
	}
}

// TestFormValues ensures that all the URL parameters given to the http.Request are the same as values.
func TestFormValues(t *testing.T, r *http.Request, values map[string]string) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	if err := r.ParseForm(); err != nil {
		t.Errorf("Failed to parse request form %v", r)
	}
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
	if len(r.Header.Values(header)) == 0 {
		t.Errorf("Header %s not found, expected %q", header, expected)
		return
	}
	for _, actual := range r.Header.Values(header) {
		if expected != actual {
			t.Errorf("Header %s = %q, expected %q", header, actual, expected)
		}
	}
}

// TestHeaderUnset checks that the header on the http.Request doesn't exist.
func TestHeaderUnset(t *testing.T, r *http.Request, header string) {
	if len(r.Header.Values(header)) > 0 {
		t.Errorf("Header %s is not expected", header)
	}
}

// TestBody verifies that the request body matches an expected body.
func TestBody(t *testing.T, r *http.Request, expected string) {
	b, err := io.ReadAll(r.Body)
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
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read request body: %v", err)
	}

	var actualJSON any
	err = json.Unmarshal(b, &actualJSON)
	if err != nil {
		t.Errorf("Unable to parse request body as JSON: %v", err)
	}

	CheckJSONEquals(t, expected, actualJSON)
}
