package volumeTypes

import (
	"fmt"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type CreateOpts struct {
	ExtraSpecs map[string]string
	Name       string
}

func Create(client *gophercloud.ServiceClient, opts CreateOpts) (*VolumeType, error) {
	type volumeType struct {
		ExtraSpecs map[string]string `json:"extra_specs,omitempty"`
		Name       *string           `json:"name,omitempty"`
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

	fmt.Printf("req: %+v\n", reqBody)
	fmt.Printf("res: %+v\n", respBody)

	return &respBody.VolumeType, nil

}
