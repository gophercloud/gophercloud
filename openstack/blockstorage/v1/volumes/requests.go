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

func Create(client *gophercloud.ServiceClient, opts CreateOpts) (*Volume, error) {

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

	type response struct {
		Volume Volume `json:"volume"`
	}

	var respBody response

	_, err := perigee.Request("POST", volumesURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200, 201},
		ReqBody:     &reqBody,
		Results:     &respBody,
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Volume, nil
}

func List(client *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, volumesURL(client), createPage)
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var gr GetResult
	_, err := perigee.Request("GET", volumeURL(client, id), perigee.Options{
		Results:     &gr.r,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	gr.err = err
	return gr
}

type UpdateOpts struct {
	Name        string
	Description string
	Metadata    map[string]string
}

func Update(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (*Volume, error) {
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

	type response struct {
		Volume Volume `json:"volume"`
	}

	var respBody response

	_, err := perigee.Request("PUT", volumeURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &respBody,
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Volume, nil

}

func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", volumeURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return err
}
