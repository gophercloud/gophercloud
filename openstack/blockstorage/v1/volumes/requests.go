package volumes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVolumeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume. This object is passed to
// the volumes.Create function. For more information about these parameters,
// see the Volume object.
type CreateOpts struct {
	// OPTIONAL
	Availability string
	// OPTIONAL
	Description string
	// OPTIONAL
	Metadata map[string]string
	// OPTIONAL
	Name string
	// REQUIRED
	Size int
	// OPTIONAL
	SnapshotID, SourceVolID, ImageID string
	// OPTIONAL
	VolumeType string
}

// ToVolumeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToVolumeCreateMap() (map[string]interface{}, error) {
	v := make(map[string]interface{})

	if opts.Size == 0 {
		err := &gophercloud.ErrMissingInput{}
		err.Argument = "CreateOpts.Size"
		err.Function = "volumes.ToVolumeCreateMap"
		return nil, err
	}
	v["size"] = opts.Size

	if opts.Availability != "" {
		v["availability_zone"] = opts.Availability
	}
	if opts.Description != "" {
		v["display_description"] = opts.Description
	}
	if opts.ImageID != "" {
		v["imageRef"] = opts.ImageID
	}
	if opts.Metadata != nil {
		v["metadata"] = opts.Metadata
	}
	if opts.Name != "" {
		v["display_name"] = opts.Name
	}
	if opts.SourceVolID != "" {
		v["source_volid"] = opts.SourceVolID
	}
	if opts.SnapshotID != "" {
		v["snapshot_id"] = opts.SnapshotID
	}
	if opts.VolumeType != "" {
		v["volume_type"] = opts.VolumeType
	}

	return map[string]interface{}{"volume": v}, nil
}

// Create will create a new Volume based on the values in CreateOpts. To extract
// the Volume object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToVolumeCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(createURL(client), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return res
}

// Delete will delete the existing Volume with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = client.Delete(deleteURL(client, id), nil)
	return res
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getURL(client, id), &res.Body, nil)
	return res
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeListQuery() (string, error)
}

// ListOpts holds options for listing Volumes. It is passed to the volumes.List
// function.
type ListOpts struct {
	// admin-only option. Set it to true to see all tenant volumes.
	AllTenants bool `q:"all_tenants"`
	// List only volumes that contain Metadata.
	Metadata map[string]string `q:"metadata"`
	// List only volumes that have Name as the display name.
	Name string `q:"name"`
	// List only volumes that have a status of Status.
	Status string `q:"status"`
}

// ToVolumeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns Volumes optionally limited by the conditions provided in ListOpts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	createPage := func(r pagination.PageResult) pagination.Page {
		return VolumePage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPage)
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVolumeUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Volume. This object is passed
// to the volumes.Update function. For more information about the parameters, see
// the Volume object.
type UpdateOpts struct {
	// OPTIONAL
	Name string
	// OPTIONAL
	Description string
	// OPTIONAL
	Metadata map[string]string
}

// ToVolumeUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToVolumeUpdateMap() (map[string]interface{}, error) {
	v := make(map[string]interface{})

	if opts.Description != "" {
		v["display_description"] = opts.Description
	}
	if opts.Metadata != nil {
		v["metadata"] = opts.Metadata
	}
	if opts.Name != "" {
		v["display_name"] = opts.Name
	}

	return map[string]interface{}{"volume": v}, nil
}

// Update will update the Volume with provided information. To extract the updated
// Volume from the response, call the Extract method on the UpdateResult.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToVolumeUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Put(updateURL(client, id), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// IDFromName is a convienience function that returns a server's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""
	if name == "" {
		err := &gophercloud.ErrMissingInput{}
		err.Function = "volumes.IDFromName"
		err.Argument = "name"
		return "", err
	}

	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractVolumes(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		err := &gophercloud.ErrResourceNotFound{}
		err.Name = name
		err.ResourceType = "volume"
		err.Function = "volumes.IDFromName"
		return "", err
	case 1:
		return id, nil
	default:
		err := &gophercloud.ErrMultipleResourcesFound{}
		err.Count = count
		err.Name = name
		err.ResourceType = "volume"
		err.Function = "volumes.IDFromName"
		return "", err
	}
}
