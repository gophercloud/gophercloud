package sharetransfers

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

const (
	invalidMarker = "-1"
)

// Transfer represents a Share Transfer record.
type Transfer struct {
	ID                   string              `json:"id"`
	Accepted             bool                `json:"accepted"`
	AuthKey              string              `json:"auth_key"`
	Name                 string              `json:"name"`
	SourceProjectID      string              `json:"source_project_id"`
	DestinationProjectID string              `json:"destination_project_id"`
	ResourceID           string              `json:"resource_id"`
	ResourceType         string              `json:"resource_type"`
	CreatedAt            time.Time           `json:"-"`
	ExpiresAt            time.Time           `json:"-"`
	Links                []map[string]string `json:"links"`
}

// UnmarshalJSON is our unmarshalling helper.
func (r *Transfer) UnmarshalJSON(b []byte) error {
	type tmp Transfer
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		ExpiresAt gophercloud.JSONRFC3339MilliNoZ `json:"expires_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Transfer(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.ExpiresAt = time.Time(s.ExpiresAt)

	return err
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Transfer object out of the commonResult object.
func (r commonResult) Extract() (*Transfer, error) {
	var s Transfer
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a transfer struct.
func (r commonResult) ExtractInto(v any) error {
	return r.Result.ExtractIntoStructPtr(v, "transfer")
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AcceptResult contains the response body and error from an Accept request.
type AcceptResult struct {
	gophercloud.ErrResult
}

// ExtractTransfers extracts and returns Transfers. It is used while iterating over a transfers.List call.
func ExtractTransfers(r pagination.Page) ([]Transfer, error) {
	var s []Transfer
	err := ExtractTransfersInto(r, &s)
	return s, err
}

// ExtractTransfersInto similar to ExtractInto but operates on a `list` of transfers
func ExtractTransfersInto(r pagination.Page, v any) error {
	return r.(TransferPage).Result.ExtractIntoSlicePtr(v, "transfers")
}

// TransferPage is a pagination.pager that is returned from a call to the List function.
type TransferPage struct {
	pagination.MarkerPageBase
}

// NextPageURL generates the URL for the page of results after this one.
func (r TransferPage) NextPageURL() (string, error) {
	currentURL := r.URL
	mark, err := r.Owner.LastMarker()
	if err != nil {
		return "", err
	}
	if mark == invalidMarker {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("offset", mark)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// LastMarker returns the last offset in a ListResult.
func (r TransferPage) LastMarker() (string, error) {
	replicas, err := ExtractTransfers(r)
	if err != nil {
		return invalidMarker, err
	}
	if len(replicas) == 0 {
		return invalidMarker, nil
	}

	u, err := url.Parse(r.URL.String())
	if err != nil {
		return invalidMarker, err
	}
	queryParams := u.Query()
	offset := queryParams.Get("offset")
	limit := queryParams.Get("limit")

	// Limit is not present, only one page required
	if limit == "" {
		return invalidMarker, nil
	}

	iOffset := 0
	if offset != "" {
		iOffset, err = strconv.Atoi(offset)
		if err != nil {
			return invalidMarker, err
		}
	}
	iLimit, err := strconv.Atoi(limit)
	if err != nil {
		return invalidMarker, err
	}
	iOffset = iOffset + iLimit
	offset = strconv.Itoa(iOffset)

	return offset, nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface.
func (r TransferPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	replicas, err := ExtractTransfers(r)
	return len(replicas) == 0, err
}
