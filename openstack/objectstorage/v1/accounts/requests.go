package accounts

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// GetOpts is a structure that contains parameters for getting an account's metadata.
type GetOpts struct {
	Newest bool `h:"X-Newest"`
}

// Get is a function that retrieves an account's metadata. To extract just the custom
// metadata, call the ExtractMetadata method on the GetResult. To extract all the headers that are
// returned (including the metadata), call the ExtractHeaders method on the GetResult.
func Get(c *gophercloud.ServiceClient, opts *GetOpts) GetResult {
	var res GetResult
	h := c.Provider.AuthenticatedHeaders()

	if opts != nil {
		headers, err := gophercloud.BuildHeaders(opts)
		if err != nil {
			res.Err = err
			return res
		}

		for k, v := range headers {
			h[k] = v
		}
	}

	resp, err := perigee.Request("HEAD", getURL(c), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
}

// UpdateOpts is a structure that contains parameters for updating, creating, or deleting an
// account's metadata.
type UpdateOpts struct {
	Metadata          map[string]string
	ContentType       string `h:"Content-Type"`
	DetectContentType bool   `h:"X-Detect-Content-Type"`
	TempURLKey        string `h:"X-Account-Meta-Temp-URL-Key"`
	TempURLKey2       string `h:"X-Account-Meta-Temp-URL-Key-2"`
}

// Update is a function that creates, updates, or deletes an account's metadata. To extract the
// headers returned, call the ExtractHeaders method on the UpdateResult.
func Update(c *gophercloud.ServiceClient, opts *UpdateOpts) UpdateResult {
	var res UpdateResult
	h := c.Provider.AuthenticatedHeaders()

	if opts != nil {
		headers, err := gophercloud.BuildHeaders(opts)
		if err != nil {
			res.Err = err
			return res
		}
		for k, v := range headers {
			h[k] = v
		}

		for k, v := range opts.Metadata {
			h["X-Account-Meta-"+k] = v
		}
	}

	resp, err := perigee.Request("POST", updateURL(c), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
}
