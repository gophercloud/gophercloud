package accounts

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// UpdateOpts is a structure that contains parameters for updating, creating, or deleting an
// account's metadata.
type UpdateOpts struct {
	Metadata map[string]string
	Headers  map[string]string
}

// Update is a function that creates, updates, or deletes an account's metadata.
func Update(c *gophercloud.ServiceClient, opts UpdateOpts) UpdateResult {
	headers := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		headers[k] = v
	}

	for k, v := range opts.Metadata {
		headers["X-Account-Meta-"+k] = v
	}

	var res UpdateResult

	var resp *perigee.Response

	resp, res.Err = perigee.Request("POST", accountURL(c), perigee.Options{
		MoreHeaders: headers,
		OkCodes:     []int{204},
	})

	res.Resp = &resp.HttpResponse

	return res
}

// GetOpts is a structure that contains parameters for getting an account's metadata.
type GetOpts struct {
	Headers map[string]string
}

// Get is a function that retrieves an account's metadata. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, opts GetOpts) GetResult {
	headers := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		headers[k] = v
	}

	var res GetResult
	var resp *perigee.Response

	resp, res.Err = perigee.Request("HEAD", accountURL(c), perigee.Options{
		MoreHeaders: headers,
		Results:     &res.Resp,
		OkCodes:     []int{204},
	})

	res.Resp = &resp.HttpResponse

	return res
}
