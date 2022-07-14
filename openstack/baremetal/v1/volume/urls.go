package bmvolume

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("volume")
}

func createConnectorsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("volume", "connectors")
}

func listConnectorsURL(client *gophercloud.ServiceClient) string {
	return createConnectorsURL(client)
}

func getConnectorsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("volume", "connectors", id)
}

func patchConnectorsURL(client *gophercloud.ServiceClient, id string) string {
	return getConnectorsURL(client, id)
}

func deleteConnectorsURL(client *gophercloud.ServiceClient, id string) string {
	return getConnectorsURL(client, id)
}

func createTargetsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("volume", "targets")
}

func listTargetsURL(client *gophercloud.ServiceClient) string {
	return createTargetsURL(client)
}

func getTargetsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("volume", "targets", id)
}

func patchTargetsURL(client *gophercloud.ServiceClient, id string) string {
	return getTargetsURL(client, id)
}

func deleteTargetsURL(client *gophercloud.ServiceClient, id string) string {
	return getTargetsURL(client, id)
}
