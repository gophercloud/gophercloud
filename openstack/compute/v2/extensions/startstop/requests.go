package startstop

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

func actionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

func Start(client *gophercloud.ServiceClient, id string) gophercloud.ErrResult {
	var res gophercloud.ErrResult

	reqBody := map[string](*interface{}){"os-start": nil}

	_, res.Err = perigee.Request("POST", actionURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})

	return res
}

func Stop(client *gophercloud.ServiceClient, id string) gophercloud.ErrResult {
	var res gophercloud.ErrResult

	reqBody := map[string](*interface{}){"os-stop": nil}

	_, res.Err = perigee.Request("POST", actionURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})

	return res
}
