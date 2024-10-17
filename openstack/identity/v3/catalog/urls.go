package catalog

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("auth", "catalog")
}
