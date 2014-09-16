package pagination

import "github.com/rackspace/gophercloud"

// SinglePageBase may be embedded in a Page that contains all of the results from an operation at once.
type SinglePageBase LastHTTPResponse

// NextPageURL always returns "" to indicate that there are no more pages to return.
func (current SinglePageBase) NextPageURL() (string, error) {
	return "", nil
}

// NewSinglePager constructs a Pager that "iterates" over a single Page.
// Supply the URL to request and a function that creates a Page of the appropriate type.
func NewSinglePager(client *gophercloud.ServiceClient, onlyURL string, createPage func(resp LastHTTPResponse) Page) Pager {
	consumed := false
	single := func(_ string) (Page, error) {
		if !consumed {
			consumed = true
			resp, err := Request(client, onlyURL)
			if err != nil {
				return nullPage{}, err
			}

			cp, err := RememberHTTPResponse(resp)
			if err != nil {
				return nullPage{}, err
			}
			return createPage(cp), nil
		}
		return nullPage{}, ErrPageNotAvailable
	}

	return Pager{
		initialURL:    "",
		fetchNextPage: single,
	}
}
