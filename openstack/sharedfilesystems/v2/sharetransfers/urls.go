package sharetransfers

import "github.com/gophercloud/gophercloud/v2"

func transferURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-transfers")
}

func acceptURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-transfers", id, "accept")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-transfers", id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-transfers")
}

func listDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-transfers", "detail")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-transfers", id)
}
