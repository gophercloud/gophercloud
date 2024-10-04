package sharenetworks

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("share-networks")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-networks", id)
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("share-networks", "detail")
}

func getURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func updateURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func addSecurityServiceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-networks", id, "action")
}

func removeSecurityServiceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-networks", id, "action")
}
