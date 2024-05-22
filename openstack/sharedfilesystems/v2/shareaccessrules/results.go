package shareaccessrules

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// ShareAccess contains information associated with an OpenStack share access rule.
type ShareAccess struct {
	// The UUID of the share to which you are granted or denied access.
	ShareID string `json:"share_id"`
	// The date and time stamp when the resource was created within the service’s database.
	CreatedAt time.Time `json:"-"`
	// The date and time stamp when the resource was last updated within the service’s database.
	UpdatedAt time.Time `json:"-"`
	// The access rule type.
	AccessType string `json:"access_type"`
	// The value that defines the access. The back end grants or denies the access to it.
	AccessTo string `json:"access_to"`
	// The access credential of the entity granted share access.
	AccessKey string `json:"access_key"`
	// The state of the access rule.
	State string `json:"state"`
	// The access level to the share.
	AccessLevel string `json:"access_level"`
	// The access rule ID.
	ID string `json:"id"`
	// Access rule metadata.
	Metadata map[string]any `json:"metadata"`
}

func (r *ShareAccess) UnmarshalJSON(b []byte) error {
	type tmp ShareAccess
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ShareAccess(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	gophercloud.Result
}

// Extract will get the ShareAccess object from the GetResult.
func (r GetResult) Extract() (*ShareAccess, error) {
	var s struct {
		ShareAccess *ShareAccess `json:"access"`
	}
	err := r.ExtractInto(&s)
	return s.ShareAccess, err
}

// ListResult contains the response body and error from a List request.
type ListResult struct {
	gophercloud.Result
}

func (r ListResult) Extract() ([]ShareAccess, error) {
	var s struct {
		AccessList []ShareAccess `json:"access_list"`
	}
	err := r.ExtractInto(&s)
	return s.AccessList, err
}
