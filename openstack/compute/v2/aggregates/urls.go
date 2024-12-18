package aggregates

import "github.com/gophercloud/gophercloud/v2"

func aggregatesListURL(c gophercloud.Client) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesCreateURL(c gophercloud.Client) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesDeleteURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesGetURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesUpdateURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesAddHostURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}

func aggregatesRemoveHostURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}

func aggregatesSetMetadataURL(c gophercloud.Client, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID, "action")
}
