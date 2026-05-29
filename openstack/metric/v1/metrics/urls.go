package metrics

import "github.com/gophercloud/gophercloud/v2"

func queryURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("query")
}

func labelsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("labels")
}

func labelValuesURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL("label", name, "values")
}

func seriesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("series")
}

func targetsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("targets")
}

func runtimeInfoURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("status", "runtimeinfo")
}

func cleanTombstonesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("admin", "tsdb", "clean_tombstones")
}

func deleteSeriesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("admin", "tsdb", "delete_series")
}

func snapshotURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("admin", "tsdb", "snapshot")
}
