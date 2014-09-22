package volumeTypes

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOpts struct {
	ExtraSpecs map[string]interface{}
	Name       string
}

func Create(client *gophercloud.ServiceClient, opts CreateOpts) (*VolumeType, error) {
	type volumeType struct {
		ExtraSpecs map[string]interface{} `json:"extra_specs,omitempty"`
		Name       *string                `json:"name,omitempty"`
	}

	type request struct {
		VolumeType volumeType `json:"volume_type"`
	}

	reqBody := request{
		VolumeType: volumeType{},
	}

	reqBody.VolumeType.Name = utils.MaybeString(opts.Name)
	reqBody.VolumeType.ExtraSpecs = opts.ExtraSpecs

	type response struct {
		VolumeType VolumeType `json:"volume_type"`
	}

	var respBody response

	_, err := perigee.Request("POST", volumeTypesURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &respBody,
	})
	if err != nil {
		return nil, err
	}

	return &respBody.VolumeType, nil

}

func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", volumeTypeURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return err
}

func Get(client *gophercloud.ServiceClient, id string) (GetResult, error) {
	var gr GetResult
	_, err := perigee.Request("GET", volumeTypeURL(client, id), perigee.Options{
		Results:     &gr,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return gr, err
}

func List(client *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, volumeTypesURL(client), createPage)
}
