package accounts

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
)

// UpdateResult is returned from a call to the Update function.
type UpdateResult struct {
	gophercloud.HeaderResult
}

// UpdateHeader represents the headers returned in the response from an Update request.
type UpdateHeader struct {
	ContentLength int64     `json:"-"`
	ContentType   string    `json:"Content-Type"`
	TransID       string    `json:"X-Trans-Id"`
	Date          time.Time `json:"-"`
}

func (s *UpdateHeader) UnmarshalJSON(b []byte) error {
	type tmp UpdateHeader
	var p *struct {
		tmp
		ContentLength string                  `json:"Content-Length"`
		Date          gophercloud.JSONRFC1123 `json:"Date"`
	}
	err := json.Unmarshal(b, &p)
	if err != nil {
		return err
	}

	*s = UpdateHeader(p.tmp)

	switch p.ContentLength {
	case "":
		s.ContentLength = 0
	default:
		s.ContentLength, err = strconv.ParseInt(p.ContentLength, 10, 64)
		if err != nil {
			return err
		}
	}

	s.Date = time.Time(p.Date)

	return err
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (ur UpdateResult) Extract() (*UpdateHeader, error) {
	var uh *UpdateHeader
	err := ur.ExtractInto(&uh)
	return uh, err
}

// GetHeader represents the headers returned in the response from a Get request.
type GetHeader struct {
	BytesUsed      int64     `json:"-"`
	ContainerCount int64     `json:"-"`
	ContentLength  int64     `json:"-"`
	ObjectCount    int64     `json:"-"`
	ContentType    string    `json:"Content-Type"`
	TransID        string    `json:"X-Trans-Id"`
	TempURLKey     string    `json:"X-Account-Meta-Temp-URL-Key"`
	TempURLKey2    string    `json:"X-Account-Meta-Temp-URL-Key-2"`
	Date           time.Time `json:"-"`
}

func (s *GetHeader) UnmarshalJSON(b []byte) error {
	type tmp GetHeader
	var p *struct {
		tmp
		BytesUsed      string `json:"X-Account-Bytes-Used"`
		ContentLength  string `json:"Content-Length"`
		ContainerCount string `json:"X-Account-Container-Count"`
		ObjectCount    string `json:"X-Account-Object-Count"`
		Date           string `json:"Date"`
	}
	err := json.Unmarshal(b, &p)
	if err != nil {
		return err
	}

	*s = GetHeader(p.tmp)

	switch p.BytesUsed {
	case "":
		s.BytesUsed = 0
	default:
		s.BytesUsed, err = strconv.ParseInt(p.BytesUsed, 10, 64)
		if err != nil {
			return err
		}
	}

	switch p.ContentLength {
	case "":
		s.ContentLength = 0
	default:
		s.ContentLength, err = strconv.ParseInt(p.ContentLength, 10, 64)
		if err != nil {
			return err
		}
	}

	switch p.ObjectCount {
	case "":
		s.ObjectCount = 0
	default:
		s.ObjectCount, err = strconv.ParseInt(p.ObjectCount, 10, 64)
		if err != nil {
			return err
		}
	}

	switch p.ContainerCount {
	case "":
		s.ContainerCount = 0
	default:
		s.ContainerCount, err = strconv.ParseInt(p.ContainerCount, 10, 64)
		if err != nil {
			return err
		}
	}

	if p.Date != "" {
		s.Date, err = time.Parse(time.RFC1123, p.Date)
	}

	return err
}

// GetResult is returned from a call to the Get function.
type GetResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Get. To obtain
// a map of headers, call the ExtractHeader method on the GetResult.
func (r GetResult) Extract() (*GetHeader, error) {
	var s *GetHeader
	err := r.ExtractInto(&s)
	return s, err
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metatdata associated with the account.
func (r GetResult) ExtractMetadata() (map[string]string, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	metadata := make(map[string]string)
	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-Account-Meta-") {
			key := strings.TrimPrefix(k, "X-Account-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}
