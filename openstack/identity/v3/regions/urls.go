package regions

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("regions")
}

func getURL(client gophercloud.Client, regionID string) string {
	return client.ServiceURL("regions", regionID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("regions")
}

func updateURL(client gophercloud.Client, regionID string) string {
	return client.ServiceURL("regions", regionID)
}

func deleteURL(client gophercloud.Client, regionID string) string {
	return client.ServiceURL("regions", regionID)
}
