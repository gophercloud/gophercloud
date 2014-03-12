package containers

import (
	"encoding/json"
	"strings"
)

type Container struct {
	Bytes int
	Count int
	Name  string
}

type ListOpts struct {
	Full   bool
	Params map[string]string
}

type CreateOpts struct {
	Name     string
	Metadata map[string]string
	Headers  map[string]string
}

type DeleteOpts struct {
	Name   string
	Params map[string]string
}

type UpdateOpts struct {
	Name     string
	Metadata map[string]string
	Headers  map[string]string
}

type GetOpts struct {
	Name     string
	Metadata map[string]string
}

// GetInfo is a function that takes a ListResult (of type *perigee.Response)
// and returns the containers' information.
func GetInfo(lr ListResult) ([]Container, error) {
	var ci []Container
	err := json.Unmarshal(lr.JsonResult, &ci)
	return ci, err
}

// GetNames is a function that takes a ListResult (of type *perigee.Response)
// and returns the containers' names.
func GetNames(lr ListResult) ([]string, error) {
	jr := string(lr.JsonResult)
	cns := strings.Split(jr, "\n")
	cns = cns[:len(cns)-1]
	return cns, nil
}

// GetMetadata is a function that takes a GetResult (of type *perigee.Response)
// and returns the custom metadata associated with the container.
func GetMetadata(gr GetResult) map[string]string {
	metadata := make(map[string]string)
	for k, v := range gr.HttpResponse.Header {
		if strings.HasPrefix(k, "X-Container-Meta-") {
			key := strings.TrimPrefix(k, "X-Container-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata
}
