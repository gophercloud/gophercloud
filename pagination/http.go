package pagination

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// LastHTTPResponse stores generic information derived from an HTTP response.
type LastHTTPResponse struct {
	url.URL
	http.Header
	Body interface{}
}

// RememberHTTPResponse parses an HTTP response as JSON and returns a LastHTTPResponse containing the results.
// The main reason to do this instead of holding the response directly is that a response body can only be read once.
// Also, this centralizes the JSON decoding.
func RememberHTTPResponse(resp http.Response) (LastHTTPResponse, error) {
	var parsedBody interface{}

	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LastHTTPResponse{}, err
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal(rawBody, &parsedBody)
		if err != nil {
			return LastHTTPResponse{}, err
		}
	} else {
		parsedBody = rawBody
	}

	return LastHTTPResponse{
		URL:    *resp.Request.URL,
		Header: resp.Header,
		Body:   parsedBody,
	}, err
}

// Request performs a Perigee request and extracts the http.Response from the result.
func Request(client *gophercloud.ServiceClient, url string) (http.Response, error) {
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200, 204},
	})
	if err != nil {
		return http.Response{}, err
	}
	return resp.HttpResponse, nil
}
