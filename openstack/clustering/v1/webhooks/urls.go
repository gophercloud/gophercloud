package webhooks

import "github.com/gophercloud/gophercloud/v2"

func triggerURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("v1", "webhooks", id, "trigger")
}
