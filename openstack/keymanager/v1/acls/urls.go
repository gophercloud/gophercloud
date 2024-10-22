package acls

import "github.com/gophercloud/gophercloud/v2"

func containerURL(client gophercloud.Client, containerID string) string {
	return client.ServiceURL("containers", containerID, "acl")
}

func secretURL(client gophercloud.Client, secretID string) string {
	return client.ServiceURL("secrets", secretID, "acl")
}
