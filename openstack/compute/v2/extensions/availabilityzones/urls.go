// Package availabilityzones provides information and interaction with Availability zones
// extension that works with the OpenStack Compute service.
package availabilityzones

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-availability-zone")
}
