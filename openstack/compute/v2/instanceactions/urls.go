package instanceactions

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "os-instance-actions")
}

func instanceActionsURL(client gophercloud.Client, serverID, requestID string) string {
	return client.ServiceURL("servers", serverID, "os-instance-actions", requestID)
}
