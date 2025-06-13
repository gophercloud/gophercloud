package images

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/assert"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Image represents an image found in the OpenStack Image service.
type Image struct {
	// ID is the image UUID.
	ID string `json:"id"`

	// Name is the human-readable display name for the image.
	Name string `json:"-"`

	// Status is the image status. It can be "queued" or "active"
	// See image/v2/images/type.go
	Status ImageStatus `json:"status"`

	// Tags is a list of image tags. Tags are arbitrarily defined strings
	// attached to an image.
	Tags []string `json:"tags"`

	// ContainerFormat is the format of the container.
	// Valid values are ami, ari, aki, bare, and ovf.
	ContainerFormat string `json:"-"`

	// DiskFormat is the format of the disk.
	// If set, valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi,
	// and iso.
	DiskFormat string `json:"-"`

	// MinDiskGigabytes is the amount of disk space in GB that is required to
	// boot the image.
	MinDiskGigabytes int `json:"-"`

	// MinRAMMegabytes [optional] is the amount of RAM in MB that is required to
	// boot the image.
	MinRAMMegabytes int `json:"-"`

	// Owner is the tenant ID the image belongs to.
	Owner string `json:"-"`

	// Protected is whether the image is deletable or not.
	Protected bool `json:"protected"`

	// Visibility defines who can see/use the image.
	Visibility ImageVisibility `json:"visibility"`

	// Hidden is whether the image is listed in default image list or not.
	Hidden bool `json:"-"`

	// Checksum is the checksum of the data that's associated with the image.
	Checksum string `json:"-"`

	// SizeBytes is the size of the data that's associated with the image.
	SizeBytes int64 `json:"-"`

	// Metadata is a set of metadata associated with the image.
	// Image metadata allow for meaningfully define the image properties
	// and tags.
	// See http://docs.openstack.org/developer/glance/metadefs-concepts.html.
	Metadata map[string]string `json:"metadata"`

	// Properties is a set of key-value pairs, if any, that are associated with
	// the image.
	Properties map[string]any `json:"-"`

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
	VirtualSize int64 `json:"-"`

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
		Name                        any    `json:"name"`
		ContainerFormat             any    `json:"container_format"`
		DiskFormat                  any    `json:"disk_format"`
		MinDiskGigabytes            any    `json:"min_disk"`
		MinRAMMegabytes             any    `json:"min_ram"`
		Owner                       any    `json:"owner"`
		Hidden                      any    `json:"os_hidden"`
		Checksum                    any    `json:"checksum"`
		SizeBytes                   any    `json:"size"`
		VirtualSize                 any    `json:"virtual_size"`
		OpenStackImageImportMethods string `json:"openstack-image-import-methods"`
		OpenStackImageStoreIDs      string `json:"openstack-image-store-ids"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Image(s.tmp)

	// some of the values may be nil, or bool as a string format, or int as
	// a string format so we need to assert them to the correct type using
	// the assert package
	r.Name, err = assert.String(s.Name, "Name")
	if err != nil {
		return err
	}
	r.ContainerFormat, err = assert.String(s.ContainerFormat, "ContainerFormat")
	if err != nil {
		return err
	}
	r.DiskFormat, err = assert.String(s.DiskFormat, "DiskFormat")
	if err != nil {
		return err
	}
	r.MinDiskGigabytes, err = assert.Int(s.MinDiskGigabytes, "MinDiskGigabytes")
	if err != nil {
		return err
	}
	r.MinRAMMegabytes, err = assert.Int(s.MinRAMMegabytes, "MinRAMMegabytes")
	if err != nil {
		return err
	}
	r.Owner, err = assert.String(s.Owner, "Owner")
	if err != nil {
		return err
	}
	r.Hidden, err = assert.Bool(s.Hidden, "Hidden")
	if err != nil {
		return err
	}
	r.Checksum, err = assert.String(s.Checksum, "Checksum")
	if err != nil {
		return err
	}
	r.SizeBytes, err = assert.Int64(s.SizeBytes, "SizeBytes")
	if err != nil {
		return err
	}
	r.VirtualSize, err = assert.Int64(s.VirtualSize, "VirtualSize")
	if err != nil {
		return err
	}

	// Bundle all other fields into Properties
	var result any
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	if resultMap, ok := result.(map[string]any); ok {
		delete(resultMap, "self")
		delete(resultMap, "name")
		delete(resultMap, "container_format")
		delete(resultMap, "disk_format")
		delete(resultMap, "min_disk")
		delete(resultMap, "min_ram")
		delete(resultMap, "owner")
		delete(resultMap, "os_hidden")
		delete(resultMap, "checksum")
		delete(resultMap, "size")
		delete(resultMap, "virtual_size")
		delete(resultMap, "openstack-image-import-methods")
		delete(resultMap, "openstack-image-store-ids")
		r.Properties = gophercloud.RemainingKeys(Image{}, resultMap)
		if m, ok := resultMap["properties"]; ok {
			r.Properties["properties"] = m
		}
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
