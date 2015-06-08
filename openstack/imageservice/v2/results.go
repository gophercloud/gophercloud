package v2

import (
	"github.com/rackspace/gophercloud"
	//"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

// does not include the literal image data; just metadata.
// returned by listing images, and by fetching a specific image.
type Image struct {
	Id string `mapstructure:"id"`
	
	Name string `mapstructure:"name"`
	
	Status ImageStatus `mapstructure:"status"`
	
	Tags []string `mapstructure:"tags"`
	
	ContainerFormat string `mapstructure:"container_format"`
	DiskFormat string `mapstructure:"disk_format"`
	
	MinDiskGigabytes int `mapstructure:"min_disk"`
	MinRamMegabytes int `mapstructure:"min_ram"`
	
	Owner string `mapstructure:"owner"`
	
	Protected bool `mapstructure:"protected"`
	Visibility ImageVisibility `mapstructure:"visibility"`

	Checksum string `mapstructure:"checksum"`
	SizeBytes int `mapstructure:"size"`
	
	Metadata map[string]string `mapstructure:"metadata"`
	Properties map[string]string `mapstructure:"properties"`
}

// implements pagination.Page<Image>, pagination.MarkerPage
// DOESN'T implement Page. Why? How does openstack/compute/v2
// type ImagePage struct {
// 	pagination.MarkerPageBase  // pagination.MarkerPageBase<Image>
// }

type CreateResult struct {
	gophercloud.ErrResult
}

// The response to `POST /images` follows the same schema as `GET /images/:id`.
func extractImage(res gophercloud.ErrResult) (*Image, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	var image Image

	err := mapstructure.Decode(res.Body, &image)
	
	return &image, err
}

func (c CreateResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type GetResult struct {
	gophercloud.ErrResult
}

func (c GetResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

type UpdateResult struct {
	gophercloud.ErrResult
}
