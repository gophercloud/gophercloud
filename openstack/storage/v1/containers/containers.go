package containers

import (
	"encoding/json"
	"io/ioutil"
	"strings"
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
	Name     string
	Metadata map[string]string
}

// ExtractInfo is a function that takes a ListResult (of type *http.Response)
// and returns the containers' information.
func ExtractInfo(lr ListResult) ([]Container, error) {
	var ci []Container
	defer lr.Body.Close()
	body, err := ioutil.ReadAll(lr.Body)
	if err != nil {
		return ci, err
	}
	err = json.Unmarshal(body, &ci)
	return ci, err
}

// ExtractNames is a function that takes a ListResult (of type *http.Response)
// and returns the containers' names.
func ExtractNames(lr ListResult) ([]string, error) {
	var cns []string
	defer lr.Body.Close()
	body, err := ioutil.ReadAll(lr.Body)
	if err != nil {
		return cns, err
	}
	jr := string(body)
	cns = strings.Split(jr, "\n")
	cns = cns[:len(cns)-1]
	return cns, nil
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
