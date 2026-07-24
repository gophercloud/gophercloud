package shareattach

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ShareAttachment contains attachment information between a share
// and server.
type ShareAttachment struct {
	// ShareID is the UUID of the attached Manila share.
	ShareID string `json:"share_id"`

	// Status is the status of the attached share.
	Status string `json:"status"`

	// Tag is a tag that can be applied to the attached share.
	Tag string `json:"tag"`

	// UUID is the UUID of the share attachment. This field is only visible to admins.
	UUID *string `json:"uuid"`

	// ExportLocation is the export location used to attach the share to the underlying host.
	// This field is only visible to admins.
	ExportLocation *string `json:"export_location"`
}

// ShareAttachmentPage stores a single page all of ShareAttachment
// results from a List call.
type ShareAttachmentPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a ShareAttachmentPage is empty.
func (page ShareAttachmentPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	sa, err := ExtractShareAttachments(page)
	return len(sa) == 0, err
}

// ExtractShareAttachments interprets a page of results as a slice of
// ShareAttachment.
func ExtractShareAttachments(r pagination.Page) ([]ShareAttachment, error) {
	var s struct {
		Shares []ShareAttachment `json:"shares"`
	}
	err := (r.(ShareAttachmentPage)).ExtractInto(&s)
	return s.Shares, err
}

// ShareAttachmentResult is the result from a share attachment operation.
type ShareAttachmentResult struct {
	gophercloud.Result
}

// Extract is a method that attempts to interpret any ShareAttachment resource
// response as a ShareAttachment struct.
func (r ShareAttachmentResult) Extract() (*ShareAttachment, error) {
	var s struct {
		Share *ShareAttachment `json:"share"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a ShareAttachment.
type CreateResult struct {
	ShareAttachmentResult
}

// GetResult is the response from a Get operation. Call its Extract method to
// interpret it as a ShareAttachment.
type GetResult struct {
	ShareAttachmentResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
