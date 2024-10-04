package extensions

import "github.com/gophercloud/gophercloud/v2"

func ActionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "action")
}
