package accounts

import (
	"net/http"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// Update is a function that creates, updates, or deletes an account's metadata.
func Update(c *gophercloud.ServiceClient, opts UpdateOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Account-Meta-"+k] = v
	}

	_, err := perigee.Request("POST", accountURL(c), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Get is a function that retrieves an account's metadata. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, opts GetOpts) (GetResult, error) {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	resp, err := perigee.Request("HEAD", accountURL(c), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return &resp.HttpResponse, err
}
