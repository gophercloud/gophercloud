package volumetypes

import (
	"fmt"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOpts struct {
	ExtraSpecs map[string]interface{}
	Name       string
}

func Create(client *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
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

	var res CreateResult
	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &res.Resp,
	})
	return res
}

func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, err := perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	res.Err = err
	return res
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	resp, err := perigee.Request("GET", getURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		Results:     &res.Resp,
	})
	res.Err = err
	fmt.Printf("resp: %+v\n", resp)
	return res
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, listURL(client), createPage)
}
