package objects

import (
	"fmt"
	"net/http"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

// ListResult is a single page of objects that is returned from a call to the List function.
type ListResult struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult contains no object names.
func (r ListResult) IsEmpty() (bool, error) {
	names, err := ExtractNames(r)
	if err != nil {
		return true, err
	}
	return len(names) == 0, nil
}

// LastMarker returns the last object name in a ListResult.
func (r ListResult) LastMarker() (string, error) {
	names, err := ExtractNames(r)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", nil
	}
	return names[len(names)-1], nil
}

// DownloadResult is a *http.Response that is returned from a call to the Download function.
type DownloadResult *http.Response

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// List is a function that retrieves all objects in a container. It also returns the details
// for the container. To extract only the object information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	var headers map[string]string

	query := utils.BuildQuery(opts.Params)

	if !opts.Full {
		headers = map[string]string{"Content-Type": "text/plain"}
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ListResult{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	url := containerURL(c, opts.Container) + query
	pager := pagination.NewPager(c, url, createPage)
	pager.Headers = headers
	return pager
}

// Download is a function that retrieves the content and metadata for an object.
// To extract just the content, pass the DownloadResult response to the ExtractContent
// function.
func Download(c *gophercloud.ServiceClient, opts DownloadOpts) (DownloadResult, error) {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	query := utils.BuildQuery(opts.Params)

	url := objectURL(c, opts.Container, opts.Name) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	return &resp.HttpResponse, err
}

// Create is a function that creates a new object or replaces an existing object.
func Create(c *gophercloud.ServiceClient, opts CreateOpts) error {
	var reqBody []byte

	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	query := utils.BuildQuery(opts.Params)

	content := opts.Content
	if content != nil {
		reqBody = make([]byte, 0)
		_, err := content.Read(reqBody)
		if err != nil {
			return err
		}
	}

	url := objectURL(c, opts.Container, opts.Name) + query
	_, err := perigee.Request("PUT", url, perigee.Options{
		ReqBody:     reqBody,
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// Copy is a function that copies one object to another.
func Copy(c *gophercloud.ServiceClient, opts CopyOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	h["Destination"] = fmt.Sprintf("/%s/%s", opts.NewContainer, opts.NewName)

	url := objectURL(c, opts.Container, opts.Name)
	_, err := perigee.Request("COPY", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// Delete is a function that deletes an object.
func Delete(c *gophercloud.ServiceClient, opts DeleteOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	query := utils.BuildQuery(opts.Params)

	url := objectURL(c, opts.Container, opts.Name) + query
	_, err := perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Get is a function that retrieves the metadata of an object. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, opts GetOpts) (GetResult, error) {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	url := objectURL(c, opts.Container, opts.Name)
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return &resp.HttpResponse, err
}

// Update is a function that creates, updates, or deletes an object's metadata.
func Update(c *gophercloud.ServiceClient, opts UpdateOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	url := objectURL(c, opts.Container, opts.Name)
	_, err := perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{202},
	})
	return err
}
