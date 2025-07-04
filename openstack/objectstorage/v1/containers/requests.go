package containers

import (
	"bytes"
	"context"
	"net/url"

	"github.com/gophercloud/gophercloud/v2"
	v1 "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToContainerListParams() (string, error)
}

// ListOpts is a structure that holds options for listing containers.
type ListOpts struct {
	// Full has been removed from the Gophercloud API. Gophercloud will now
	// always request the "full" (json) listing, because simplified listing
	// (plaintext) returns false results when names contain end-of-line
	// characters.

	Limit     int    `q:"limit"`
	Marker    string `q:"marker"`
	EndMarker string `q:"end_marker"`
	Format    string `q:"format"`
	Prefix    string `q:"prefix"`
	Delimiter string `q:"delimiter"`
}

// ToContainerListParams formats a ListOpts into a query string.
func (opts ListOpts) ToContainerListParams() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List is a function that retrieves containers associated with the account as
// well as account metadata. It returns a pager which can be iterated with the
// EachPage function.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	headers := map[string]string{"Accept": "application/json", "Content-Type": "application/json"}

	url := listURL(c)
	if opts != nil {
		query, err := opts.ToContainerListParams()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ContainerPage{pagination.MarkerPageBase{PageResult: r}}
		p.Owner = p
		return p
	})
	pager.Headers = headers
	return pager
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToContainerCreateMap() (map[string]string, error)
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
	HistoryLocation   string `h:"X-History-Location"`
	TempURLKey        string `h:"X-Container-Meta-Temp-URL-Key"`
	TempURLKey2       string `h:"X-Container-Meta-Temp-URL-Key-2"`
	StoragePolicy     string `h:"X-Storage-Policy"`
	VersionsEnabled   bool   `h:"X-Versions-Enabled"`
}

// ToContainerCreateMap formats a CreateOpts into a map of headers.
func (opts CreateOpts) ToContainerCreateMap() (map[string]string, error) {
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, err
	}
	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}
	return h, nil
}

// Create is a function that creates a new container.
func Create(ctx context.Context, c *gophercloud.ServiceClient, containerName string, opts CreateOptsBuilder) (r CreateResult) {
	url, err := createURL(c, containerName)
	if err != nil {
		r.Err = err
		return
	}
	h := make(map[string]string)
	if opts != nil {
		headers, err := opts.ToContainerCreateMap()
		if err != nil {
			r.Err = err
			return
		}
		for k, v := range headers {
			h[k] = v
		}
	}
	resp, err := c.Request(ctx, "PUT", url, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{201, 202, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// BulkDelete is a function that bulk deletes containers.
func BulkDelete(ctx context.Context, c *gophercloud.ServiceClient, containers []string) (r BulkDeleteResult) {
	var body bytes.Buffer

	for i := range containers {
		if err := v1.CheckContainerName(containers[i]); err != nil {
			r.Err = err
			return
		}
		body.WriteString(url.PathEscape(containers[i]))
		body.WriteRune('\n')
	}

	resp, err := c.Post(ctx, bulkDeleteURL(c), &body, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "text/plain",
		},
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete is a function that deletes a container.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, containerName string) (r DeleteResult) {
	url, err := deleteURL(c, containerName)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Delete(ctx, url, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToContainerUpdateMap() (map[string]string, error)
}

// UpdateOpts is a structure that holds parameters for updating, creating, or
// deleting a container's metadata.
type UpdateOpts struct {
	Metadata               map[string]string
	RemoveMetadata         []string
	ContainerRead          *string `h:"X-Container-Read"`
	ContainerSyncTo        *string `h:"X-Container-Sync-To"`
	ContainerSyncKey       *string `h:"X-Container-Sync-Key"`
	ContainerWrite         *string `h:"X-Container-Write"`
	ContentType            *string `h:"Content-Type"`
	DetectContentType      *bool   `h:"X-Detect-Content-Type"`
	RemoveVersionsLocation string  `h:"X-Remove-Versions-Location"`
	VersionsLocation       string  `h:"X-Versions-Location"`
	RemoveHistoryLocation  string  `h:"X-Remove-History-Location"`
	HistoryLocation        string  `h:"X-History-Location"`
	TempURLKey             string  `h:"X-Container-Meta-Temp-URL-Key"`
	TempURLKey2            string  `h:"X-Container-Meta-Temp-URL-Key-2"`
	VersionsEnabled        *bool   `h:"X-Versions-Enabled"`
}

// ToContainerUpdateMap formats a UpdateOpts into a map of headers.
func (opts UpdateOpts) ToContainerUpdateMap() (map[string]string, error) {
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	for _, k := range opts.RemoveMetadata {
		h["X-Remove-Container-Meta-"+k] = "remove"
	}

	return h, nil
}

// Update is a function that creates, updates, or deletes a container's
// metadata.
func Update(ctx context.Context, c *gophercloud.ServiceClient, containerName string, opts UpdateOptsBuilder) (r UpdateResult) {
	url, err := updateURL(c, containerName)
	if err != nil {
		r.Err = err
		return
	}
	h := make(map[string]string)
	if opts != nil {
		headers, err := opts.ToContainerUpdateMap()
		if err != nil {
			r.Err = err
			return
		}

		for k, v := range headers {
			h[k] = v
		}
	}
	resp, err := c.Request(ctx, "POST", url, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{201, 202, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetOptsBuilder allows extensions to add additional parameters to the Get
// request.
type GetOptsBuilder interface {
	ToContainerGetMap() (map[string]string, error)
}

// GetOpts is a structure that holds options for listing containers.
type GetOpts struct {
	Newest bool `h:"X-Newest"`
}

// ToContainerGetMap formats a GetOpts into a map of headers.
func (opts GetOpts) ToContainerGetMap() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// Get is a function that retrieves the metadata of a container. To extract just
// the custom metadata, pass the GetResult response to the ExtractMetadata
// function.
func Get(ctx context.Context, c *gophercloud.ServiceClient, containerName string, opts GetOptsBuilder) (r GetResult) {
	url, err := getURL(c, containerName)
	if err != nil {
		r.Err = err
		return
	}
	h := make(map[string]string)
	if opts != nil {
		headers, err := opts.ToContainerGetMap()
		if err != nil {
			r.Err = err
			return
		}

		for k, v := range headers {
			h[k] = v
		}
	}
	resp, err := c.Head(ctx, url, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
