package stacks

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("stacks")
}

func adoptURL(c gophercloud.Client) string {
	return createURL(c)
}

func listURL(c gophercloud.Client) string {
	return createURL(c)
}

func getURL(c gophercloud.Client, name, id string) string {
	return c.ServiceURL("stacks", name, id)
}

func findURL(c gophercloud.Client, identity string) string {
	return c.ServiceURL("stacks", identity)
}

func updateURL(c gophercloud.Client, name, id string) string {
	return getURL(c, name, id)
}

func deleteURL(c gophercloud.Client, name, id string) string {
	return getURL(c, name, id)
}

func previewURL(c gophercloud.Client) string {
	return c.ServiceURL("stacks", "preview")
}

func abandonURL(c gophercloud.Client, name, id string) string {
	return c.ServiceURL("stacks", name, id, "abandon")
}
