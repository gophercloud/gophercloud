package executions

import "github.com/gophercloud/gophercloud"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("executions")
}
