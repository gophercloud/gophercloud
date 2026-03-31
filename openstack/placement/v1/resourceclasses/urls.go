package resourceclasses

import (
	"github.com/gophercloud/gophercloud/v2"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("resource_classes")
}

func getURL(client *gophercloud.ServiceClient, name string) string {
	return client.ServiceURL("resource_classes", name)
}
