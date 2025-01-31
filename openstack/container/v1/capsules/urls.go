package capsules

import "github.com/gophercloud/gophercloud/v2"

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("capsules", id)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("capsules")
}

// `listURL` is a pure function. `listURL(c)` is a URL for which a GET
// request will respond with a list of capsules in the service `c`.
func listURL(c gophercloud.Client) string {
	return c.ServiceURL("capsules")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("capsules", id)
}
