package hosts

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("os-hosts")
}
