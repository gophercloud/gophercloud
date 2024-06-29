package amphorae

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "octavia"
	resourcePath = "amphorae"
	failoverPath = "failover"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func failoverRootURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, failoverPath)
}
