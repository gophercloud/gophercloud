package flavors

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id int) string {
	return client.ServiceURL("flavors", strconv.Itoa(id))
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors")
}
