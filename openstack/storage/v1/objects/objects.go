package objects

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Object map[string]interface{}

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
