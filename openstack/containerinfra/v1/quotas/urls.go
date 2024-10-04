package quotas

import (
	"github.com/gophercloud/gophercloud/v2"
)

var apiName = "quotas"

func commonURL(client gophercloud.Client) string {
	return client.ServiceURL(apiName)
}

func createURL(client gophercloud.Client) string {
	return commonURL(client)
}
