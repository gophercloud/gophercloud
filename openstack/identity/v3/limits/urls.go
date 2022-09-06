package limits

import "github.com/gophercloud/gophercloud"

const (
	rootPath             = "limits"
	enforcementModelPath = "model"
)

func enforcementModelURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rootPath, enforcementModelPath)
}

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
