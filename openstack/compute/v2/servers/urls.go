package servers

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("servers")
}

func listURL(client gophercloud.Client) string {
	return createURL(client)
}

func listDetailURL(client gophercloud.Client) string {
	return client.ServiceURL("servers", "detail")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id)
}

func getURL(client gophercloud.Client, id string) string {
	return deleteURL(client, id)
}

func updateURL(client gophercloud.Client, id string) string {
	return deleteURL(client, id)
}

func actionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "action")
}

func metadatumURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("servers", id, "metadata", key)
}

func metadataURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "metadata")
}

func listAddressesURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "ips")
}

func listAddressesByNetworkURL(client gophercloud.Client, id, network string) string {
	return client.ServiceURL("servers", id, "ips", network)
}

func passwordURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "os-server-password")
}
