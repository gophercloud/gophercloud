package stackevents

import "github.com/gophercloud/gophercloud/v2"

func findURL(c gophercloud.Client, stackName string) string {
	return c.ServiceURL("stacks", stackName, "events")
}

func listURL(c gophercloud.Client, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "events")
}

func listResourceEventsURL(c gophercloud.Client, stackName, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources", resourceName, "events")
}

func getURL(c gophercloud.Client, stackName, stackID, resourceName, eventID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "resources", resourceName, "events", eventID)
}
