package instances

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
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
