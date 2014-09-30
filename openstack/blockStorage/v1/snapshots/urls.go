package snapshots

import "github.com/rackspace/gophercloud"

func snapshotsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("snapshots")
}

func snapshotURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id)
}
