package objects

import (
	"fmt"
	"io"
	"time"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts is a structure that holds parameters for listing objects.
type ListOpts struct {
	Full      bool
	Limit     int     `q:"limit"`
	Marker    string  `q:"marker"`
	EndMarker string  `q:"end_marker"`
	Format    string  `q:"format"`
	Prefix    string  `q:"prefix"`
	Delimiter [1]byte `q:"delimiter"`
	Path      string  `q:"path"`
}

// List is a function that retrieves all objects in a container. It also returns the details
// for the container. To extract only the object information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *gophercloud.ServiceClient, containerName string, opts ListOpts) pagination.Pager {
	var headers map[string]string

	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		fmt.Printf("Error building query string: %v", err)
		return pagination.Pager{Err: err}
	}

	if !opts.Full {
		headers = map[string]string{"Accept": "text/plain"}
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ObjectPage{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	url := containerURL(c, containerName) + query
	pager := pagination.NewPager(c, url, createPage)
	pager.Headers = headers
	return pager
}

// DownloadOpts is a structure that holds parameters for downloading an object.
type DownloadOpts struct {
	IfMatch           string    `h:"If-Match"`
	IfModifiedSince   time.Time `h:"If-Modified-Since"`
	IfNoneMatch       string    `h:"If-None-Match"`
	IfUnmodifiedSince time.Time `h:"If-Unmodified-Since"`
	Range             string    `h:"Range"`
	Expires           string    `q:"expires"`
	MultipartManifest string    `q:"multipart-manifest"`
	Signature         string    `q:"signature"`
}

// Download is a function that retrieves the content and metadata for an object.
// To extract just the content, pass the DownloadResult response to the ExtractContent
// function.
func Download(c *gophercloud.ServiceClient, containerName, objectName string, opts DownloadOpts) DownloadResult {
	var dr DownloadResult

	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		dr.Err = err
		return dr
	}

	for k, v := range headers {
		h[k] = v
	}

	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		dr.Err = err
		return dr
	}

	url := objectURL(c, containerName, objectName) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	dr.Err = err
	dr.Resp = &resp.HttpResponse
	return dr
}

// CreateOpts is a structure that holds parameters for creating an object.
type CreateOpts struct {
	Metadata           map[string]string
	ContentDisposition string `h:"Content-Disposition"`
	ContentEncoding    string `h:"Content-Encoding"`
	ContentLength      int    `h:"Content-Length"`
	ContentType        string `h:"Content-Type"`
	CopyFrom           string `h:"X-Copy-From"`
	DeleteAfter        int    `h:"X-Delete-After"`
	DeleteAt           int    `h:"X-Delete-At"`
	DetectContentType  string `h:"X-Detect-Content-Type"`
	ETag               string `h:"ETag"`
	IfNoneMatch        string `h:"If-None-Match"`
	ObjectManifest     string `h:"X-Object-Manifest"`
	TransferEncoding   string `h:"Transfer-Encoding"`
	Expires            string `q:"expires"`
	MultipartManifest  string `q:"multiple-manifest"`
	Signature          string `q:"signature"`
}

// Create is a function that creates a new object or replaces an existing object.
func Create(c *gophercloud.ServiceClient, containerName, objectName string, content io.Reader, opts CreateOpts) error {
	var reqBody []byte

	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil
	}

	for k, v := range headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return err
	}

	if content != nil {
		reqBody = make([]byte, 0)
		_, err := content.Read(reqBody)
		if err != nil {
			return err
		}
	}

	url := objectURL(c, containerName, objectName) + query
	_, err = perigee.Request("PUT", url, perigee.Options{
		ReqBody:     reqBody,
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// CopyOpts is a structure that holds parameters for copying one object to another.
type CopyOpts struct {
	Metadata           map[string]string
	ContentDisposition string `h:"Content-Disposition"`
	ContentEncoding    string `h:"Content-Encoding"`
	ContentType        string `h:"Content-Type"`
	Destination        string `h:"Destination,required"`
}

// Copy is a function that copies one object to another.
func Copy(c *gophercloud.ServiceClient, containerName, objectName string, opts CopyOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return err
	}
	for k, v := range headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	url := objectURL(c, containerName, objectName)
	_, err = perigee.Request("COPY", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	return err
}

// DeleteOpts is a structure that holds parameters for deleting an object.
type DeleteOpts struct {
	MultipartManifest string `q:"multipart-manifest"`
}

// Delete is a function that deletes an object.
func Delete(c *gophercloud.ServiceClient, containerName, objectName string, opts DeleteOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return err
	}

	url := objectURL(c, containerName, objectName) + query
	_, err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// GetOpts is a structure that holds parameters for getting an object's metadata.
type GetOpts struct {
	Expires   string `q:"expires"`
	Signature string `q:"signature"`
}

// Get is a function that retrieves the metadata of an object. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *gophercloud.ServiceClient, containerName, objectName string, opts GetOpts) GetResult {
	var gr GetResult
	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		gr.Err = err
		return gr
	}

	url := objectURL(c, containerName, objectName) + query
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200, 204},
	})
	gr.Err = err
	gr.Resp = &resp.HttpResponse
	return gr
}

// UpdateOpts is a structure that holds parameters for updating, creating, or deleting an
// object's metadata.
type UpdateOpts struct {
	Metadata           map[string]string
	ContentDisposition string `h:"Content-Disposition"`
	ContentEncoding    string `h:"Content-Encoding"`
	ContentType        string `h:"Content-Type"`
	DeleteAfter        int    `h:"X-Delete-After"`
	DeleteAt           int    `h:"X-Delete-At"`
	DetectContentType  bool   `h:"X-Detect-Content-Type"`
}

// Update is a function that creates, updates, or deletes an object's metadata.
func Update(c *gophercloud.ServiceClient, containerName, objectName string, opts UpdateOpts) error {
	h := c.Provider.AuthenticatedHeaders()

	headers, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil
	}

	for k, v := range headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	url := objectURL(c, containerName, objectName)
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{202},
	})
	return err
}
