package containers

import (
	"net/http"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

// ListResult is a *http.Response that is returned from a call to the List function.
type ListResult *http.Response

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// List is a function that retrieves all objects in a container. It also returns the details
// for the account. To extract just the container information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *gophercloud.ServiceClient, opts ListOpts) (ListResult, error) {
	contentType := ""

	h := c.Provider.AuthenticatedHeaders()

	query := utils.BuildQuery(opts.Params)

	if !opts.Full {
		contentType = "text/plain"
	}

	url := getAccountURL(c) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
		Accept:      contentType,
		OkCodes:     []int{200, 204},
	})
	return &resp.HttpResponse, err
}

// Create is a function that creates a new container.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) (Container, error) {
	var ci Container

	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := getContainerURL(c, opts.Name)
	_, err := perigee.Request("PUT", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201, 204},
	})
	if err == nil {
		ci = Container{
			"name": opts.Name,
		}
	}
	return ci, err
}

// Delete is a function that deletes a container.
func Delete(c *gophercloud.ServiceClient, opts DeleteOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	query := utils.BuildQuery(opts.Params)

	url := getContainerURL(c, opts.Name) + query
	_, err := perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Update is a function that creates, updates, or deletes a container's metadata.
func Update(c *gophercloud.ServiceClient, opts UpdateOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := getContainerURL(c, opts.Name)
	_, err := perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Get is a function that retrieves the metadata of a container. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, opts GetOpts) (GetResult, error) {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := getContainerURL(c, opts.Name)
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return &resp.HttpResponse, err
}
