package replicas

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("share-replicas")
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("share-replicas")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("share-replicas", "detail")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-replicas", id)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-replicas", id)
}

func listExportLocationsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-replicas", id, "export-locations")
}

func getExportLocationURL(c gophercloud.Client, replicaID, id string) string {
	return c.ServiceURL("share-replicas", replicaID, "export-locations", id)
}

func actionURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-replicas", id, "action")
}
