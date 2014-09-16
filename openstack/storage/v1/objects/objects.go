package objects

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/rackspace/gophercloud/pagination"
)

// Object is a structure that holds information related to a storage object.
type Object map[string]interface{}

// ListOpts is a structure that holds parameters for listing objects.
type ListOpts struct {
	Container string
	Full      bool
	Params    map[string]string
}

// DownloadOpts is a structure that holds parameters for downloading an object.
type DownloadOpts struct {
	Container string
	Name      string
	Headers   map[string]string
	Params    map[string]string
}

// CreateOpts is a structure that holds parameters for creating an object.
type CreateOpts struct {
	Container string
	Name      string
	Content   io.Reader
	Metadata  map[string]string
	Headers   map[string]string
	Params    map[string]string
}

// CopyOpts is a structure that holds parameters for copying one object to another.
type CopyOpts struct {
	Container    string
	Name         string
	NewContainer string
	NewName      string
	Metadata     map[string]string
	Headers      map[string]string
}

// DeleteOpts is a structure that holds parameters for deleting an object.
type DeleteOpts struct {
	Container string
	Name      string
	Params    map[string]string
}

// GetOpts is a structure that holds parameters for getting an object's metadata.
type GetOpts struct {
	Container string
	Name      string
	Params    map[string]string
}

// UpdateOpts is a structure that holds parameters for updating, creating, or deleting an
// object's metadata.
type UpdateOpts struct {
	Container string
	Name      string
	Metadata  map[string]string
	Headers   map[string]string
}

// ExtractInfo is a function that takes a page of objects and returns their full information.
func ExtractInfo(page pagination.Page) ([]Object, error) {
	untyped := page.(ListResult).Body.([]interface{})
	results := make([]Object, len(untyped))
	for index, each := range untyped {
		results[index] = Object(each.(map[string]interface{}))
	}
	return results, nil
}

// ExtractNames is a function that takes a page of objects and returns only their names.
func ExtractNames(page pagination.Page) ([]string, error) {
	casted := page.(ListResult)
	ct := casted.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		parsed, err := ExtractInfo(page)
		if err != nil {
			return nil, err
		}

		names := make([]string, 0, len(parsed))
		for _, object := range parsed {
			names = append(names, object["name"].(string))
		}
		return names, nil
	case strings.HasPrefix(ct, "text/plain"):
		names := make([]string, 0, 50)

		body := string(page.(ListResult).Body.([]uint8))
		for _, name := range strings.Split(body, "\n") {
			if len(name) > 0 {
				names = append(names, name)
			}
		}

		return names, nil
	default:
		return nil, fmt.Errorf("Cannot extract names from response with content-type: [%s]", ct)
	}
}

// ExtractContent is a function that takes a DownloadResult (of type *http.Response)
// and returns the object's content.
func ExtractContent(dr DownloadResult) ([]byte, error) {
	var body []byte
	defer dr.Body.Close()
	body, err := ioutil.ReadAll(dr.Body)
	if err != nil {
		return body, fmt.Errorf("Error trying to read DownloadResult body: %v", err)
	}
	return body, nil
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metadata associated with the object.
func ExtractMetadata(gr GetResult) map[string]string {
	metadata := make(map[string]string)
	for k, v := range gr.Header {
		if strings.HasPrefix(k, "X-Object-Meta-") {
			key := strings.TrimPrefix(k, "X-Object-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata
}
