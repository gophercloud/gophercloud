package images

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Image represents an image found in the OpenStack Image service.
type Image struct {
	// ID is the image UUID.
	ID string `json:"id"`

	// Name is the human-readable display name for the image.
	Name string `json:"name"`

	// Status is the image status. It can be "queued" or "active"
	// See image/v2/images/type.go
	Status ImageStatus `json:"status"`

	// Tags is a list of image tags. Tags are arbitrarily defined strings
	// attached to an image.
	Tags []string `json:"tags"`

	// ContainerFormat is the format of the container.
	// Valid values are ami, ari, aki, bare, and ovf.
	ContainerFormat string `json:"container_format"`

	// DiskFormat is the format of the disk.
	// If set, valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi,
	// and iso.
	DiskFormat string `json:"disk_format"`

	// MinDiskGigabytes is the amount of disk space in GB that is required to
	// boot the image.
	MinDiskGigabytes int `json:"min_disk"`

	// MinRAMMegabytes [optional] is the amount of RAM in MB that is required to
	// boot the image.
	MinRAMMegabytes int `json:"min_ram"`

	// Owner is the tenant ID the image belongs to.
	Owner string `json:"owner"`

	// Protected is whether the image is deletable or not.
	Protected bool `json:"protected"`

	// Visibility defines who can see/use the image.
	Visibility ImageVisibility `json:"visibility"`

	// Hidden is whether the image is listed in default image list or not.
	Hidden bool `json:"os_hidden"`

	// Checksum is the checksum of the data that's associated with the image.
	Checksum string `json:"checksum"`

	// SizeBytes is the size of the data that's associated with the image.
	SizeBytes int64 `json:"-"`

	// Metadata is a set of metadata associated with the image.
	// Image metadata allow for meaningfully define the image properties
	// and tags.
	// See http://docs.openstack.org/developer/glance/metadefs-concepts.html.
	Metadata map[string]string `json:"metadata"`

	// Properties is a set of key-value pairs, if any, that are associated with
	// the image.
	Properties map[string]any

	// CreatedAt is the date when the image has been created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the date when the last change has been made to the image or
	// its properties.
	UpdatedAt time.Time `json:"updated_at"`

	// File is the trailing path after the glance endpoint that represent the
	// location of the image or the path to retrieve it.
	File string `json:"file"`

	// Schema is the path to the JSON-schema that represent the image or image
	// entity.
	Schema string `json:"schema"`

	// VirtualSize is the virtual size of the image
	VirtualSize int64 `json:"virtual_size"`

	// OpenStackImageImportMethods is a slice listing the types of import
	// methods available in the cloud.
	OpenStackImageImportMethods []string `json:"-"`
	// OpenStackImageStoreIDs is a slice listing the store IDs available in
	// the cloud.
	OpenStackImageStoreIDs []string `json:"-"`
}

func (r *Image) UnmarshalJSON(b []byte) error {
	type tmp Image
	var s struct {
		tmp
		SizeBytes                   any    `json:"size"`
		OpenStackImageImportMethods string `json:"openstack-image-import-methods"`
		OpenStackImageStoreIDs      string `json:"openstack-image-store-ids"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Image(s.tmp)

	switch t := s.SizeBytes.(type) {
	case nil:
		r.SizeBytes = 0
	case float32:
		r.SizeBytes = int64(t)
	case float64:
		r.SizeBytes = int64(t)
	default:
		return fmt.Errorf("unknown type for SizeBytes: %v (value: %v)", reflect.TypeOf(t), t)
	}

	// Bundle all other fields into Properties
	var result any
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	if resultMap, ok := result.(map[string]any); ok {
		delete(resultMap, "self")
		delete(resultMap, "size")
		delete(resultMap, "openstack-image-import-methods")
		delete(resultMap, "openstack-image-store-ids")
		r.Properties = gophercloud.RemainingKeys(Image{}, resultMap)
	}

	if v := strings.FieldsFunc(strings.TrimSpace(s.OpenStackImageImportMethods), splitFunc); len(v) > 0 {
		r.OpenStackImageImportMethods = v
	}
	if v := strings.FieldsFunc(strings.TrimSpace(s.OpenStackImageStoreIDs), splitFunc); len(v) > 0 {
		r.OpenStackImageStoreIDs = v
	}

	return err
}

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any commonResult as an Image.
func (r commonResult) Extract() (*Image, error) {
	var s *Image
	if v, ok := r.Body.(map[string]any); ok {
		for k, h := range r.Header {
			if strings.ToLower(k) == "openstack-image-import-methods" {
				for _, s := range h {
					v["openstack-image-import-methods"] = s
				}
			}
			if strings.ToLower(k) == "openstack-image-store-ids" {
				for _, s := range h {
					v["openstack-image-store-ids"] = s
				}
			}
		}
	}
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult represents the result of a Create operation. Call its Extract
// method to interpret it as an Image.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of an Update operation. Call its Extract
// method to interpret it as an Image.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract
// method to interpret it as an Image.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a Delete operation. Call its
// ExtractErr method to interpret it as an Image.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ImagePage represents the results of a List request.
type ImagePage struct {
	serviceURL string
	pagination.LinkedPageBase
}

// IsEmpty returns true if an ImagePage contains no Images results.
func (r ImagePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	images, err := ExtractImages(r)
	return len(images) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to
// the next page of results.
func (r ImagePage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	if s.Next == "" {
		return "", nil
	}

	return nextPageURL(r.serviceURL, s.Next)
}

// ExtractImages interprets the results of a single page from a List() call,
// producing a slice of Image entities.
func ExtractImages(r pagination.Page) ([]Image, error) {
	var s struct {
		Images []Image `json:"images"`
	}
	err := (r.(ImagePage)).ExtractInto(&s)
	return s.Images, err
}

// splitFunc is a helper function used to avoid a slice of empty strings.
func splitFunc(c rune) bool {
	return c == ','
}
