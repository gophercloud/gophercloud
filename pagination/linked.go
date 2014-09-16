package pagination

import "github.com/mitchellh/mapstructure"

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
