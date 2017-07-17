package volumeactions

import (
	"time"

	"github.com/gophercloud/gophercloud"
)

// AttachResult contains the response body and error from a Get request.
type AttachResult struct {
	gophercloud.ErrResult
}

// BeginDetachingResult contains the response body and error from a Get request.
type BeginDetachingResult struct {
	gophercloud.ErrResult
}

// DetachResult contains the response body and error from a Get request.
type DetachResult struct {
	gophercloud.ErrResult
}

// UploadImageResult contains the response body and error from a UploadImage request.
type UploadImageResult struct {
	gophercloud.ErrResult
}

// ReserveResult contains the response body and error from a Get request.
type ReserveResult struct {
	gophercloud.ErrResult
}

// UnreserveResult contains the response body and error from a Get request.
type UnreserveResult struct {
	gophercloud.ErrResult
}

// TerminateConnectionResult contains the response body and error from a Get request.
type TerminateConnectionResult struct {
	gophercloud.ErrResult
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (map[string]interface{}, error) {
	var s struct {
		ConnectionInfo map[string]interface{} `json:"connection_info"`
	}
	err := r.ExtractInto(&s)
	return s.ConnectionInfo, err
}

// VolumeImage contains information about volume upload to an image service.
type VolumeImage struct {
	// The ID of a volume an image is created from.
	VolumeId string `json:"id"`
	// Container format, may be bare, ofv, ova, etc.
	ContainerFormat string `json:"container_format"`
	// Disk format, may be raw, qcow2, vhd, vdi, vmdk, etc.
	DiskFormat string `json:"disk_format"`
	// Human-readable description for the volume.
	Description string `json:"display_description"`
	// The ID of an image being created.
	ImageId string `json:"image_id"`
	// Human-readable display name for the image.
	ImageName string `json:"image_name"`
	// Size of the volume in GB.
	Size int `json:"size"`
	// Current status of the volume.
	Status string `json:"status"`
	// The date when this volume was last updated.
	UpdatedAt time.Time `json:"-"`
}

// Extract will get an object with info about image being uploaded out of the UploadImageResult object.
func (r UploadImageResult) Extract() (VolumeImage, error) {
	var s struct {
		VolumeUploadImage VolumeImage `json:"os-volume_upload_image"`
	}
	err := r.ExtractInto(&s)
	return s.VolumeUploadImage, err
}

// InitializeConnectionResult contains the response body and error from a Get request.
type InitializeConnectionResult struct {
	commonResult
}

// ExtendSizeResult contains the response body and error from an ExtendSize request.
type ExtendSizeResult struct {
	gophercloud.ErrResult
}
