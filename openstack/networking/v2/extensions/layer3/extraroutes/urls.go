package extraroutes

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "routers"

func addExtraRoutesURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "add_extraroutes")
}

func removeExtraRoutesURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_extraroutes")
}
