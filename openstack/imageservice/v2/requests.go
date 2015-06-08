package v2

import (
	"github.com/rackspace/gophercloud"
	//"github.com/rackspace/gophercloud/pagination"
)

// List : (*gophercloud.ServiceClient) -> pagination.Pager<ImagePage>
// func List(c *gophercloud.ServiceClient) pagination.Pager {
// 	return pagination.NewPager(c, listURL(c), func (r pagination.PageResult) pagination.Page {
// 		p := ImagePage{pagination.MarkerPageBase{PageResult: r}}
// 		p.MarkerPageBase.Owner = p
// 		return p
// 	})
// }

func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult
	body, err := opts.ToImageCreateMap()
	if err != nil {
		res.Err = err
		return res
	}
	response, err := client.Post(createURL(client), body, &res.Body, nil)
	// TODO check response..?
	return res
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
// The CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	// Returns value that implements json.Marshaler
	ToImageCreateMap() (map[string]interface{}, error)
}

// implements CreateOptsBuilder
type CreateOpts struct {
	// Name [required] is the name of the new image.
	Name string

	// Id [optional] is the the image ID.
	Id string

	// Visibility [optional] defines who can see/use the image.
	Visibility ImageVisibility

	// Tags [optional] is a set of image tags.
	Tags []string

	// ContainerFormat [optional] is the format of the
	// container. Valid values are ami, ari, aki, bare, and ovf.
	ContainerFormat string

	// DiskFormat [optional] is the format of the disk. If set,
	// valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi,
	// and iso.
	DiskFormat string

	// MinDiskGigabytes [optional] is the amount of disk space in
	// GB that is required to boot the image.
	MinDiskGigabytes int

	// MinRamMegabytes [optional] is the amount of RAM in MB that
	// is required to boot the image.
	MinRamMegabytes int

	// protected [optional] is whether the image is not deletable.
	Protected bool

	// properties [optional] is a set of properties, if any, that
	// are associated with the image.
	Properties map[string]string
}

func setIfNotNil(m map[string]interface{}, k string, v interface{}) {
	if v != nil {
		m[k] = v
	}
}

// ToImageCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToImageCreateMap() (map[string]interface{}, error) {
	body := map[string]interface{}{}
	body["name"] = opts.Name
	setIfNotNil(body, "id", opts.Id)
	setIfNotNil(body, "visibility", opts.Visibility)
	setIfNotNil(body, "tags", opts.Tags)
	setIfNotNil(body, "container_format", opts.ContainerFormat)
	setIfNotNil(body, "disk_format", opts.DiskFormat)
	setIfNotNil(body, "min_disk", opts.MinDiskGigabytes)
	setIfNotNil(body, "min_ram", opts.MinRamMegabytes)
	setIfNotNil(body, "protected", opts.Protected)
	setIfNotNil(body, "properties", opts.Properties)
	return body, nil
}

func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = client.Delete(deleteURL(client, id), nil)
	return res
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	client.Get(getURL(client, id), &res.Body, nil)
	return res
}

func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult
	reqBody := opts.ToImageUpdateMap()

	_, res.Err = client.Patch(updateURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

type UpdateOptsBuilder interface {
	// returns value implementing json.Marshaler which when marshaled matches the patch schema:
	// http://specs.openstack.org/openstack/glance-specs/specs/api/v2/http-patch-image-api-v2.html
	ToImageUpdateMap() []interface{}
}

// implements UpdateOptsBuilder
type UpdateOpts []Patch

func (opts UpdateOpts) ToImageUpdateMap() []interface{} {
	m := make([]interface{}, len(opts))
	for i, patch := range opts {
		patchJson := patch.ToImagePatchMap()
		m[i] = patchJson
	}
	return m
}

// Patch represents a single update to an existing image. Multiple updates to an image can be
// submitted at the same time.
type Patch interface {
	ToImagePatchMap() map[string]interface{}
}

//implements Patch
type ReplaceImageName struct {
	NewName string
}

func (r *ReplaceImageName) ToImagePatchMap() map[string]interface{} {
	m := map[string]interface{}{}
	m["op"] = "replace"
	m["path"] = "/name"
	m["value"] = r.NewName
	return m
}

// Things implementing Patch can also implement UpdateOptsBuilder.
// Unfortunately we have to specify each of these manually.

func (r *ReplaceImageName) ToImageUpdateMap() map[string]interface{} {
	return r.ToImagePatchMap()
}
