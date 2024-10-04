package routers

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "routers"

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func addInterfaceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "add_router_interface")
}

func removeInterfaceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_router_interface")
}

func listl3AgentsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "l3-agents")
}
