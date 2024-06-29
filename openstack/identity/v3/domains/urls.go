package domains

import "github.com/gophercloud/gophercloud/v2"

func listAvailableURL(client gophercloud.Client) string {
	return client.ServiceURL("auth", "domains")
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("domains")
}

func getURL(client gophercloud.Client, domainID string) string {
	return client.ServiceURL("domains", domainID)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("domains")
}

func deleteURL(client gophercloud.Client, domainID string) string {
	return client.ServiceURL("domains", domainID)
}

func updateURL(client gophercloud.Client, domainID string) string {
	return client.ServiceURL("domains", domainID)
}
