package introspection

import "github.com/bizflycloud/gophercloud"

func listIntrospectionsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("introspection")
}

func introspectionURL(client *gophercloud.ServiceClient, nodeID string) string {
	return client.ServiceURL("introspection", nodeID)
}

func abortIntrospectionURL(client *gophercloud.ServiceClient, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "abort")
}

func introspectionDataURL(client *gophercloud.ServiceClient, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "data")
}

func introspectionUnprocessedDataURL(client *gophercloud.ServiceClient, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "data", "unprocessed")
}
