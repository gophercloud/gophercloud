package volumes

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type VolumeOpts struct {
	Availability                     string
	Description                      string
	Metadata                         map[string]string
	Name                             string
	Size                             int
	SnapshotID, SourceVolID, ImageID string
	VolumeType                       string
}

func Create(client *gophercloud.ServiceClient, opts VolumeOpts) (*Volume, error) {

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

/*
func Get(c *blockstorage.Client, opts GetOpts) (Volume, error) {
	var v Volume
	h, err := c.GetHeaders()
	if err != nil {
		return v, err
	}
	url := c.GetVolumeURL(opts["id"])
	_, err = perigee.Request("GET", url, perigee.Options{
		Results: &struct {
			Volume *Volume `json:"volume"`
		}{&v},
		MoreHeaders: h,
	})
	return v, err
}
*/

func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", volumeURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return err
}
