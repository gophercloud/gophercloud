package introspection

import "github.com/gophercloud/gophercloud/v2"

func listIntrospectionsURL(client gophercloud.Client) string {
	return client.ServiceURL("introspection")
}

func introspectionURL(client gophercloud.Client, nodeID string) string {
	return client.ServiceURL("introspection", nodeID)
}

func abortIntrospectionURL(client gophercloud.Client, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "abort")
}

func introspectionDataURL(client gophercloud.Client, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "data")
}

func introspectionUnprocessedDataURL(client gophercloud.Client, nodeID string) string {
	return client.ServiceURL("introspection", nodeID, "data", "unprocessed")
}
