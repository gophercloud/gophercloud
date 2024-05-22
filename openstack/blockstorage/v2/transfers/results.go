package transfers

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Transfer represents a Volume Transfer record
type Transfer struct {
	ID        string              `json:"id"`
	AuthKey   string              `json:"auth_key"`
	Name      string              `json:"name"`
	VolumeID  string              `json:"volume_id"`
	CreatedAt time.Time           `json:"-"`
	Links     []map[string]string `json:"links"`
}

// UnmarshalJSON is our unmarshalling helper
func (r *Transfer) UnmarshalJSON(b []byte) error {
	type tmp Transfer
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Transfer(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)

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

// ExtractInto converts our response data into a transfer struct
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
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Transfers.
func (r TransferPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	transfers, err := ExtractTransfers(r)
	return len(transfers) == 0, err
}

func (page TransferPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"transfers_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}
