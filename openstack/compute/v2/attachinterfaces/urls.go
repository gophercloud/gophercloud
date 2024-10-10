package attachinterfaces

import "github.com/gophercloud/gophercloud/v2"

func listInterfaceURL(client gophercloud.Client, serverID string) string {
	return client.ServiceURL("servers", serverID, "os-interface")
}

func getInterfaceURL(client gophercloud.Client, serverID, portID string) string {
	return client.ServiceURL("servers", serverID, "os-interface", portID)
}

func createInterfaceURL(client gophercloud.Client, serverID string) string {
	return client.ServiceURL("servers", serverID, "os-interface")
}
func deleteInterfaceURL(client gophercloud.Client, serverID, portID string) string {
	return client.ServiceURL("servers", serverID, "os-interface", portID)
}
