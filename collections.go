package gophercloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

var (
	// ErrPageNotAvailable is returned from a Pager when a next or previous page is requested, but does not exist.
	ErrPageNotAvailable = errors.New("The requested Collection page does not exist.")
)

// Page must be satisfied by the result type of any resource collection.
// It allows clients to interact with the resource uniformly, regardless of whether or not or how it's paginated.
type Page interface {

	// NextPageURL generates the URL for the page of data that follows this collection.
	// Return "" if no such page exists.
	NextPageURL() (string, error)
}

// Pager knows how to advance through a specific resource collection, one page at a time.
type Pager struct {
	initialURL string

	advance func(string) (Page, error)
}

// NewPager constructs a manually-configured pager.
// Supply the URL for the first page and a function that requests a specific page given a URL.
func NewPager(initialURL string, advance func(string) (Page, error)) Pager {
	return Pager{
		initialURL: initialURL,
		advance:    advance,
	}
}

// EachPage iterates over each page returned by a Pager, yielding one at a time to a handler function.
// Return "false" from the handler to prematurely stop iterating.
func (p Pager) EachPage(handler func(Page) bool) error {
	currentURL := p.initialURL
	for {
		currentPage, err := p.advance(currentURL)
		if err != nil {
			return err
		}

		if !handler(currentPage) {
			return nil
		}

		currentURL, err = currentPage.NextPageURL()
		if err != nil {
			return err
		}
		if currentURL == "" {
			return nil
		}
	}
}

// ConcretePage stores generic information derived from an HTTP response.
type ConcretePage struct {
	http.Header
	Body map[string]interface{}
}

// NewConcretePage parses an HTTP response as JSON and returns a ConcretePage containing the results.
func NewConcretePage(resp http.Response) (ConcretePage, error) {
	var parsedBody map[string]interface{}

	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ConcretePage{}, err
	}
	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		return ConcretePage{}, err
	}

	return ConcretePage{Header: resp.Header, Body: parsedBody}, err
}

// SinglePage is a page that contains all of the results from an operation.
type SinglePage ConcretePage

// NextPageURL always returns "" to indicate that there are no more pages to return.
func (current SinglePage) NextPageURL() (string, error) {
	return "", nil
}

// NewSinglePager constructs a Pager that "iterates" over a single Page.
// Supply a function that returns the only page.
func NewSinglePager(only func() (http.Response, error)) Pager {
	consumed := false
	single := func(_ string) (Page, error) {
		if !consumed {
			consumed = true
			resp, err := only()
			if err != nil {
				return SinglePage{}, err
			}

			cp, err := NewConcretePage(resp)
			if err != nil {
				return SinglePage{}, err
			}
			return SinglePage(cp), nil
		}
		return SinglePage{}, ErrPageNotAvailable
	}

	return Pager{
		initialURL: "",
		advance:    single,
	}
}

// LinkedPage is a page in a collection that provides navigational "Next" and "Previous" links within its result.
type LinkedPage ConcretePage

// NextPageURL extracts the pagination structure from a JSON response and returns the "next" link, if one is present.
func (current LinkedPage) NextPageURL() (string, error) {
	type response struct {
		Links struct {
			Next *string `mapstructure:"next,omitempty"`
		} `mapstructure:"links"`
	}

	var r response
	err := mapstructure.Decode(current.Body, &r)
	if err != nil {
		return "", err
	}

	if r.Links.Next == nil {
		return "", nil
	}

	return *r.Links.Next, nil
}

// NewLinkedPager creates a Pager that uses a "links" element in the JSON response to locate the next page.
func NewLinkedPager(initialURL string, request func(string) (http.Response, error)) Pager {
	advance := func(url string) (Page, error) {
		resp, err := request(url)
		if err != nil {
			return nil, err
		}

		cp, err := NewConcretePage(resp)
		if err != nil {
			return nil, err
		}

		return LinkedPage(cp), nil
	}

	return Pager{
		initialURL: initialURL,
		advance:    advance,
	}
}
