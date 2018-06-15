package clustertemplates

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

var apiVersion = "v1"
var apiName = "clustertemplates"

func commonURL(client *gophercloud.ServiceClient) string {
	if strings.HasSuffix(client.ResourceBaseURL(), apiVersion+"/") {
		return client.ServiceURL(apiName)
	} else {
		return client.ServiceURL(apiVersion, apiName)
	}
}

func createURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}
