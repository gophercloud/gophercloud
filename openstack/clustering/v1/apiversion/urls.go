package apiversion

import "github.com/gophercloud/gophercloud"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL()
}

func versionURL(client *gophercloud.ServiceClient, version string) string {
	return client.ServiceURL(version)
}

func listURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}

func getURL(client *gophercloud.ServiceClient, version string) string {
	return versionURL(client, version)
}
