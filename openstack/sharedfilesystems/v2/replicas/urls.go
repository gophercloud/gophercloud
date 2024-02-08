package replicas

import "github.com/gophercloud/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-replicas")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-replicas")
}

func listDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-replicas", "detail")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-replicas", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-replicas", id)
}

func listExportLocationsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-replicas", id, "export-locations")
}

func getExportLocationURL(c *gophercloud.ServiceClient, replicaID, id string) string {
	return c.ServiceURL("share-replicas", replicaID, "export-locations", id)
}

func actionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-replicas", id, "action")
}
