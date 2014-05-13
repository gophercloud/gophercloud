package accounts

import (
	"github.com/racker/perigee"
	storage "github.com/rackspace/gophercloud/openstack/storage/v1"
	"net/http"
)

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// Update is a function that creates, updates, or deletes an account's metadata.
func Update(c *storage.Client, opts UpdateOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}
	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Account-Meta-"+k] = v
	}

	url := c.GetAccountURL()
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}

// Get is a function that retrieves an account's metadata. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *storage.Client, opts GetOpts) (GetResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	url := c.GetAccountURL()
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
	})
	return &resp.HttpResponse, err
}
