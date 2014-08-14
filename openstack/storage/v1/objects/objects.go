package objects

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
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
	Headers   map[string]string
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

// ExtractInfo is a function that takes a ListResult (of type *http.Response)
// and returns the objects' information.
func ExtractInfo(lr ListResult) ([]Object, error) {
	var oi []Object
	defer lr.Body.Close()
	body, err := ioutil.ReadAll(lr.Body)
	if err != nil {
		return oi, err
	}
	err = json.Unmarshal(body, &oi)
	return oi, err
}

// ExtractNames is a function that takes a ListResult (of type *http.Response)
// and returns the objects' names.
func ExtractNames(lr ListResult) ([]string, error) {
	var ons []string
	defer lr.Body.Close()
	body, err := ioutil.ReadAll(lr.Body)
	if err != nil {
		return ons, err
	}
	jr := string(body)
	ons = strings.Split(jr, "\n")
	ons = ons[:len(ons)-1]
	return ons, nil
}

// ExtractContent is a function that takes a DownloadResult (of type *http.Response)
// and returns the object's content.
func ExtractContent(dr DownloadResult) ([]byte, error) {
	var body []byte
	defer dr.Body.Close()
	body, err := ioutil.ReadAll(dr.Body)
	return body, err
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
