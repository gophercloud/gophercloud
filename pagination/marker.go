package pagination

import "github.com/rackspace/gophercloud"

// MarkerPage is a stricter Page interface that describes additional functionality required for use with NewMarkerPager.
// For convenience, embed the MarkedPageBase struct.
type MarkerPage interface {
	Page

	// LastMark returns the last "marker" value on this page.
	LastMark() (string, error)
}

// MarkerPageBase is a page in a collection that's paginated by "limit" and "marker" query parameters.
type MarkerPageBase struct {
	LastHTTPResponse

	// A reference to the embedding struct.
	Self MarkerPage
}

// NextPageURL generates the URL for the page of results after this one.
func (current MarkerPageBase) NextPageURL() (string, error) {
	currentURL := current.URL

	mark, err := current.Self.LastMark()
	if err != nil {
		return "", err
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// NewMarkerPager creates a Pager that iterates over successive pages by issuing requests with a "marker" parameter set to the
// final element of the previous Page.
func NewMarkerPager(client *gophercloud.ServiceClient, initialURL string, createPage func(resp LastHTTPResponse) MarkerPage) Pager {

	fetchNextPage := func(currentURL string) (Page, error) {
		resp, err := Request(client, currentURL)
		if err != nil {
			return nullPage{}, err
		}

		last, err := RememberHTTPResponse(resp)
		if err != nil {
			return nullPage{}, err
		}

		return createPage(last), nil
	}

	return Pager{
		initialURL:    initialURL,
		fetchNextPage: fetchNextPage,
	}
}
