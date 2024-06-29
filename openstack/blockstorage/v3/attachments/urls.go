package attachments

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("attachments")
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("attachments", "detail")
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("attachments", id)
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("attachments", id)
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("attachments", id)
}

func completeURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("attachments", id, "action")
}
