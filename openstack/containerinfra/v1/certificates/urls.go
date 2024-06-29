package certificates

import (
	"github.com/gophercloud/gophercloud/v2"
)

var apiName = "certificates"

func commonURL(client gophercloud.Client) string {
	return client.ServiceURL(apiName)
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL(apiName, id)
}

func createURL(client gophercloud.Client) string {
	return commonURL(client)
}

func updateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL(apiName, id)
}
