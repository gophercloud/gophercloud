package volumes

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/racker/perigee"
)

// CreateOpts contains options for creating a Volume. This object is passed to
// the volumes.Create function. For more information about these parameters,
// see the Volume object.
type CreateOpts struct {
	Availability                     string            // OPTIONAL
	Description                      string            // OPTIONAL
	Metadata                         map[string]string // OPTIONAL
	Name                             string            // OPTIONAL
	Size                             int               // REQUIRED
	SnapshotID, SourceVolID, ImageID string            // REQUIRED (one of them)
	VolumeType                       string            // OPTIONAL
}

// Create will create a new Volume based on the values in CreateOpts. To extract
// the Volume object from the response, call the Extract method on the CreateResult.
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

	reqBody.Volume.Availability = gophercloud.MaybeString(opts.Availability)
	reqBody.Volume.Description = gophercloud.MaybeString(opts.Description)
	reqBody.Volume.ImageID = gophercloud.MaybeString(opts.ImageID)
	reqBody.Volume.Name = gophercloud.MaybeString(opts.Name)
	reqBody.Volume.Size = gophercloud.MaybeInt(opts.Size)
	reqBody.Volume.SnapshotID = gophercloud.MaybeString(opts.SnapshotID)
	reqBody.Volume.SourceVolID = gophercloud.MaybeString(opts.SourceVolID)
	reqBody.Volume.VolumeType = gophercloud.MaybeString(opts.VolumeType)

	var res CreateResult
	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200, 201},
	})
	return res
}

// Delete will delete the existing Volume with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202, 204},
	})
	return err
}

// Get retrieves the Volume with the provided ID. To extract the Volume object from
// the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		Results:     &res.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})
	return res
}

// ListOpts holds options for listing Volumes. It is passed to the volumes.List
// function.
type ListOpts struct {
	AllTenants bool              `q:"all_tenants"` // admin-only option. Set it to true to see all tenant volumes.
	Metadata   map[string]string `q:"metadata"`    // List only volumes that contain Metadata.
	Name       string            `q:"name"`        // List only volumes that have Name as the display name.
	Status     string            `q:"status"`      // List only volumes that have a status of Status.
}

// List returns Volumes optionally limited by the conditions provided in ListOpts.
func List(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query.String()
	}
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, listURL(client), createPage)
}

// UpdateOpts contain options for updating an existing Volume. This object is passed
// to the volumes.Update function. For more information about the parameters, see
// the Volume object.
type UpdateOpts struct {
	Name        string            // OPTIONAL
	Description string            // OPTIONAL
	Metadata    map[string]string // OPTIONAL
}

// Update will update the Volume with provided information. To extract the updated
// Volume from the response, call the Extract method on the UpdateResult.
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

	reqBody.Volume.Description = gophercloud.MaybeString(opts.Description)
	reqBody.Volume.Name = gophercloud.MaybeString(opts.Name)

	var res UpdateResult

	_, res.Err = perigee.Request("PUT", updateURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &res.Resp,
	})
	return res
}
