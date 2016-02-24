package volumeattach

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// VolumeAttachment controls the attachment of a volume to an instance.
type VolumeAttachment struct {
	// ID is a unique id of the attachment
	ID string `json:"id"`

	// Device is what device the volume is attached as
	Device string `json:"device"`

	// VolumeID is the ID of the attached volume
	VolumeID string `json:"volumeId"`

	// ServerID is the ID of the instance that has the volume attached
	ServerID string `json:"serverId"`
}

// VolumeAttachmentsPage stores a single, only page of VolumeAttachments
// results from a List call.
type VolumeAttachmentsPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a VolumeAttachmentsPage is empty.
func (page VolumeAttachmentsPage) IsEmpty() (bool, error) {
	va, err := ExtractVolumeAttachments(page)
	return len(va) == 0, err
}

// ExtractVolumeAttachments interprets a page of results as a slice of
// VolumeAttachments.
func ExtractVolumeAttachments(page pagination.Page) ([]VolumeAttachment, error) {
	r := page.(VolumeAttachmentsPage)
	var s struct {
		VolumeAttachments []VolumeAttachment `json:"volumeAttachments"`
	}
	err := r.ExtractInto(&s)
	return s.VolumeAttachments, err
}

type VolumeAttachmentResult struct {
	gophercloud.Result
}

// Extract is a method that attempts to interpret any VolumeAttachment resource
// response as a VolumeAttachment struct.
func (r VolumeAttachmentResult) Extract() (*VolumeAttachment, error) {
	var s struct {
		VolumeAttachment *VolumeAttachment `json:"volumeAttachment"`
	}
	err := r.ExtractInto(&s)
	return s.VolumeAttachment, err
}

// CreateResult is the response from a Create operation. Call its Extract method to interpret it
// as a VolumeAttachment.
type CreateResult struct {
	VolumeAttachmentResult
}

// GetResult is the response from a Get operation. Call its Extract method to interpret it
// as a VolumeAttachment.
type GetResult struct {
	VolumeAttachmentResult
}

// DeleteResult is the response from a Delete operation. Call its Extract method to determine if
// the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
