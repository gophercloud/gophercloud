package snapshots

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("snapshots")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id)
}

func getURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func listURL(c gophercloud.Client) string {
	return createURL(c)
}

func updateURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func metadataURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "metadata")
}

func updateMetadataURL(c gophercloud.Client, id string) string {
	return metadataURL(c, id)
}

func resetStatusURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}

func updateStatusURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}

func forceDeleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}
