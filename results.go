package gophercloud

import (
	"encoding/json"
	"net/http"
)

// Result acts as a base struct that other results can embed.
type Result struct {
	// Body is the payload of the HTTP response from the server. In most cases, this will be the
	// deserialized JSON structure.
	Body interface{}

	// Header contains the HTTP header structure from the original response.
	Header http.Header

	// Err is an error that occurred during the operation. It's deferred until extraction to make
	// it easier to chain operations.
	Err error
}

// PrettyPrintJSON creates a string containing the full response body as pretty-printed JSON.
func (r Result) PrettyPrintJSON() string {
	pretty, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return string(pretty)
}

// RFC3339Milli describes a time format used by API responses.
const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

// Link represents a structure that enables paginated collections how to
// traverse backward or forward. The "Rel" field is usually either "next".
type Link struct {
	Href string `mapstructure:"href"`
	Rel  string `mapstructure:"rel"`
}

// ExtractNextURL attempts to extract the next URL from a JSON structure. It
// follows the common structure of nesting back and next links.
func ExtractNextURL(links []Link) (string, error) {
	var url string

	for _, l := range links {
		if l.Rel == "next" {
			url = l.Href
		}
	}

	if url == "" {
		return "", nil
	}

	return url, nil
}
