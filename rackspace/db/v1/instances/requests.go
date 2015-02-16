package instances

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

func GetDefaultConfig(client *gophercloud.ServiceClient, id string) ConfigResult {
	var res ConfigResult

	resp, err := perigee.Request("GET", configURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

func AssociateWithConfigGroup(client *gophercloud.ServiceClient, instanceID, configGroupID string) UpdateResult {
	reqBody := map[string]string{
		"configuration": configGroupID,
	}

	var res UpdateResult

	resp, err := perigee.Request("PUT", resourceURL(client, instanceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     map[string]map[string]string{"instance": reqBody},
		OkCodes:     []int{202},
	})

	res.Header = resp.HttpResponse.Header
	res.Err = err

	return res
}

func ListBackups(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {

}
