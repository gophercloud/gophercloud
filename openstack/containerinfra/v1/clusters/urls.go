package clusters

import (
	"github.com/gophercloud/gophercloud/v2"
)

var apiName = "clusters"

func commonURL(client gophercloud.Client) string {
	return client.ServiceURL(apiName)
}

func idURL(client gophercloud.Client, id string) string {
	return client.ServiceURL(apiName, id)
}

func createURL(client gophercloud.Client) string {
	return commonURL(client)
}

func deleteURL(client gophercloud.Client, id string) string {
	return idURL(client, id)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("clusters", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("clusters")
}

func listDetailURL(client gophercloud.Client) string {
	return client.ServiceURL("clusters", "detail")
}

func updateURL(client gophercloud.Client, id string) string {
	return idURL(client, id)
}

func upgradeURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("clusters", id, "actions/upgrade")
}

func resizeURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("clusters", id, "actions/resize")
}
