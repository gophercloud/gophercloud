package servers

import "github.com/rackspace/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("servers")
}

func detailURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("servers", "detail")
}

func serverURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("servers", id)
}

func actionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}
