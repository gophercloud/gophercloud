package registeredlimits

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath             = "registered_limits"
	enforcementModelPath = "model"
)

func rootURL(client gophercloud.Client) string {
	return client.ServiceURL(rootPath)
}

func resourceURL(client gophercloud.Client, id string) string {
	return client.ServiceURL(rootPath, id)
}
