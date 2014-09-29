package containers

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts is a structure that holds options for listing containers.
type ListOpts struct {
	Full      bool
	Limit     int    `q:"limit"`
	Marker    string `q:"marker"`
	EndMarker string `q:"end_marker"`
	Format    string `q:"format"`
	Prefix    string `q:"prefix"`
	Delimiter []byte `q:"delimiter"`
}

// List is a function that retrieves all objects in a container. It also returns the details
// for the account. To extract just the container information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	var headers map[string]string

	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}

	if !opts.Full {
		headers = map[string]string{"Accept": "text/plain"}
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ContainerPage{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	url := accountURL(c) + query
	pager := pagination.NewPager(c, url, createPage)
	pager.Headers = headers
	return pager
}

// CreateOpts is a structure that holds parameters for creating a container.
type CreateOpts struct {
	Metadata          map[string]string
	ContainerRead     string `h:"X-Container-Read"`
	ContainerSyncTo   string `h:"X-Container-Sync-To"`
	ContainerSyncKey  string `h:"X-Container-Sync-Key"`
	ContainerWrite    string `h:"X-Container-Write"`
	ContentType       string `h:"Content-Type"`
	DetectContentType bool   `h:"X-Detect-Content-Type"`
	IfNoneMatch       string `h:"If-None-Match"`
	VersionsLocation  string `h:"X-Versions-Location"`
}

// Create is a function that creates a new container.
func Create(c *gophercloud.ServiceClient, containerName string, opts CreateOpts) (Container, error) {
	var container Container
	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return container, err
	}

	for k, v := range headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	_, err = perigee.Request("PUT", containerURL(c, containerName), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201, 204},
	})
	if err == nil {
		container = Container{Name: containerName}
	}
	return container, err
}

// Delete is a function that deletes a container.
func Delete(c *gophercloud.ServiceClient, containerName string) error {
	_, err := perigee.Request("DELETE", containerURL(c, containerName), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}

// UpdateOpts is a structure that holds parameters for updating, creating, or deleting a
// container's metadata.
type UpdateOpts struct {
	Metadata               map[string]string
	ContainerRead          string `h:"X-Container-Read"`
	ContainerSyncTo        string `h:"X-Container-Sync-To"`
	ContainerSyncKey       string `h:"X-Container-Sync-Key"`
	ContainerWrite         string `h:"X-Container-Write"`
	ContentType            string `h:"Content-Type"`
	DetectContentType      bool   `h:"X-Detect-Content-Type"`
	RemoveVersionsLocation string `h:"X-Remove-Versions-Location"`
	VersionsLocation       string `h:"X-Versions-Location"`
}

// Update is a function that creates, updates, or deletes a container's metadata.
func Update(c *gophercloud.ServiceClient, containerName string, opts UpdateOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return err
	}

	for k, v := range headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	url := containerURL(c, containerName)
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// Get is a function that retrieves the metadata of a container. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, containerName string) GetResult {
	var gr GetResult
	resp, err := perigee.Request("HEAD", containerURL(c, containerName), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	gr.Err = err
	gr.Resp = &resp.HttpResponse
	return gr
}
