package registeredlimits

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath             = "registered_limits"
	enforcementModelPath = "model"
)

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

func resourceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL(rootPath, id)
}
