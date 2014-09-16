package containers

import (
	"fmt"
	"strings"

	"github.com/rackspace/gophercloud/pagination"
)

// Container is a structure that holds information related to a storage container.
type Container map[string]interface{}

// ListOpts is a structure that holds parameters for listing containers.
type ListOpts struct {
	Full   bool
	Params map[string]string
}

// CreateOpts is a structure that holds parameters for creating a container.
type CreateOpts struct {
	Name     string
	Metadata map[string]string
	Headers  map[string]string
}

// DeleteOpts is a structure that holds parameters for deleting a container.
type DeleteOpts struct {
	Name   string
	Params map[string]string
}

// UpdateOpts is a structure that holds parameters for updating, creating, or deleting a
// container's metadata.
type UpdateOpts struct {
	Name     string
	Metadata map[string]string
	Headers  map[string]string
}

// GetOpts is a structure that holds parameters for getting a container's metadata.
type GetOpts struct {
	Name string
}

// ExtractInfo is a function that takes a ListResult and returns the containers' information.
func ExtractInfo(page pagination.Page) ([]Container, error) {
	untyped := page.(ListResult).Body.([]interface{})
	results := make([]Container, len(untyped))
	for index, each := range untyped {
		results[index] = Container(each.(map[string]interface{}))
	}
	return results, nil
}

// ExtractNames is a function that takes a ListResult and returns the containers' names.
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
		for _, container := range parsed {
			names = append(names, container["name"].(string))
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

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metadata associated with the container.
func ExtractMetadata(gr GetResult) map[string]string {
	metadata := make(map[string]string)
	for k, v := range gr.Header {
		if strings.HasPrefix(k, "X-Container-Meta-") {
			key := strings.TrimPrefix(k, "X-Container-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata
}
