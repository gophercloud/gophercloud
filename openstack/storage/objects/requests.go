package objects

import (
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud/openstack/storage"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type ListResult *perigee.Response
type DownloadResult *perigee.Response
type GetResult *perigee.Response

// List is a function that retrieves all objects in a container. It also returns the details
// for the container. To extract only the object information or names, pass the ListResult
// response to the GetInfo or GetNames function, respectively.
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

	url := c.GetContainerURL(opts.Container) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		Results:     true,
		MoreHeaders: h,
		OkCodes:     []int{200, 204},
		Accept:      contentType,
	})
	return resp, err
}

// Download is a function that retrieves the content and metadata for an object.
// To extract just the content, pass the DownloadResult response to the GetContent
// function.
func Download(c *storage.Client, opts DownloadOpts) (DownloadResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	query := utils.BuildQuery(opts.Params)

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		Results:     true,
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	return resp, err
}

// Create is a function that creates a new object or replaces an existing object.
func Create(c *storage.Client, opts CreateOpts) error {
	var reqBody []byte

	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	query := utils.BuildQuery(opts.Params)

	content := opts.Content
	if content != nil {
		reqBody = make([]byte, content.Len())
		_, err = content.Read(reqBody)
		if err != nil {
			return err
		}
	}

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	_, err = perigee.Request("PUT", url, perigee.Options{
		ReqBody:     reqBody,
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// Copy is a function that copies one object to another.
func Copy(c *storage.Client, opts CopyOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	h["Destination"] = fmt.Sprintf("/%s/%s", opts.NewContainer, opts.NewName)

	url := c.GetObjectURL(opts.Container, opts.Name)
	_, err = perigee.Request("COPY", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// Delete is a function that deletes an object.
func Delete(c *storage.Client, opts DeleteOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	query := utils.BuildQuery(opts.Params)

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	_, err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Get is a function that retrieves the metadata of an object. To extract just the custom
// metadata, pass the GetResult response to the GetMetadata function.
func Get(c *storage.Client, opts GetOpts) (GetResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	url := c.GetObjectURL(opts.Container, opts.Name)
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return resp, err
}

// Update is a function that creates, updates, or deletes an object's metadata.
func Update(c *storage.Client, opts UpdateOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	url := c.GetObjectURL(opts.Container, opts.Name)
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{202, 204},
	})
	return err
}
