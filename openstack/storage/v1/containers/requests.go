package containers

import (
	"github.com/racker/perigee"
	storage "github.com/rackspace/gophercloud/openstack/storage/v1"
	"github.com/rackspace/gophercloud/openstack/utils"
	"net/http"
)

// ListResult is a *http.Response that is returned from a call to the List function.
type ListResult *http.Response

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// List is a function that retrieves all objects in a container. It also returns the details
// for the account. To extract just the container information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *storage.Client, opts ListOpts) (ListResult, error) {
	contentType := ""

	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	query := utils.BuildQuery(opts.Params)

	if !opts.Full {
		contentType = "text/plain"
	}

	url := c.GetAccountURL() + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
		Accept:      contentType,
	})
	return &resp.HttpResponse, err
}

// Create is a function that creates a new container.
func Create(c *storage.Client, opts CreateOpts) (Container, error) {
	var ci Container

	h, err := c.GetHeaders()
	if err != nil {
		return Container{}, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := c.GetContainerURL(opts.Name)
	_, err = perigee.Request("PUT", url, perigee.Options{
		MoreHeaders: h,
	})
	if err == nil {
		ci = Container{
			"name": opts.Name,
		}
	}
	return ci, err
}

// Delete is a function that deletes a container.
func Delete(c *storage.Client, opts DeleteOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	query := utils.BuildQuery(opts.Params)

	url := c.GetContainerURL(opts.Name) + query
	_, err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}

// Update is a function that creates, updates, or deletes a container's metadata.
func Update(c *storage.Client, opts UpdateOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := c.GetContainerURL(opts.Name)
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}

// Get is a function that retrieves the metadata of a container. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *storage.Client, opts GetOpts) (GetResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := c.GetContainerURL(opts.Name)
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
	})
	return &resp.HttpResponse, err
}
