package containers

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts is a structure that holds options for listing containers.
type ListOpts struct {
	Full      bool
	Limit     int     `q:"limit"`
	Marker    string  `q:"marker"`
	EndMarker string  `q:"end_marker"`
	Format    string  `q:"format"`
	Prefix    string  `q:"prefix"`
	Delimiter [1]byte `q:"delimiter"`
}

// List is a function that retrieves containers associated with the account as well as account
// metadata. It returns a pager which can be iterated with the EachPage function.
func List(c *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	var headers map[string]string

	url := listURL(c)
	if opts != nil {
		query, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query.String()

		if !opts.Full {
			headers = map[string]string{"Accept": "text/plain", "Content-Type": "text/plain"}
		}
	} else {
		headers = map[string]string{"Accept": "text/plain", "Content-Type": "text/plain"}
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ContainerPage{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

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
func Create(c *gophercloud.ServiceClient, containerName string, opts *CreateOpts) CreateResult {
	var res CreateResult
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
			h["X-Container-Meta-"+k] = v
		}
	}

	resp, err := perigee.Request("PUT", createURL(c, containerName), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201, 204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
}

// Delete is a function that deletes a container.
func Delete(c *gophercloud.ServiceClient, containerName string) DeleteResult {
	var res DeleteResult
	resp, err := perigee.Request("DELETE", deleteURL(c, containerName), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
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
func Update(c *gophercloud.ServiceClient, containerName string, opts *UpdateOpts) UpdateResult {
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
			h["X-Container-Meta-"+k] = v
		}
	}

	resp, err := perigee.Request("POST", updateURL(c, containerName), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
}

// Get is a function that retrieves the metadata of a container. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, containerName string) GetResult {
	var res GetResult
	resp, err := perigee.Request("HEAD", getURL(c, containerName), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Resp = &resp.HttpResponse
	res.Err = err
	return res
}
