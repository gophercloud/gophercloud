package pagination

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// LinkedPageBase may be embedded to implement a page that provides navigational "Next" and "Previous" links within its result.
type LinkedPageBase LastHTTPResponse

// NextPageURL extracts the pagination structure from a JSON response and returns the "next" link, if one is present.
// It assumes that the links are available in a "links" element of the top-level response object.
// If this is not the case, override NextPageURL on your result type.
func (current LinkedPageBase) NextPageURL() (string, error) {
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
func NewLinkedPager(client *gophercloud.ServiceClient, initialURL string, createPage func(resp LastHTTPResponse) Page) Pager {
	fetchNextPage := func(url string) (Page, error) {
		resp, err := Request(client, url)
		if err != nil {
			return nil, err
		}

		cp, err := RememberHTTPResponse(resp)
		if err != nil {
			return nil, err
		}

		return createPage(cp), nil
	}

	return Pager{
		initialURL:    initialURL,
		fetchNextPage: fetchNextPage,
	}
}
