// Package attachments provides access to OpenStack Block Storage Attachment API's
// Use of this package requires Cinder version 3.45 at a minimum, and recommends
// 3.50 in order to utilize multi attach capability
package attachments

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Attachment contains all the information associated with an OpenStack Attachment.
type Attachment struct {
	// ID is the Unique identifier for the attachment.
	ID string `json:"id"`
	// VolumeID is the UUID of the Volume associated with this attachment
	VolumeID string `json:"volume_id"`
	// InstanceID is the Instance/Server UUID associated with this attachment
	ServerID string `json:"instance_uuid"`
	// Mountpoint is the requested Mountpoint of the volume on the Instance
	Mountpoint string `json:"mountpoint"`
	// AttachTime is the time the attachment was created
	AttachTime time.Time `json:"-"`
	// DetachTime is the time the attachment was created
	DetachTime time.Time `json:"-"`
	// AttachStatus is the current attach status
	AttachStatus string `json:"attach_status"`
	// AttachMode includes things like Read Only etc
	AttachMode string `json:"attach_mode"`
	// Add a volume sub-type here?
	// ConnectionInfo is the required info for a node to make a connection provided by the driver
	ConnectionInfo map[string]string `json:"connection_info"`
	// Connector is the initiator side of the connection info from the Node
	Connector map[string]string `json:"connector"`
}

// UnmarshalJSON is our unmarshalling helper
func (r *Attachment) UnmarshalJSON(b []byte) error {
	type tmp Attachment
	var s struct {
		tmp
		AttachedAt gophercloud.JSONRFC3339MilliNoZ `json:"attached_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Attachment(s.tmp)

	r.AttachTime = time.Time(s.AttachedAt)

	return err
}

// AttachmentPage is a pagination.pager that is returned from a call to the List function.
type AttachmentPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Attachments.
func (r AttachmentPage) IsEmpty() (bool, error) {
	attachments, err := ExtractAttachments(r)
	return len(attachments) == 0, err
}

// ExtractAttachments extracts and returns Attachments. It is used while iterating over a attachment.List call.
func ExtractAttachments(r pagination.Page) ([]Attachment, error) {
	var s []Attachment
	err := ExtractAttachmentsInto(r, &s)
	return s, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Attachment object out of the commonResult object.
func (r commonResult) Extract() (*Attachment, error) {
	var s Attachment
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a attachment struct
func (r commonResult) ExtractInto(a interface{}) error {
	return r.Result.ExtractIntoStructPtr(a, "attachment")
}

// ExtractAttachmentsInto similar to ExtractInto but operates on a `list` of attachments
func ExtractAttachmentsInto(r pagination.Page, a interface{}) error {
	return r.(AttachmentPage).Result.ExtractIntoSlicePtr(a, "attachments")
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}
