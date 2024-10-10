package policies

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "fwaas"
	resourcePath = "firewall_policies"
	insertPath   = "insert_rule"
	removePath   = "remove_rule"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func insertURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, insertPath)
}

func removeURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, removePath)
}
