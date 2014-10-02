package volumes

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOpts struct {
	Availability                     string
	Description                      string
	Metadata                         map[string]string
	Name                             string
	Size                             int
	SnapshotID, SourceVolID, ImageID string
	VolumeType                       string
}

func Create(client *gophercloud.ServiceClient, opts *CreateOpts) CreateResult {

	type volume struct {
		Availability *string           `json:"availability_zone,omitempty"`
		Description  *string           `json:"display_description,omitempty"`
		ImageID      *string           `json:"imageRef,omitempty"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		Name         *string           `json:"display_name,omitempty"`
		Size         *int              `json:"size,omitempty"`
		SnapshotID   *string           `json:"snapshot_id,omitempty"`
		SourceVolID  *string           `json:"source_volid,omitempty"`
		VolumeType   *string           `json:"volume_type,omitempty"`
	}

	type request struct {
		Volume volume `json:"volume"`
	}

	reqBody := request{
		Volume: volume{},
	}

	reqBody.Volume.Availability = utils.MaybeString(opts.Availability)
	reqBody.Volume.Description = utils.MaybeString(opts.Description)
	reqBody.Volume.ImageID = utils.MaybeString(opts.ImageID)
	reqBody.Volume.Name = utils.MaybeString(opts.Name)
	reqBody.Volume.Size = utils.MaybeInt(opts.Size)
	reqBody.Volume.SnapshotID = utils.MaybeString(opts.SnapshotID)
	reqBody.Volume.SourceVolID = utils.MaybeString(opts.SourceVolID)
	reqBody.Volume.VolumeType = utils.MaybeString(opts.VolumeType)

	var res CreateResult
	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200, 201},
	})
	return res
}

// ListOpts holds options for listing volumes. It is passed to the volumes.List function.
type ListOpts struct {
	// AllTenants is an admin-only option. Set it to true to see a tenant volumes.
	AllTenants bool
	// List only volumes that contain Metadata.
	Metadata map[string]string
	// List only volumes that have Name as the display name.
	Name string
	// List only volumes that have a status of Status.
	Status string
}

func List(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, listURL(client), createPage)
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		Results:     &res.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})
	return res
}

type UpdateOpts struct {
	Name        string
	Description string
	Metadata    map[string]string
}

func Update(client *gophercloud.ServiceClient, id string, opts *UpdateOpts) UpdateResult {
	type update struct {
		Description *string           `json:"display_description,omitempty"`
		Metadata    map[string]string `json:"metadata,omitempty"`
		Name        *string           `json:"display_name,omitempty"`
	}

	type request struct {
		Volume update `json:"volume"`
	}

	reqBody := request{
		Volume: update{},
	}

	reqBody.Volume.Description = utils.MaybeString(opts.Description)
	reqBody.Volume.Name = utils.MaybeString(opts.Name)

	var res UpdateResult

	_, res.Err = perigee.Request("PUT", updateURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &res.Resp,
	})
	return res

}

func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
