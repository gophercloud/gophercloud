package credentials

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("credentials")
}

func getURL(client gophercloud.Client, credentialID string) string {
	return client.ServiceURL("credentials", credentialID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("credentials")
}

func deleteURL(client gophercloud.Client, credentialID string) string {
	return client.ServiceURL("credentials", credentialID)
}

func updateURL(client gophercloud.Client, credentialID string) string {
	return client.ServiceURL("credentials", credentialID)
}
