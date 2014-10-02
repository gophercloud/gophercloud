package tokens

import "github.com/rackspace/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tokens")
}
