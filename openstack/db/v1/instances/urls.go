package instances

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL("instances")
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("instances", id)
}

func userRootURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("instances", id, "root")
}

func actionURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("instances", id, "action")
}
