package endpointgroups

import "github.com/gophercloud/gophercloud"

const (
	endpointGroupPath = "OS-EP-FILTER/endpoint_groups"
)

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(endpointGroupPath)
}

func resourceURL(client *gophercloud.ServiceClient, endointGroupID string) string {
	return client.ServiceURL(endpointGroupPath, endointGroupID)
}
