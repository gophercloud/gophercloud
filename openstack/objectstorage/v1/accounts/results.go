package accounts

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// UpdateResult is returned from a call to the Update function.
type UpdateResult struct {
	gophercloud.HeaderResult
}

// UpdateHeader represents the headers returned in the response from an Update request.
type UpdateHeader struct {
	ContentLength string    `json:"Content-Length"`
	ContentType   []string  `json:"Content-Type"`
	Date          time.Time `json:"-"`
	TransID       string    `json:"X-Trans_ID"`
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (ur UpdateResult) Extract() (UpdateHeader, error) {
	var uh UpdateHeader
	if ur.Err != nil {
		return uh, ur.Err
	}

	b, err := json.Marshal(ur.Header)
	if err != nil {
		return uh, err
	}

	err = json.Unmarshal(b, &uh)
	if err != nil {
		return uh, err
	}

	date, err := time.Parse(time.RFC1123, ur.Header["Date"][0])
	if err != nil {
		return uh, err
	}

	uh.Date = date

	return uh, nil
}

// GetHeader represents the headers returned in the response from a Get request.
type GetHeader struct {
	BytesUsed      int64     `json:"X-Account-Bytes-Used"`
	ContainerCount int       `json:"X-Accound-Container-Count"`
	ContentLength  int64     `json:"Content-Length"`
	ContentType    string    `json:"Content-Type"`
	Date           time.Time `mapstructure:"-" json:"-"`
	ObjectCount    int64     `json:"X-Account-Object-Count"`
	TransID        string    `json:"X-Trans-Id"`
}

// GetResult is returned from a call to the Get function.
type GetResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (gr GetResult) Extract() (GetHeader, error) {
	fmt.Printf("raw response header: %+v\n", gr.Header)

	var gh GetHeader

	if err := mapstructure.Decode(gr.Header, &gh); err != nil {
		return gh, err
	}

	t, err := time.Parse(time.RFC1123, gr.Header["Date"][0])
	if err != nil {
		return gh, err
	}
	gh.Date = t

	return gh, nil
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (gr GetResult) ExtractMetadata() (map[string]string, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}

	metadata := make(map[string]string)
	for k, v := range gr.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}
