package stacks

import "github.com/rackspace/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("stacks")
}

func adoptURL(c *gophercloud.ServiceClient) string {
	return createURL(c)
}

func listURL(c *gophercloud.ServiceClient) string {
	return createURL(c)
}

func getURL(c *gophercloud.ServiceClient, name, id string) string {
	return c.ServiceURL("stacks", name, id)
}

func updateURL(c *gophercloud.ServiceClient, name, id string) string {
	return getURL(c, name, id)
}

func deleteURL(c *gophercloud.ServiceClient, name, id string) string {
	return getURL(c, name, id)
}

func previewURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("stacks", "preview")
}

func abandonURL(c *gophercloud.ServiceClient, name, id string) string {
	return c.ServiceURL("stacks", name, id, "abandon")
}

func createSnapshotURL(c *gophercloud.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "snapshots")
}

func listSnapshotsURL(c *gophercloud.ServiceClient, stackName, stackID string) string {
	return createSnapshotURL(c, stackName, stackID)
}

func getSnapshotURL(c *gophercloud.ServiceClient, stackName, stackID, snapshotID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "snapshots", snapshotID)
}

func restoreSnapshotURL(c *gophercloud.ServiceClient, stackName, stackID, snapshotID string) string {
	return getSnapshotURL(c, stackName, stackID, snapshotID)
}

func deleteSnapshotURL(c *gophercloud.ServiceClient, stackName, stackID, snapshotID string) string {
	return getSnapshotURL(c, stackName, stackID, snapshotID)
}
