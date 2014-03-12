package objects

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Object struct {
	Name          string
	Hash          string
	Bytes         int
	Content_type  string
	Last_modified string
}

type ListOpts struct {
	Container string
	Full      bool
	Params    map[string]string
}

type DownloadOpts struct {
	Container string
	Name      string
	Headers   map[string]string
	Params    map[string]string
}

type CreateOpts struct {
	Container string
	Name      string
	Content   *bytes.Buffer
	Metadata  map[string]string
	Headers   map[string]string
	Params    map[string]string
}

type CopyOpts struct {
	Container    string
	Name         string
	NewContainer string
	NewName      string
	Metadata     map[string]string
	Headers      map[string]string
}

type DeleteOpts struct {
	Container string
	Name      string
	Params    map[string]string
}

type GetOpts struct {
	Container string
	Name      string
	Headers   map[string]string
	Params    map[string]string
}

type UpdateOpts struct {
	Container string
	Name      string
	Metadata  map[string]string
	Headers   map[string]string
}

// GetInfo is a function that takes a ListResult (of type *perigee.Response)
// and returns the objects' information.
func GetInfo(lr ListResult) ([]Object, error) {
	var oi []Object
	err := json.Unmarshal(lr.JsonResult, &oi)
	return oi, err
}

// GetNames is a function that takes a ListResult (of type *perigee.Response)
// and returns the objects' names.
func GetNames(lr ListResult) []string {
	jr := string(lr.JsonResult)
	ons := strings.Split(jr, "\n")
	ons = ons[:len(ons)-1]
	return ons
}

// GetContent is a function that takes a DownloadResult (of type *perigee.Response)
// and returns the object's content.
func GetContent(dr DownloadResult) []byte {
	return dr.JsonResult
}

// GetMetadata is a function that takes a GetResult (of type *perifee.Response)
// and returns the custom metadata associated with the object.
func GetMetadata(gr GetResult) map[string]string {
	metadata := make(map[string]string)
	for k, v := range gr.HttpResponse.Header {
		if strings.HasPrefix(k, "X-Object-Meta-") {
			key := strings.TrimPrefix(k, "X-Object-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata
}
