package objects

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Object is a structure that holds information related to a storage object.
type Object struct {
	// Bytes is the total number of bytes that comprise the object.
	Bytes int64 `json:"bytes"`

	// ContentType is the content type of the object.
	ContentType string `json:"content_type"`

	// Hash represents the MD5 checksum value of the object's content.
	Hash string `json:"hash"`

	// LastModified is the time the object was last modified.
	LastModified time.Time `json:"-"`

	// Name is the unique name for the object.
	Name string `json:"name"`

	// Subdir denotes if the result contains a subdir.
	Subdir string `json:"subdir"`

	// IsLatest indicates whether the object version is the latest one.
	IsLatest bool `json:"is_latest"`

	// VersionID contains a version ID of the object, when container
	// versioning is enabled.
	VersionID string `json:"version_id"`
}

func (r *Object) UnmarshalJSON(b []byte) error {
	type tmp Object
	var s *struct {
		tmp
		LastModified string `json:"last_modified"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Object(s.tmp)

	if s.LastModified != "" {
		t, err := time.Parse(gophercloud.RFC3339MilliNoZ, s.LastModified)
		if err != nil {
			t, err = time.Parse(gophercloud.RFC3339Milli, s.LastModified)
			if err != nil {
				return err
			}
		}
		r.LastModified = t
	}

	return nil
}

// ObjectPage is a single page of objects that is returned from a call to the
// List function.
type ObjectPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult contains no object names.
func (r ObjectPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	names, err := ExtractNames(r)
	return len(names) == 0, err
}

// LastMarker returns the last object name in a ListResult.
func (r ObjectPage) LastMarker() (string, error) {
	return extractLastMarker(r)
}

// ExtractInfo is a function that takes a page of objects and returns their
// full information.
func ExtractInfo(r pagination.Page) ([]Object, error) {
	var s []Object
	err := (r.(ObjectPage)).ExtractInto(&s)
	return s, err
}

// ExtractNames is a function that takes a page of objects and returns only
// their names.
func ExtractNames(r pagination.Page) ([]string, error) {
	casted := r.(ObjectPage)
	ct := casted.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		parsed, err := ExtractInfo(r)
		if err != nil {
			return nil, err
		}

		names := make([]string, 0, len(parsed))
		for _, object := range parsed {
			if object.Subdir != "" {
				names = append(names, object.Subdir)
			} else {
				names = append(names, object.Name)
			}
		}

		return names, nil
	case strings.HasPrefix(ct, "text/plain"):
		names := make([]string, 0, 50)

		body := string(r.(ObjectPage).Body.([]uint8))
		for _, name := range strings.Split(body, "\n") {
			if len(name) > 0 {
				names = append(names, name)
			}
		}

		return names, nil
	case strings.HasPrefix(ct, "text/html"):
		return []string{}, nil
	default:
		return nil, fmt.Errorf("cannot extract names from response with content-type: [%s]", ct)
	}
}

// DownloadHeader represents the headers returned in the response from a
// Download request.
type DownloadHeader struct {
	AcceptRanges       string    `json:"Accept-Ranges"`
	ContentDisposition string    `json:"Content-Disposition"`
	ContentEncoding    string    `json:"Content-Encoding"`
	ContentLength      int64     `json:"Content-Length,string"`
	ContentType        string    `json:"Content-Type"`
	Date               time.Time `json:"-"`
	DeleteAt           time.Time `json:"-"`
	ETag               string    `json:"Etag"`
	LastModified       time.Time `json:"-"`
	ObjectManifest     string    `json:"X-Object-Manifest"`
	StaticLargeObject  bool      `json:"-"`
	TransID            string    `json:"X-Trans-Id"`
	ObjectVersionID    string    `json:"X-Object-Version-Id"`
}

func (r *DownloadHeader) UnmarshalJSON(b []byte) error {
	type tmp DownloadHeader
	var s struct {
		tmp
		Date              gophercloud.JSONRFC1123 `json:"Date"`
		DeleteAt          gophercloud.JSONUnix    `json:"X-Delete-At"`
		LastModified      gophercloud.JSONRFC1123 `json:"Last-Modified"`
		StaticLargeObject any                     `json:"X-Static-Large-Object"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = DownloadHeader(s.tmp)

	switch t := s.StaticLargeObject.(type) {
	case string:
		if t == "True" || t == "true" {
			r.StaticLargeObject = true
		}
	case bool:
		r.StaticLargeObject = t
	}

	r.Date = time.Time(s.Date)
	r.DeleteAt = time.Time(s.DeleteAt)
	r.LastModified = time.Time(s.LastModified)

	return nil
}

// DownloadResult is a *http.Response that is returned from a call to the
// Download function.
type DownloadResult struct {
	gophercloud.HeaderResult
	Body io.ReadCloser
}

// Extract will return a struct of headers returned from a call to Download.
func (r DownloadResult) Extract() (*DownloadHeader, error) {
	var s DownloadHeader
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractContent is a function that takes a DownloadResult's io.Reader body
// and reads all available data into a slice of bytes. Please be aware that due
// the nature of io.Reader is forward-only - meaning that it can only be read
// once and not rewound. You can recreate a reader from the output of this
// function by using bytes.NewReader(downloadBytes)
func (r *DownloadResult) ExtractContent() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetHeader represents the headers returned in the response from a Get request.
type GetHeader struct {
	ContentDisposition string    `json:"Content-Disposition"`
	ContentEncoding    string    `json:"Content-Encoding"`
	ContentLength      int64     `json:"Content-Length,string"`
	ContentType        string    `json:"Content-Type"`
	Date               time.Time `json:"-"`
	DeleteAt           time.Time `json:"-"`
	ETag               string    `json:"Etag"`
	LastModified       time.Time `json:"-"`
	ObjectManifest     string    `json:"X-Object-Manifest"`
	StaticLargeObject  bool      `json:"-"`
	TransID            string    `json:"X-Trans-Id"`
	ObjectVersionID    string    `json:"X-Object-Version-Id"`
}

func (r *GetHeader) UnmarshalJSON(b []byte) error {
	type tmp GetHeader
	var s struct {
		tmp
		Date              gophercloud.JSONRFC1123 `json:"Date"`
		DeleteAt          gophercloud.JSONUnix    `json:"X-Delete-At"`
		LastModified      gophercloud.JSONRFC1123 `json:"Last-Modified"`
		StaticLargeObject any                     `json:"X-Static-Large-Object"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = GetHeader(s.tmp)

	switch t := s.StaticLargeObject.(type) {
	case string:
		if t == "True" || t == "true" {
			r.StaticLargeObject = true
		}
	case bool:
		r.StaticLargeObject = t
	}

	r.Date = time.Time(s.Date)
	r.DeleteAt = time.Time(s.DeleteAt)
	r.LastModified = time.Time(s.LastModified)

	return nil
}

// GetResult is a *http.Response that is returned from a call to the Get
// function.
type GetResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Get.
func (r GetResult) Extract() (*GetHeader, error) {
	var s GetHeader
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractMetadata is a function that takes a GetResult (of type *http.Response)
// and returns the custom metadata associated with the object.
func (r GetResult) ExtractMetadata() (map[string]string, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	metadata := make(map[string]string)
	for k, v := range r.Header {
		if strings.HasPrefix(k, "X-Object-Meta-") {
			key := strings.TrimPrefix(k, "X-Object-Meta-")
			metadata[key] = v[0]
		}
	}
	return metadata, nil
}

// CreateHeader represents the headers returned in the response from a
// Create request.
type CreateHeader struct {
	ContentLength   int64     `json:"Content-Length,string"`
	ContentType     string    `json:"Content-Type"`
	Date            time.Time `json:"-"`
	ETag            string    `json:"Etag"`
	LastModified    time.Time `json:"-"`
	TransID         string    `json:"X-Trans-Id"`
	ObjectVersionID string    `json:"X-Object-Version-Id"`
}

func (r *CreateHeader) UnmarshalJSON(b []byte) error {
	type tmp CreateHeader
	var s struct {
		tmp
		Date         gophercloud.JSONRFC1123 `json:"Date"`
		LastModified gophercloud.JSONRFC1123 `json:"Last-Modified"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = CreateHeader(s.tmp)

	r.Date = time.Time(s.Date)
	r.LastModified = time.Time(s.LastModified)

	return nil
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Create.
func (r CreateResult) Extract() (*CreateHeader, error) {
	var s CreateHeader
	err := r.ExtractInto(&s)
	return &s, err
}

// UpdateHeader represents the headers returned in the response from a
// Update request.
type UpdateHeader struct {
	ContentLength   int64     `json:"Content-Length,string"`
	ContentType     string    `json:"Content-Type"`
	Date            time.Time `json:"-"`
	TransID         string    `json:"X-Trans-Id"`
	ObjectVersionID string    `json:"X-Object-Version-Id"`
}

func (r *UpdateHeader) UnmarshalJSON(b []byte) error {
	type tmp UpdateHeader
	var s struct {
		tmp
		Date gophercloud.JSONRFC1123 `json:"Date"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = UpdateHeader(s.tmp)

	r.Date = time.Time(s.Date)

	return nil
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Update.
func (r UpdateResult) Extract() (*UpdateHeader, error) {
	var s UpdateHeader
	err := r.ExtractInto(&s)
	return &s, err
}

// DeleteHeader represents the headers returned in the response from a
// Delete request.
type DeleteHeader struct {
	ContentLength          int64     `json:"Content-Length,string"`
	ContentType            string    `json:"Content-Type"`
	Date                   time.Time `json:"-"`
	TransID                string    `json:"X-Trans-Id"`
	ObjectVersionID        string    `json:"X-Object-Version-Id"`
	ObjectCurrentVersionID string    `json:"X-Object-Current-Version-Id"`
}

func (r *DeleteHeader) UnmarshalJSON(b []byte) error {
	type tmp DeleteHeader
	var s struct {
		tmp
		Date gophercloud.JSONRFC1123 `json:"Date"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = DeleteHeader(s.tmp)

	r.Date = time.Time(s.Date)

	return nil
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Delete.
func (r DeleteResult) Extract() (*DeleteHeader, error) {
	var s DeleteHeader
	err := r.ExtractInto(&s)
	return &s, err
}

// CopyHeader represents the headers returned in the response from a
// Copy request.
type CopyHeader struct {
	ContentLength          int64     `json:"Content-Length,string"`
	ContentType            string    `json:"Content-Type"`
	CopiedFrom             string    `json:"X-Copied-From"`
	CopiedFromLastModified time.Time `json:"-"`
	Date                   time.Time `json:"-"`
	ETag                   string    `json:"Etag"`
	LastModified           time.Time `json:"-"`
	TransID                string    `json:"X-Trans-Id"`
	ObjectVersionID        string    `json:"X-Object-Version-Id"`
}

func (r *CopyHeader) UnmarshalJSON(b []byte) error {
	type tmp CopyHeader
	var s struct {
		tmp
		CopiedFromLastModified gophercloud.JSONRFC1123 `json:"X-Copied-From-Last-Modified"`
		Date                   gophercloud.JSONRFC1123 `json:"Date"`
		LastModified           gophercloud.JSONRFC1123 `json:"Last-Modified"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = CopyHeader(s.tmp)

	r.Date = time.Time(s.Date)
	r.CopiedFromLastModified = time.Time(s.CopiedFromLastModified)
	r.LastModified = time.Time(s.LastModified)

	return nil
}

// CopyResult represents the result of a copy operation.
type CopyResult struct {
	gophercloud.HeaderResult
}

// Extract will return a struct of headers returned from a call to Copy.
func (r CopyResult) Extract() (*CopyHeader, error) {
	var s CopyHeader
	err := r.ExtractInto(&s)
	return &s, err
}

type BulkDeleteResponse struct {
	ResponseStatus string     `json:"Response Status"`
	ResponseBody   string     `json:"Response Body"`
	Errors         [][]string `json:"Errors"`
	NumberDeleted  int        `json:"Number Deleted"`
	NumberNotFound int        `json:"Number Not Found"`
}

// BulkDeleteResult represents the result of a bulk delete operation. To extract
// the response object from the HTTP response, call its Extract method.
type BulkDeleteResult struct {
	gophercloud.Result
}

// Extract will return a BulkDeleteResponse struct returned from a BulkDelete
// call.
func (r BulkDeleteResult) Extract() (*BulkDeleteResponse, error) {
	var s BulkDeleteResponse
	err := r.ExtractInto(&s)
	return &s, err
}

// extractLastMarker is a function that takes a page of objects and returns the
// marker for the page. This can either be a subdir or the last object's name.
func extractLastMarker(r pagination.Page) (string, error) {
	casted := r.(ObjectPage)

	// If a delimiter was requested, check if a subdir exists.
	queryParams, err := url.ParseQuery(casted.RawQuery)
	if err != nil {
		return "", err
	}

	var delimeter bool
	if v, ok := queryParams["delimiter"]; ok && len(v) > 0 {
		delimeter = true
	}

	ct := casted.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		parsed, err := ExtractInfo(r)
		if err != nil {
			return "", err
		}

		var lastObject Object
		if len(parsed) > 0 {
			lastObject = parsed[len(parsed)-1]
		}

		if !delimeter {
			return lastObject.Name, nil
		}

		if lastObject.Name != "" {
			return lastObject.Name, nil
		}

		return lastObject.Subdir, nil
	case strings.HasPrefix(ct, "text/plain"):
		names := make([]string, 0, 50)

		body := string(r.(ObjectPage).Body.([]uint8))
		for _, name := range strings.Split(body, "\n") {
			if len(name) > 0 {
				names = append(names, name)
			}
		}

		return names[len(names)-1], err
	case strings.HasPrefix(ct, "text/html"):
		return "", nil
	default:
		return "", fmt.Errorf("cannot extract names from response with content-type: [%s]", ct)
	}
}
