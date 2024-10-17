package snapshots

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("snapshots")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("snapshots", "detail")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id)
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id)
}

func resetStatusURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}

func forceDeleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}
