package secrets

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("secrets")
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("secrets", id)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("secrets")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("secrets", id)
}

func updateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("secrets", id)
}

func payloadURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("secrets", id, "payload")
}

func metadataURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("secrets", id, "metadata")
}

func metadatumURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("secrets", id, "metadata", key)
}
