package throttle

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

const (
	path    = "loadbalancers"
	sslPath = "ssltermination"
)

func rootURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL(path, strconv.Itoa(id), sslPath)
}
