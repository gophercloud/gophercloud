package limits

import "github.com/bizflycloud/gophercloud"

func enforcementModelURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("limits", "model")
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("limits")
}
