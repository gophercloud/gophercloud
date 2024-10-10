package sharetransfers

import "github.com/gophercloud/gophercloud/v2"

func transferURL(c gophercloud.Client) string {
	return c.ServiceURL("share-transfers")
}

func acceptURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-transfers", id, "accept")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-transfers", id)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("share-transfers")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("share-transfers", "detail")
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("share-transfers", id)
}
