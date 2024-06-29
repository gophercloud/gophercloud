package clustertemplates

import (
	"github.com/gophercloud/gophercloud/v2"
)

var apiName = "clustertemplates"

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

func listURL(client gophercloud.Client) string {
	return commonURL(client)
}

func getURL(client gophercloud.Client, id string) string {
	return idURL(client, id)
}

func updateURL(client gophercloud.Client, id string) string {
	return idURL(client, id)
}
